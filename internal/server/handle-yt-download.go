package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
)

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

	mediaInfo, ok := s.dlapi.Cache.Data[mediaURL]
	if !ok {
		http.Error(w, "Please fetch video information before downloading", http.StatusBadRequest)
		return
	}

	log.Println(formatID)

	var cmd *exec.Cmd
	if strings.Contains(formatID, "+") {
		formatIDs := strings.Split(formatID, "+")
		if len(formatIDs) == 2 {
			audioFormatID := formatIDs[0]
			videoFormatID := formatIDs[1]
			cmd = dlapi.GetMediaByFormatIDS(mediaURL, audioFormatID, videoFormatID)
		} else {
			http.Error(w, "Invalid media format options", http.StatusBadRequest)
			return
		}
	} else {
		cmd = dlapi.GetMediaByFormatID(mediaURL, formatID)
	}

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to create stdout pipe", http.StatusInternalServerError)
		return
	}
	defer stdout.Close()

	// Get stderr pipe for logging
	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, "Failed to create stderr pipe", http.StatusInternalServerError)
		return
	}

	// Start download
	if err = cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start download: %v", err), http.StatusInternalServerError)
		return
	}

	// Log stderr in background
	go func() {
		if _, err = io.Copy(os.Stderr, stderr); err != nil {
			log.Printf("Error reading stderr: %v", err)
		}
	}()

	// Set headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.%s"`, sanitizeFilename(mediaInfo.Title), "mp4"))
	w.Header().Set("Transfer-Encoding", "chunked")

	// Stream to response
	written, err := io.Copy(w, stdout)
	if err != nil {
		log.Printf("Error streaming download: %v", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command error: %v", err)
		return
	}

	log.Printf("Successfully downloaded %s (%d bytes)", mediaInfo.Title, written)
}

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
