package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
)

func (s *Server) handleYTInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	mediaUrl := r.URL.Query().Get("url")
	if mediaUrl == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	details, ok := s.dlapi.Cache.Data[mediaUrl]
	if ok {
		log.Println("Cache Hit")
		if err := json.NewEncoder(w).Encode(details); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	log.Println("Cache Miss")
	cmd := dlapi.GetVideoInfoCmd(mediaUrl)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to run yt-dlp: %v", err), http.StatusInternalServerError)
		return
	}

	if err := dlapi.ParseYTMediaInfo(output, &details); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing media info: %v", err), http.StatusInternalServerError)
		return
	}

	s.dlapi.Cache.Data[mediaUrl] = details
	if err := s.dlapi.Cache.SaveCache(); err != nil {
		log.Printf("Failed to save cache: %v", err)
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
