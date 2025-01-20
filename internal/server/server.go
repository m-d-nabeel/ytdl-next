package server

import (
	"log"
	"net/http"
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
	log.Printf("Server starting on port %s...", s.port)
	return http.ListenAndServe(s.port, nil)
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
