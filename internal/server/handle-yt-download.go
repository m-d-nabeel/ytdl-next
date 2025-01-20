package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/m-d-nabeel/ytdl-web/internal/api"
)

func (s *Server) handleYTDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mediaUrl := r.URL.Query().Get("url")
	formatID := r.URL.Query().Get("format_id")
	if mediaUrl == "" || formatID == "" {
		http.Error(w, "URL and formatID parameter is required", http.StatusBadRequest)
		return
	}

	mediaInfo, ok := s.api.Cache.Data[mediaUrl]
	if !ok {
		http.Error(w, "Please fetch video information before downloading", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.mp4"`, sanitizeString(mediaInfo.Title)))
	w.Header().Set("Transfer-Encoding", "chunked")

	cmd := api.GetMediaByFormatID(mediaUrl, formatID)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to create stdout pipe", http.StatusInternalServerError)
		return
	}
	defer stdout.Close()

	if err = cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start download: %v", err), http.StatusInternalServerError)
		return
	}

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

func sanitizeString(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' ||
			r == ' ' {
			return r
		}
		return '-'
	}, s)
}
