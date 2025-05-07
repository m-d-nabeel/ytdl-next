// filepath: /home/m-d-nabeel/Projects/ytdl-next/internal/server/download_manager.go
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DownloadRequest represents a single download request
type DownloadRequest struct {
	MediaURL     string
	FormatID     string
	Title        string
	ResponseChan chan *DownloadResponse
	Context      context.Context
}

// DownloadManager handles concurrent downloads
type DownloadManager struct {
	maxWorkers       int
	downloadQueue    chan *DownloadRequest
	activeDownloads  int
	activeDownloadMu sync.Mutex
	workerWg         sync.WaitGroup
	shutdownCh       chan struct{}
}

// NewDownloadManager creates a new download manager with the specified number of workers
func NewDownloadManager(maxWorkers int, queueSize int) *DownloadManager {
	if maxWorkers <= 0 {
		maxWorkers = 3 // Default to 3 workers if not specified
	}
	if queueSize <= 0 {
		queueSize = 50 // Default queue size
	}

	dm := &DownloadManager{
		maxWorkers:    maxWorkers,
		downloadQueue: make(chan *DownloadRequest, queueSize),
		shutdownCh:    make(chan struct{}),
	}

	// Start worker pool
	for i := 0; i < maxWorkers; i++ {
		dm.workerWg.Add(1)
		go dm.worker(i)
	}

	log.Printf("Download manager started with %d workers and queue size %d", maxWorkers, queueSize)
	return dm
}

// EnqueueDownload adds a download request to the queue
func (dm *DownloadManager) EnqueueDownload(req *DownloadRequest) error {
	// Use a non-blocking send to avoid hanging if queue is full
	select {
	case dm.downloadQueue <- req:
		return nil
	default:
		return fmt.Errorf("download queue is full, try again later")
	}
}

// Shutdown gracefully shuts down the download manager
func (dm *DownloadManager) Shutdown(timeout time.Duration) {
	log.Println("Shutting down download manager...")
	close(dm.shutdownCh)

	// Use a timeout context for shutdown
	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Wait for workers to finish with timeout
	doneCh := make(chan struct{})
	go func() {
		dm.workerWg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		log.Println("Download manager shut down gracefully")
	case <-timeoutCtx.Done():
		log.Println("Download manager shutdown timed out, some downloads may be interrupted")
	}
}

// GetActiveDownloadCount returns the current number of active downloads
func (dm *DownloadManager) GetActiveDownloadCount() int {
	dm.activeDownloadMu.Lock()
	defer dm.activeDownloadMu.Unlock()
	return dm.activeDownloads
}

// worker processes download requests from the queue
func (dm *DownloadManager) worker(id int) {
	defer dm.workerWg.Done()
	log.Printf("Download worker %d started", id)

	for {
		select {
		case <-dm.shutdownCh:
			log.Printf("Download worker %d shutting down", id)
			return
		case req := <-dm.downloadQueue:
			dm.activeDownloadMu.Lock()
			dm.activeDownloads++
			dm.activeDownloadMu.Unlock()

			log.Printf("Worker %d processing download for: %s", id, req.Title)
			result := dm.processDownload(req)
			// Send result back through response channel
			req.ResponseChan <- result

			dm.activeDownloadMu.Lock()
			dm.activeDownloads--
			dm.activeDownloadMu.Unlock()
		}
	}
}

// processDownload handles the actual download process
func (dm *DownloadManager) processDownload(req *DownloadRequest) *DownloadResponse {
	// Check if context is already canceled
	if req.Context.Err() != nil {
		return &DownloadResponse{
			Status: "error",
			Error:  "Request canceled",
		}
	}

	// Get the file size first using yt-dlp to retrieve file information
	fileSize := getFileSize(req.MediaURL, req.FormatID)
	log.Printf("Estimated file size for %s: %d bytes", req.Title, fileSize)

	// Create a pipe to capture the download output
	reader, writer := io.Pipe()

	// Create a cancellable context for the download process
	downloadCtx, cancelDownload := context.WithCancel(req.Context)

	// Start the download in a goroutine
	errCh := make(chan error, 1)
	var cmd *exec.Cmd

	// Set up and start the download command
	go func() {
		defer writer.Close()

		// Create the command with context for proper cancellation
		cmd = createDownloadCommandWithContext(downloadCtx, req.MediaURL, req.FormatID)
		if cmd == nil {
			errCh <- fmt.Errorf("failed to create download command")
			return
		}

		// Get stdout pipe
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			errCh <- fmt.Errorf("failed to create stdout pipe: %w", err)
			return
		}
		defer stdoutPipe.Close()

		// Get stderr pipe for logging
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			errCh <- fmt.Errorf("failed to create stderr pipe: %w", err)
			return
		}
		defer stderrPipe.Close()

		// Log stderr in background
		go func() {
			if _, err = io.Copy(os.Stderr, stderrPipe); err != nil {
				log.Printf("Error reading stderr: %v", err)
			}
		}()

		// Start the command
		if err = cmd.Start(); err != nil {
			errCh <- fmt.Errorf("failed to start download: %w", err)
			return
		}

		// Copy the output to our pipe
		written, err := io.Copy(writer, stdoutPipe)
		if err != nil {
			// Check if the error is due to cancellation
			if downloadCtx.Err() != nil {
				errCh <- fmt.Errorf("download canceled: %w", downloadCtx.Err())
			} else {
				errCh <- fmt.Errorf("error streaming download: %w", err)
			}
			return
		}

		log.Printf("Download completed: %s (%d bytes)", req.Title, written)

		// Wait for command to complete
		if err := cmd.Wait(); err != nil {
			// Check if this was due to cancellation
			if downloadCtx.Err() != nil {
				errCh <- fmt.Errorf("download canceled: %w", downloadCtx.Err())
			} else {
				errCh <- fmt.Errorf("command error: %w", err)
			}
			return
		}

		errCh <- nil
	}()

	return &DownloadResponse{
		Status:   "success",
		Error:    "",
		URL:      "", // This will be replaced with actual URL in the handler
		Reader:   reader,
		ErrorCh:  errCh,
		FileSize: fileSize,
		CancelFn: cancelDownload,
	}
}

