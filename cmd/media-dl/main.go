package main

import (
	"log"

	"github.com/m-d-nabeel/ytdl-web/internal/api"
	"github.com/m-d-nabeel/ytdl-web/internal/cache"
	"github.com/m-d-nabeel/ytdl-web/internal/server"
	"github.com/m-d-nabeel/ytdl-web/internal/types"
)

func main() {
	cache := cache.Cache{
		Path: "./.cache",
		Data: make(map[string]types.YTMediaInfo),
	}
	cache.LoadCache()
	// tempUrl := "https://www.youtube.com/shorts/l_4OinnwnS4"
	api := api.API{
		CachedData: make(map[string]types.YTMediaInfo),
	}
	api.CachedData = cache.Data
	defer func() {
		cache.Data = api.CachedData
		err := cache.SaveCache()
		if err != nil {
			log.Println("Error saving cache")
		}
	}()

	srv := server.NewServer(server.DefaultConfig(), &api)
	err := srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}
