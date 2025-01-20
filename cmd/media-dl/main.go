package main

import (
	"log"

	"github.com/m-d-nabeel/ytdl-web/internal/api"
	"github.com/m-d-nabeel/ytdl-web/internal/server"
)

func main() {
	api := api.NewAPI(".cache")
	defer func() {
		err := api.Cache.SaveCache()
		if err != nil {
			log.Println("Error saving cache")
		}
	}()

	srv := server.NewServer(api)

	err := srv.Start()
	if err != nil {
		log.Println(err.Error())
	}
}
