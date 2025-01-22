package dlapi

import (
	"encoding/json"
	"fmt"
	"sort"

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

	for _, format := range details.Formats {
		if format.Filesize > 0 {
			validFormats = append(validFormats, format)
		}
	}

	sort.Slice(validFormats, func(i, j int) bool {
		return validFormats[i].Filesize > validFormats[j].Filesize
	})

	details.Formats = validFormats
	return nil
}
