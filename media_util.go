package main

import (
	"encoding/json"
	"log"
)

func (a *API) getYTMediaInfo(mediaUrl string) (YTMediaInfo, error) {
	details, ok := a.CachedData[mediaUrl]
	if ok {
		log.Println("Cache Hit")
		return details, nil
	}

	log.Println("Cache Miss")
	cmd := a.GetYTMediaInfoCmd()
	cmd.Args = append(cmd.Args, mediaUrl)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running yt-dlp: %v", err)
	}

	err = json.Unmarshal(output, &details)
	if err != nil {
		log.Fatalf("Error parsing JSON output: %v", err)
	}

	a.CachedData[mediaUrl] = details
	return details, nil
}
