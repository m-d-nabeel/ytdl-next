package api

import (
	"os/exec"
)

// GetVideoInfoCmd returns a command to fetch video information with available formats
func GetVideoInfoCmd(url string) *exec.Cmd {
	return exec.Command("yt-dlp",
		"-j",                // Output as JSON
		"--write-thumbnail", // Get thumbnail URL/info
		"--skip-download",   // Don't download the video
		"--no-playlist",     // Don't process playlists
		url,
	)
}

// GetMediaByFormatID returns a command to download media with specific format ID
func GetMediaByFormatID(url, formatID string) *exec.Cmd {
	return exec.Command("yt-dlp",
		"-f", formatID,
		"-q",                      // Quiet mode
		"--progress-template", "", // Disable progress output
		"--no-playlist",
		"-o", "-", // Output to stdout
		url,
	)
}
