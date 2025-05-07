package dlapi

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/m-d-nabeel/ytdl-web/internal/types"
)

func ParseYTMediaInfo(output []byte, details *types.YTMediaInfo) error {
	if err := json.Unmarshal(output, &details); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if details.ID == "" {
		return fmt.Errorf("missing required field: ID")
	}

	validFormats := make([]types.YTFormat, 0, len(details.Formats))

	// Process formats to identify compatible ones and assign quality levels
	for i := range details.Formats {
		format := details.Formats[i]

		// Only include formats with filesize > 0
		if format.Filesize <= 0 {
			continue
		}

		// Mark compatibility based on codec
		format.IsCompatible = format.IsUniversallyCompatible()

		// Assign quality labels based on resolution and compatibility
		format.QualityLabel = assignQualityLabel(&format)

		validFormats = append(validFormats, format)
	}

	// Sort formats by filesize (descending)
	sort.Slice(validFormats, func(i, j int) bool {
		// First prioritize compatible formats
		if validFormats[i].IsCompatible != validFormats[j].IsCompatible {
			return validFormats[i].IsCompatible
		}

		// Then sort by filesize
		return validFormats[i].Filesize > validFormats[j].Filesize
	})

	details.Formats = validFormats
	return nil
}

// assignQualityLabel determines a human-readable quality level for a format
func assignQualityLabel(format *types.YTFormat) string {
	if !format.HasVideo() {
		// For audio-only formats
		return determineAudioQuality(format)
	}

	// Extract height from resolution (e.g. "1280x720" -> 720)
	height := 0
	if format.Resolution != "" && format.Resolution != "audio only" {
		parts := strings.Split(format.Resolution, "x")
		if len(parts) == 2 {
			if h, err := strconv.Atoi(parts[1]); err == nil {
				height = h
			}
		}
	}

	// Determine quality based on resolution
	var quality string
	switch {
	case height >= 2160:
		quality = "4K"
	case height >= 1440:
		quality = "2K"
	case height >= 1080:
		quality = "HD"
	case height >= 720:
		quality = "HD"
	case height >= 480:
		quality = "SD"
	case height > 0:
		quality = "Low"
	default:
		quality = "Unknown"
	}

	// Append compatibility indication
	if format.IsCompatible {
		quality += " (Compatible)"
	} else if format.IsVP9() {
		quality += " (VP9)"
	} else if format.IsAV1() {
		quality += " (AV1)"
	}

	return quality
}

// determineAudioQuality assigns a quality label to audio formats
func determineAudioQuality(format *types.YTFormat) string {
	// Check if the format note contains bitrate info
	note := strings.ToLower(format.Note)

	if strings.Contains(note, "high") || strings.Contains(note, "128") {
		if format.IsCompatible {
			return "High Quality Audio (Compatible)"
		}
		return "High Quality Audio"
	}

	if strings.Contains(note, "medium") {
		if format.IsCompatible {
			return "Medium Quality Audio (Compatible)"
		}
		return "Medium Quality Audio"
	}

	if strings.Contains(note, "low") {
		if format.IsCompatible {
			return "Low Quality Audio (Compatible)"
		}
		return "Low Quality Audio"
	}

	// Default audio quality label
	if format.IsCompatible {
		return "Audio (Compatible)"
	}
	return "Audio"
}
