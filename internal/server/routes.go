package server

import (
	"log"
	"net/http"

	"github.com/m-d-nabeel/ytdl-web/internal/embedfs"
)

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/health", s.handleHealth())

	// API routes
	s.router.HandleFunc("/api/yt/download", s.handleYTDownload)
	s.router.HandleFunc("/api/yt/info", s.handleYTInfo)

	// Serve static files
	var staticHandler http.Handler

	if s.devMode {
		// In development mode, serve from filesystem
		log.Println("Running in development mode, serving static files from website/dist")
		staticHandler = http.FileServer(http.Dir("website/dist"))
	} else {
		// In production mode, use embedded files
		log.Println("Running in production mode, serving embedded static files")
		staticHandler = http.FileServer(embedfs.GetHTTPFileSystem())
	}

	// Handle all other routes with the static file server
	s.router.Handle("/", staticHandler)
}
