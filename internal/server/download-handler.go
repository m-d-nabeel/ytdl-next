package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	mediaURL := r.URL.Query().Get("url")
	if mediaURL == "" {
		http.Error(w, "URL parameter is required", http.StatusBadRequest)
		return
	}

	// Get video info first
	filename, err := s.getFilename(mediaURL)
	if err != nil {
		http.Error(w, "Failed to get video info", http.StatusInternalServerError)
		return
	}
	// Set headers for streaming download
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	// Create the download command
	cmd := exec.Command("yt-dlp",
		"-o", "-", // Output to stdout
		"--newline", // Force progress to new lines
		mediaURL)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Failed to create pipe", http.StatusInternalServerError)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, "Failed to create error pipe", http.StatusInternalServerError)
		return
	}

	if err = cmd.Start(); err != nil {
		http.Error(w, "Failed to start download", http.StatusInternalServerError)
		return
	}

	// Monitor progress in a goroutine
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("Progress: %s", scanner.Text())
		}
	}()

	// Stream the download directly to the client
	_, err = io.Copy(w, stdout)
	if err != nil {
		log.Printf("Error streaming: %v", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command error: %v", err)
	}
}

func (s *Server) getFilename(url string) (string, error) {
	cmd := exec.Command("yt-dlp",
		"--get-filename",
		"-o", "%(title)s.%(ext)s",
		url)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	filename := strings.TrimSpace(string(output))
	// Sanitize filename
	filename = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`<>:"/\|?*`, r) {
			return '-'
		}
		return r
	}, filename)

	return filename, nil
}
