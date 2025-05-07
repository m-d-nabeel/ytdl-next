package dlapi

import (
	"log"

	"github.com/m-d-nabeel/ytdl-web/internal/cache"
)

type DLAPI struct {
	Cache *cache.Cache
}

func NewDLAPI(path string) *DLAPI {
	// Initialize Badger DB cache
	cacheInstance, err := cache.NewCache(path)
	if err != nil {
		log.Printf("Failed to initialize cache: %v", err)
		log.Printf("Continuing without cache...")
		return &DLAPI{Cache: nil}
	}

	return &DLAPI{
		Cache: cacheInstance,
	}
}
