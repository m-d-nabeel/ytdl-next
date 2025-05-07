package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
	"github.com/m-d-nabeel/ytdl-web/internal/types"
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

	// Check if the Cache is initialized
	if s.dlapi.Cache == nil {
		log.Println("Cache not available")
		fetchAndServeMediaInfo(s, w, r, mediaUrl)
		return
	}

	// Try to get the media info from the cache
	details, err := s.dlapi.Cache.Get(mediaUrl)
	if err == nil && details != nil {
		log.Println("Cache Hit")
		if err := json.NewEncoder(w).Encode(*details); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
		return
	}

	log.Println("Cache Miss")
	fetchAndServeMediaInfo(s, w, r, mediaUrl)
}

func fetchAndServeMediaInfo(s *Server, w http.ResponseWriter, r *http.Request, mediaUrl string) {
	var details types.YTMediaInfo

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

	// Save to cache if available
	if s.dlapi.Cache != nil {
		if err := s.dlapi.Cache.Set(mediaUrl, details); err != nil {
			log.Printf("Failed to save to cache: %v", err)
		}
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
