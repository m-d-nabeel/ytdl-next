package server

import "net/http"

// setupRoutes initializes all server routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth())

	// Serve static files from the public directory
	http.Handle("/", http.FileServer(http.Dir("website/dist")))

	// API endpoints
	http.HandleFunc("/api/yt/download", s.handleYTDownload)
	http.HandleFunc("/api/yt/info", s.handleYTInfo)
}
