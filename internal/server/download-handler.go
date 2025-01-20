package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")

	cmd := exec.Command("yt-dlp", "-o", "-", requestData.URL)
	cmd.Stdout = w

	if err := cmd.Start(); err != nil {
		json.NewEncoder(w).Encode(DownloadResponse{
			Status: "error",
			Error:  fmt.Sprintf("Failed to start download: %v", err),
		})
		return
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Download error: %v", err)
	}
}
