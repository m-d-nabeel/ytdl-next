package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/m-d-nabeel/ytdl-web/internal/types"
)

func (a *API) getYTMediaInfo(mediaUrl string) (types.YTMediaInfo, error) {
	details, ok := a.CachedData[mediaUrl]
	if ok {
		log.Println("Cache Hit")
		return details, nil
	}

	log.Println("Cache Miss")
	cmd := GetYTMediaInfoCmd()
	cmd.Args = append(cmd.Args, mediaUrl)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return types.YTMediaInfo{}, fmt.Errorf("failed to run yt-dlp: %w", err)
	}

	err = parseYTMediaInfo(output, &details)
	if err != nil {
		log.Fatalf("Error running yt-dlp: %v", err)
	}
	a.CachedData[mediaUrl] = details
	return details, nil
}

func parseYTMediaInfo(output []byte, details *types.YTMediaInfo) error {
	if err := json.Unmarshal(output, &details); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if details.ID == "" {
		return fmt.Errorf("missing required field: ID")
	}

	validFormats := make([]types.YTFormat, 0, len(details.Formats))
	for _, format := range details.Formats {
		if format.Filesize > 0 {
			validFormats = append(validFormats, format)
		}
	}

	details.Formats = validFormats
	return nil
}
