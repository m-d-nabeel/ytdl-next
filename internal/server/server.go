package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
)

type DownloadResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	URL    string `json:"url,omitempty"`
}

type Server struct {
	router  *http.ServeMux
	dlapi   *dlapi.DLAPI
	port    string
	devMode bool
}

// NewServer creates a new server instance with the given config
func NewServer(dlapi *dlapi.DLAPI) *Server {
	s := &Server{
		router:  http.NewServeMux(),
		dlapi:   dlapi,
		port:    ":8080",
		devMode: os.Getenv("DEV_MODE") == "true",
	}

	s.setupRoutes()
	return s
}

// Start initializes and starts the server
func (s *Server) Start() error {
	log.Printf("Server starting on port %s...", s.port)
	return http.ListenAndServe(s.port, s.router)
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
