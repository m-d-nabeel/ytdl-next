package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/m-d-nabeel/ytdl-web/internal/api"
)

type DownloadResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	URL    string `json:"url,omitempty"`
}

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	wg         sync.WaitGroup
	api        *api.API
	port       string
}

// NewServer creates a new server instance with the given config
func NewServer(api *api.API) *Server {
	s := &Server{
		router: http.NewServeMux(),
		api:    api,
		port:   ":8080",
	}

	s.setupRoutes()
	return s
}

// Start initializes and starts the server
func (s *Server) Start() error {
	// Serve static files from the public directory
	http.Handle("/", http.FileServer(http.Dir("public")))

	// API endpoints
	http.HandleFunc("/api/download", s.handleDownload)

	log.Printf("Server starting on port %s...", s.port)
	return http.ListenAndServe(s.port, nil)
}

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

// getContentType determines the content type based on file extension
func (s *Server) getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".ogg":
		return "audio/ogg"
	case ".webm":
		return "video/webm"
	default:
		return "application/octet-stream"
	}
}
