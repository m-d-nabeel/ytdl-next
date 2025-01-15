package main

import (
	"encoding/json"
	"log"
)

func (a *API) getYTMediaInfo(mediaUrl string) (YTMediaInfo, error) {
	cmd := a.GetYTMediaCmd()
	cmd.Args = append(cmd.Args, mediaUrl)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error running yt-dlp: %v", err)
	}

	var details YTMediaInfo
	err = json.Unmarshal(output, &details)
	if err != nil {
		log.Fatalf("Error parsing JSON output: %v", err)
	}

	return details, nil
}
