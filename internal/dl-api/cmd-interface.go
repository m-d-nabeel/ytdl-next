package dlapi

import (
	"fmt"
	"log"
	"os/exec"
)

const outputStdout = "-"

// Define base arguments as a slice instead of a single string
var baseArgs = []string{
	"--no-playlist",
	"--no-warnings",
	"--no-cookies",
}

func GetVideoInfoCmd(url string) *exec.Cmd {
	args := append([]string{
		"-j",
		"--write-thumbnail",
		"--skip-download",
		"--extractor-args", "youtube:player_client=web", // Faster extraction
	}, baseArgs...)

	return exec.Command("yt-dlp", append(args, url)...)
}

func GetMediaByFormatID(url, formatID string) *exec.Cmd {
	args := append([]string{
		"-f", formatID,
		"--progress-template", "",
		"--throttled-rate", "100K",
	}, baseArgs...)

	args = append(args, "-o", outputStdout)
	return exec.Command("yt-dlp", append(args, url)...)
}

func GetMediaByFormatIDS(url, audioFormatID, videoFormatID string) *exec.Cmd {
	log.Printf("FormatID: %v+%v\n", audioFormatID, videoFormatID)

	args := append([]string{
		"-f", fmt.Sprintf("%v+%v", audioFormatID, videoFormatID),
		"--concurrent-fragments", "4", // Parallel download
		"--audio-multistreams", // Allow multiple audio streams
		"--video-multistreams", // Allow multiple video streams
	}, baseArgs...)

	args = append(args, "-o", outputStdout)
	return exec.Command("yt-dlp", append(args, url)...)
}
