package dlapi

import (
	"github.com/m-d-nabeel/ytdl-web/internal/cache"
)

type DLAPI struct {
	Cache *cache.Cache
}

func NewDLAPI(path string) *DLAPI {
	dlapi := &DLAPI{
		Cache: &cache.Cache{
			Path: path,
		},
	}
	dlapi.Cache.LoadCache()
	return dlapi
}
