package embedfs

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web
var webFS embed.FS

// GetHTTPFileSystem returns an http.FileSystem for the embedded files
func GetHTTPFileSystem() http.FileSystem {
	// Get the subdirectory
	subFS, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal("Failed to get web subdirectory from embedded filesystem:", err)
	}
	return http.FS(subFS)
}
