package server

import (
	"fmt"
	"net/http"
)

// handleHealth returns a health check handler
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is healthy")
	}
}
