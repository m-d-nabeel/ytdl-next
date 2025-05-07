package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Temporary extension for downloads in progress
const DownloadInProgressExt = ".ytdlp"

func (s *Server) handleYTDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mediaURL := r.URL.Query().Get("url")
	formatID := r.URL.Query().Get("format_id")

	if mediaURL == "" || formatID == "" {
		http.Error(w, "URL and format_id parameters are required", http.StatusBadRequest)
		return
	}

	// Get media info from cache
	mediaInfo, ok := s.dlapi.Cache.Data[mediaURL]
	if !ok {
		http.Error(w, "Please fetch video information before downloading", http.StatusBadRequest)
		return
	}

	log.Printf("Download request received for %s with format %s", mediaInfo.Title, formatID)

	// Get number of active downloads for rate limiting information
	activeDownloads := s.downloadManager.GetActiveDownloadCount()
	log.Printf("Current active downloads: %d", activeDownloads)

	// Create a context that is canceled when the client disconnects
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Monitor for client disconnection using CloseNotifier pattern
	// This works with Go 1.17+ using context Done channel from request Context
	clientGone := r.Context().Done()

	// Start a goroutine to monitor client disconnection
	go func() {
		select {
		case <-clientGone:
			log.Printf("Client disconnected, canceling download for %s", mediaInfo.Title)
			cancel() // Cancel our download context
		case <-ctx.Done():
			// Our context was cancelled elsewhere, no action needed
		}
	}()

	// Create a channel to receive the download response
	responseChan := make(chan *DownloadResponse, 1)

	// Create and enqueue the download request
	downloadReq := &DownloadRequest{
		MediaURL:     mediaURL,
		FormatID:     formatID,
		Title:        mediaInfo.Title,
		ResponseChan: responseChan,
		Context:      ctx,
	}

	// Try to enqueue the download
	err := s.downloadManager.EnqueueDownload(downloadReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Server busy: %v", err), http.StatusServiceUnavailable)
		return
	}

	// Wait for response from worker
	select {
	case <-ctx.Done():
		http.Error(w, "Download canceled or timed out", http.StatusRequestTimeout)
		return
	case downloadResp := <-responseChan:
		if downloadResp.Status != "success" || downloadResp.Reader == nil {
			http.Error(w, fmt.Sprintf("Download failed: %s", downloadResp.Error), http.StatusInternalServerError)
			return
		}

		// Set up HTTP response headers for proper download
		filename := sanitizeFilename(mediaInfo.Title)
		extension := guessFileExtension(formatID)

			// Create the final and temporary filenames
		finalFilename := fmt.Sprintf("%s.%s", filename, extension)
		tempFilename := fmt.Sprintf("%s.%s%s", filename, extension, DownloadInProgressExt)

		// Set content type based on file extension
		w.Header().Set("Content-Type", s.getContentType(finalFilename))

		// Use temporary extension for downloads in progress
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, tempFilename))

		// Add custom header to indicate final filename for potential client-side renaming
		w.Header().Set("X-Final-Filename", finalFilename)

		// If we have a file size, set Content-Length header for better download experience
		if downloadResp.FileSize > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", downloadResp.FileSize))
			// Don't use chunked encoding when we have a content length
		} else {
			// Use chunked encoding when size is unknown
			w.Header().Set("Transfer-Encoding", "chunked")
		}

		// Additional headers for better download experience
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Accept-Ranges", "bytes") // Support resumable downloads

		// Create a custom writer that detects if the client disconnects during writing
		// and calls the cancel function if that happens
		clientWriter := &clientAwareWriter{
			ResponseWriter: w,
			cancel:         downloadResp.CancelFn,
		}

		// Stream the response to the client
		bytesWritten, err := io.Copy(clientWriter, downloadResp.Reader)

		// Log whether the download completed or was canceled
		if ctx.Err() != nil {
			log.Printf("Download for %s was canceled after %d bytes", mediaInfo.Title, bytesWritten)
		} else if err != nil {
			log.Printf("Error streaming download to client: %v", err)
		} else {
			log.Printf("Successfully streamed %d bytes for %s", bytesWritten, mediaInfo.Title)
		}

		// Check for errors from the download process
		select {
		case err := <-downloadResp.ErrorCh:
			if err != nil && ctx.Err() == nil {
				// Only log if it wasn't a cancellation
				log.Printf("Error during download process: %v", err)
			}
		default:
			// No error or not ready yet
		}
	}
}

// clientAwareWriter wraps an http.ResponseWriter to detect client disconnections
type clientAwareWriter struct {
	http.ResponseWriter
	cancel func()
}

// Write overrides the standard Write method to check for client disconnection
func (w *clientAwareWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	if err != nil {
		// If there's an error writing, the client may have disconnected
		if strings.Contains(err.Error(), "broken pipe") ||
			strings.Contains(err.Error(), "reset by peer") ||
			strings.Contains(err.Error(), "connection closed") {
			log.Println("Client disconnected during write, canceling download")
			if w.cancel != nil {
				w.cancel() // Call cancel function to stop the download
			}
		}
	}
	return n, err
}

// Other existing functions remain unchanged
func sanitizeFilename(s string) string {
	// Replace invalid characters with hyphens
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' ||
			r == ' ' {
			return r
		}
		return '-'
	}, s)

	// Trim spaces and hyphens from ends
	return strings.Trim(safe, " -")
}

func guessFileExtension(formatID string) string {
	// Some heuristics to guess file extension based on format ID

	// Audio-only formats typically have lower IDs
	if len(formatID) <= 3 && !strings.Contains(formatID, "+") {
		return "mp3"
	}

	// Video formats or compound formats
	return "mp4"
}
