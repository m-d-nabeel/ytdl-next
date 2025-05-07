package dlapi

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

const outputStdout = "-"

// Define base arguments as a slice instead of a single string
var baseArgs = []string{
	"--no-playlist",
	"--no-warnings",
	"--no-cookies",
	"--extractor-retries", "3",
	"--socket-timeout", "30",
	"--force-ipv4",                 // Use IPv4 only to avoid IPv6-related issues
	"--geo-verification-proxy", "", // Will be set by proxy manager
	"--xff", "default", // Send X-Forwarded-For header to mask real IP
}

// Anti-blocking arguments - these help avoid detection
var antiBlockingArgs = []string{
	"--sleep-requests", "1", // Sleep 1 second between requests
	"--sleep-interval", "5", // Sleep 5 seconds between videos in playlist (not used but good practice)
	"--max-sleep-interval", "30", // Maximum sleep is 30 seconds
	"--sleep-subtitles", "1", // 1 second sleep before subtitle downloads
	"--user-agent", "", // Will be set dynamically
}

// Create global proxy manager instance
var (
	proxyManager *ProxyManager
	pmOnce       sync.Once
)

// GetProxyManager initializes a singleton ProxyManager
func GetProxyManager() *ProxyManager {
	pmOnce.Do(func() {
		// Initialize with empty proxy list - will be configured from server
		proxyManager = NewProxyManager(nil, nil)
	})
	return proxyManager
}

// GetVideoInfoCmd returns a command for fetching video information
func GetVideoInfoCmd(url string) *exec.Cmd {
	// Get a proxy and user agent from the proxy manager
	proxy, userAgent := GetProxyManager().GetNextProxy()

	args := append([]string{
		"-j",
		"--write-thumbnail",
		"--skip-download",
	}, baseArgs...)

	// Apply anti-blocking settings
	args = append(args, antiBlockingArgs...)

	// Set the user agent if provided
	if userAgent != "" {
		for i, arg := range args {
			if arg == "--user-agent" {
				args[i+1] = userAgent
				break
			}
		}
		if !contains(args, "--user-agent") {
			args = append(args, "--user-agent", userAgent)
		}
	}

	// Set the proxy if provided
	if proxy != "" {
		for i, arg := range args {
			if arg == "--geo-verification-proxy" {
				args[i+1] = proxy
				break
			}
		}
		if !contains(args, "--geo-verification-proxy") {
			args = append(args, "--geo-verification-proxy", proxy)
		}

		// Also use the proxy for all network operations
		args = append(args, "--proxy", proxy)
	}

	log.Printf("Running info request via proxy: %s", maskProxyCredentials(proxy))

	return exec.Command("yt-dlp", append(args, url)...)
}

// GetMediaByFormatID returns a command for downloading media with a specific format ID
func GetMediaByFormatID(url, formatID string) *exec.Cmd {
	// Get a proxy and user agent from the proxy manager
	proxy, userAgent := GetProxyManager().GetNextProxy()

	args := append([]string{
		"-f", formatID,
		"--progress-template", "",
		"--throttled-rate", "100K",
	}, baseArgs...)

	// Apply anti-blocking settings
	args = append(args, antiBlockingArgs...)

	// Set the user agent if provided
	if userAgent != "" {
		for i, arg := range args {
			if arg == "--user-agent" {
				args[i+1] = userAgent
				break
			}
		}
		if !contains(args, "--user-agent") {
			args = append(args, "--user-agent", userAgent)
		}
	}

	// Set the proxy if provided
	if proxy != "" {
		for i, arg := range args {
			if arg == "--geo-verification-proxy" {
				args[i+1] = proxy
				break
			}
		}
		if !contains(args, "--geo-verification-proxy") {
			args = append(args, "--geo-verification-proxy", proxy)
		}

		// Also use the proxy for all network operations
		args = append(args, "--proxy", proxy)
	}

	args = append(args, "-o", outputStdout)
	log.Printf("Running download request via proxy: %s", maskProxyCredentials(proxy))

	return exec.Command("yt-dlp", append(args, url)...)
}

// GetMediaByFormatIDS returns a command for downloading media with separate audio and video format IDs
func GetMediaByFormatIDS(url, audioFormatID, videoFormatID string) *exec.Cmd {
	log.Printf("FormatID: %v+%v\n", audioFormatID, videoFormatID)

	// Get a proxy and user agent from the proxy manager
	proxy, userAgent := GetProxyManager().GetNextProxy()

	args := append([]string{
		"-f", fmt.Sprintf("%v+%v", audioFormatID, videoFormatID),
		"--concurrent-fragments", "4", // Parallel download
		"--audio-multistreams", // Allow multiple audio streams
		"--video-multistreams", // Allow multiple video streams
	}, baseArgs...)

	// Apply anti-blocking settings
	args = append(args, antiBlockingArgs...)

	// Set the user agent if provided
	if userAgent != "" {
		for i, arg := range args {
			if arg == "--user-agent" {
				args[i+1] = userAgent
				break
			}
		}
		if !contains(args, "--user-agent") {
			args = append(args, "--user-agent", userAgent)
		}
	}

	// Set the proxy if provided
	if proxy != "" {
		for i, arg := range args {
			if arg == "--geo-verification-proxy" {
				args[i+1] = proxy
				break
			}
		}
		if !contains(args, "--geo-verification-proxy") {
			args = append(args, "--geo-verification-proxy", proxy)
		}

		// Also use the proxy for all network operations
		args = append(args, "--proxy", proxy)
	}

	args = append(args, "-o", outputStdout)
	log.Printf("Running combined download request via proxy: %s", maskProxyCredentials(proxy))

	return exec.Command("yt-dlp", append(args, url)...)
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// maskProxyCredentials hides the username/password in proxy URLs for logging
func maskProxyCredentials(proxy string) string {
	if proxy == "" {
		return "none"
	}

	// Simple masking for logging purposes
	// This is a basic implementation, could be improved for edge cases
	masked := proxy
	for i := 0; i < len(proxy); i++ {
		if proxy[i] == '@' && i > 0 {
			// Find the protocol marker
			protocolEnd := 0
			for j := 0; j < i; j++ {
				if j+2 < len(proxy) && proxy[j:j+3] == "://" {
					protocolEnd = j + 3
					break
				}
			}

			if protocolEnd > 0 && protocolEnd < i {
				// Replace the credentials part with ****
				masked = proxy[:protocolEnd] + "****" + proxy[i:]
			}
			break
		}
	}

	return masked
}
