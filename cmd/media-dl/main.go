package main

import (
	"log"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
	"github.com/m-d-nabeel/ytdl-web/internal/server"
)

func main() {
	dlapi := dlapi.NewDLAPI(".cache")
	defer func() {
		err := dlapi.Cache.SaveCache()
		if err != nil {
			log.Println("Error saving cache")
		}
	}()

	srv := server.NewServer(dlapi)

	err := srv.Start()
	if err != nil {
		log.Println(err.Error())
	}
}
