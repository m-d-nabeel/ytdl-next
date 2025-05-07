package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
)

type DownloadResponse struct {
	Status   string     `json:"status"`
	Error    string     `json:"error,omitempty"`
	URL      string     `json:"url,omitempty"`
	Reader   io.Reader  // Used for streaming download content
	ErrorCh  chan error // Channel for receiving download errors
	FileSize int64      // File size in bytes for Content-Length header
	CancelFn func()     // Function to call to cancel the download process
}

type Server struct {
	router          *http.ServeMux
	dlapi           *dlapi.DLAPI
	port            string
	devMode         bool
	downloadManager *DownloadManager
}

// NewServer creates a new server instance with the given config
func NewServer(dlapiInstance *dlapi.DLAPI) *Server {
	// Calculate a reasonable number of workers based on available CPU cores
	maxWorkers := min(runtime.NumCPU(), 4)

	// Initialize and load proxy configuration
	const proxyConfigFilename = "proxies.json"

	// Use package-level functions by using the package name, not the variable
	proxyConfig, err := dlapi.LoadProxyConfig(proxyConfigFilename)
	if err != nil {
		log.Printf("Warning: Error loading proxy configuration: %v", err)
		// Continue with default proxy settings
	} else if proxyConfig != nil && len(proxyConfig.Proxies) == 0 {
		// Create a default config file if none exists
		if err := dlapi.InitializeDefaultProxyConfig(proxyConfigFilename); err != nil {
			log.Printf("Warning: Could not create default proxy config: %v", err)
		} else {
			log.Printf("Created default proxy configuration. Please edit %s in your home directory to add your own proxies.", proxyConfigFilename)
		}
	}

	s := &Server{
		router:          http.NewServeMux(),
		dlapi:           dlapiInstance,
		port:            ":8080",
		devMode:         os.Getenv("DEV_MODE") == "true",
		downloadManager: NewDownloadManager(maxWorkers, 50),
	}

	s.setupRoutes()
	return s
}

// Start initializes and starts the server
func (s *Server) Start() error {
	log.Printf("Server starting on port %s...", s.port)

	// Set up graceful shutdown handler
	server := &http.Server{
		Addr:    s.port,
		Handler: s.router,
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)

	// Register for interrupt and termination signals
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine so it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	log.Printf("Server is ready to handle requests at http://localhost%s", s.port)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, shutting down gracefully...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Shut down the download manager gracefully
	s.downloadManager.Shutdown(30 * time.Second)

	log.Println("Server shutdown complete")
	return nil
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
