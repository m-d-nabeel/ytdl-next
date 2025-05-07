package main

import (
	"log"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
	"github.com/m-d-nabeel/ytdl-web/internal/server"
)

func main() {
	dlapi := dlapi.NewDLAPI(".badger-cache")
	// Close the Badger DB when the application exits
	if dlapi.Cache != nil {
		defer func() {
			log.Println("Closing Badger DB...")
			if err := dlapi.Cache.Close(); err != nil {
				log.Printf("Error closing cache: %v", err)
			}
		}()
	}

	srv := server.NewServer(dlapi)

	err := srv.Start()
	if err != nil {
		log.Println(err.Error())
	}
}
