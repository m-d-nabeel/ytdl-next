package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleYTInfo(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the "url" parameter from the query string
	mediaUrl := r.URL.Query().Get("url")
	if mediaUrl == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch YouTube media information
	info, err := s.api.GetYTMediaInfo(mediaUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
  s.api.Cache.SaveCache()

	// Send a JSON response with the media info
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
