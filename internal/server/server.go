package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/m-d-nabeel/ytdl-web/internal/api"
)

// ServerConfig holds all server configuration options
type ServerConfig struct {
	Port            string
	MediaDirectory  string
	StaticDirectory string
	ChunkSize       int64
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	CertFile        string
	KeyFile         string
}

// DefaultConfig returns a ServerConfig with sensible defaults
func DefaultConfig() *ServerConfig {
	return &ServerConfig{
		Port:            ":8080",
		MediaDirectory:  "./media",
		StaticDirectory: "./static",
		ChunkSize:       1024 * 1024, // 1MB
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1MB
	}
}

// Server represents our HTTP server
type Server struct {
	config     *ServerConfig
	httpServer *http.Server
	router     *http.ServeMux
	wg         sync.WaitGroup
	api        *api.API
}

// NewServer creates a new server instance with the given config
func NewServer(config *ServerConfig, api *api.API) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	s := &Server{
		config: config,
		router: http.NewServeMux(),
		api:    api,
	}

	s.setupRoutes()
	return s
}

// setupRoutes initializes all server routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.HandleFunc("/health", s.handleHealth())

	// Media streaming
	s.router.HandleFunc("/stream/", s.handleStreamMedia())

	// Static files
	s.router.Handle("/static/",
		http.StripPrefix("/static/",
			s.cacheMiddleware(
				http.FileServer(http.Dir(s.config.StaticDirectory)))))
}

// Start initializes and starts the server
func (s *Server) Start() error {
	// Create required directories
	if err := s.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	s.httpServer = &http.Server{
		Addr:           s.config.Port,
		Handler:        s.router,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	// Start server in a goroutine
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		var err error
		if s.config.CertFile != "" && s.config.KeyFile != "" {
			err = s.httpServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
		} else {
			err = s.httpServer.ListenAndServe()
		}
		if err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s...", s.config.Port)
	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	s.wg.Wait()
	return nil
}

// createDirectories ensures required directories exist
func (s *Server) createDirectories() error {
	dirs := []string{s.config.MediaDirectory, s.config.StaticDirectory}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// handleHealth returns a health check handler
func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server is healthy")
	}
}

// handleStreamMedia returns a media streaming handler
func (s *Server) handleStreamMedia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Base(r.URL.Path)
		path := filepath.Join(s.config.MediaDirectory, filename)

		file, err := os.Open(path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			http.Error(w, "Failed to get file info", http.StatusInternalServerError)
			return
		}

		// Stream the file in chunks
		w.Header().Set("Content-Type", s.getContentType(filename))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

		buf := make([]byte, s.config.ChunkSize)
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading file: %v", err)
				return
			}
			if _, err := w.Write(buf[:n]); err != nil {
				log.Printf("Error writing response: %v", err)
				return
			}
		}
	}
}

// cacheMiddleware adds caching headers to static files
func (s *Server) cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	})
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
