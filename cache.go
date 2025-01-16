package main

import (
	"encoding/json"
	"log"
	"os"
)

type Cache struct {
	Path string
	Data map[string]YTMediaInfo
}

func (c *Cache) loadCache() error {
	if c.Data == nil {
		c.Data = make(map[string]YTMediaInfo)
	}

	file, err := os.Open(c.Path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Cache file does not exist, starting with an empty cache.")
			return nil
		}
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("File closing failed\n%v", err.Error())
		}
	}()

	err = json.NewDecoder(file).Decode(&c.Data)
	if err != nil {
		log.Printf("Error decoding cache: %v", err)
	}
	log.Printf("Cache loaded")
	return err
}

func (c *Cache) saveCache() error {
	file, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("File closing failed\n%v", err.Error())
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c.Data)
	return err
}
