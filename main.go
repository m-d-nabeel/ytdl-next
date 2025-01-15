package main

import (
	"fmt"
	"log"
)

type API struct{}

func main() {
	tempUrl := "https://www.youtube.com/shorts/l_4OinnwnS4"
	a := API{}
	details, err := a.getYTMediaInfo(tempUrl)
	printMediaDetails(details)
	if err != nil {
		log.Fatalf("ERROR: %v", err.Error())
	}
}

func printMediaDetails(details YTMediaInfo) {
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