// getFileSize retrieves the estimated file size for a given format ID
func getFileSize(mediaURL, formatID string) int64 {
	// If this is a compound format (audio+video), don't try to get accurate size
	// as the final muxed size can be different than the sum of parts
	if isCompoundFormat(formatID) {
		return -1 // Use chunked transfer encoding instead
	}

	// First, try to get file size from the --print format option
	args := []string{
		"--print", "%(filesize)s",
		"-f", formatID,
		"--no-playlist",
		"--no-warnings",
		"--no-download",
	}

	cmd := exec.Command("yt-dlp", append(args, mediaURL)...)
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		// Trim whitespace and parse as int64
		sizeStr := strings.TrimSpace(string(output))
		// Check if it's "NA" or empty
		if sizeStr != "NA" && sizeStr != "" {
			size, err := strconv.ParseInt(sizeStr, 10, 64)
			if err == nil && size > 0 {
				return size
			}
		}
	}

	// Fall back to best estimate if exact size can't be determined
	// Many formats report filesize in metadata
	args = []string{
		"-j",
		"-f", formatID,
		"--no-playlist",
		"--no-warnings",
		"--no-download",
	}

	cmd = exec.Command("yt-dlp", append(args, mediaURL)...)
	output, err = cmd.Output()
	if err == nil && len(output) > 0 {
		// Parse JSON to get filesize
		var result struct {
			Filesize       int64 `json:"filesize"`
			FilesizeApprox int64 `json:"filesize_approx"`
		}

		if err := json.Unmarshal(output, &result); err == nil {
			if result.Filesize > 0 {
				return result.Filesize
			}
			if result.FilesizeApprox > 0 {
				return result.FilesizeApprox
			}
		}
	}

	// If we still don't have a size, return a default estimate
	return -1 // Negative value indicates unknown size
}

// Add these helper methods for command creation
func isCompoundFormat(formatID string) bool {
	return containsChar(formatID, '+')
}

func parseCompoundFormat(formatID string) []string {
	return splitByChar(formatID, '+')
}

// These methods would need to be implemented based on your existing code
func containsChar(s string, c byte) bool {
	for i := range len(s) {
		if s[i] == c {
			return true
		}
	}
	return false
}

func splitByChar(s string, c byte) []string {
	var result []string
	start := 0
	for i := range len(s) {
		if s[i] == c {
			if i > start {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

// createDownloadCommandWithContext creates the appropriate download command with context for proper cancellation
func createDownloadCommandWithContext(ctx context.Context, mediaURL, formatID string) *exec.Cmd {
	// If format contains '+', it's a combined audio/video download
	if isCompoundFormat(formatID) {
		ids := parseCompoundFormat(formatID)
		if len(ids) != 2 {
			return nil
		}
		return createCombinedDownloadCmdWithContext(ctx, mediaURL, ids[0], ids[1])
	}
	return createSingleFormatDownloadCmdWithContext(ctx, mediaURL, formatID)
}

// Helper functions for command creation with context
func createSingleFormatDownloadCmdWithContext(ctx context.Context, mediaURL, formatID string) *exec.Cmd {
	args := []string{
		"-f", formatID,
		"--progress-template", "",
		"--throttled-rate", "100K",
		"--no-playlist",
		"--no-warnings",
		"--no-cookies",
		"-o", "-", // Output to stdout
	}

	return exec.CommandContext(ctx, "yt-dlp", append(args, mediaURL)...)
}

func createCombinedDownloadCmdWithContext(ctx context.Context, mediaURL, audioFormatID, videoFormatID string) *exec.Cmd {
	args := []string{
		"-f", fmt.Sprintf("%v+%v", audioFormatID, videoFormatID),
		"--concurrent-fragments", "4", // Parallel download
		"--audio-multistreams", // Allow multiple audio streams
		"--video-multistreams", // Allow multiple video streams
		"--no-playlist",
		"--no-warnings",
		"--no-cookies",
		"-o", "-", // Output to stdout
	}

	return exec.CommandContext(ctx, "yt-dlp", append(args, mediaURL)...)
}
