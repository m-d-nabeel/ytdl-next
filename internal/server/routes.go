package server

import "net/http"

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/health", s.handleHealth())

	http.Handle("/", http.FileServer(http.Dir("website/dist")))

	http.HandleFunc("/api/yt/download", s.handleYTDownload)
	http.HandleFunc("/api/yt/info", s.handleYTInfo)
}
