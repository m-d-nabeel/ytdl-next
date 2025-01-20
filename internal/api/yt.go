package api

import (
	"fmt"

	"github.com/m-d-nabeel/ytdl-web/internal/cache"
	"github.com/m-d-nabeel/ytdl-web/internal/types"
)

type API struct {
	Cache *cache.Cache
}

func NewAPI(path string) *API {
	api := &API{
		Cache: &cache.Cache{
			Path: path,
		},
	}
	api.Cache.LoadCache()
	return api
}

func printMediaDetails(details types.YTMediaInfo) {
	// Print general media details
	fmt.Printf("media ID: %s\n", details.ID)
	fmt.Printf("Title: %s\n", details.Title)
	fmt.Printf("Uploader: %s\n", details.Uploader)
	fmt.Printf("Duration: %d seconds\n", details.Duration)
	fmt.Printf("View Count: %d \n", details.Duration)
	fmt.Printf("Like Count: %d \n", details.LikeCount)

	// Print available quality options
	fmt.Println("\nAvailable Formats:")
	for _, format := range details.Formats {
		if format.Filesize > 0 {
			fmt.Printf("- Format ID: %s, Resolution: %s, Extension: %s, Filesize: %d bytes, Note: %s\n",
				format.FormatID, format.Resolution, format.Ext, format.Filesize, format.Note)
		}
	}
}
