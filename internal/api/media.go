package api

import (
	"github.com/m-d-nabeel/ytdl-web/internal/cache"
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
