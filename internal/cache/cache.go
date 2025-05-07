package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/m-d-nabeel/ytdl-web/internal/types"
)

type Cache struct {
	Path string
	db   *badger.DB
}

// NewCache creates a new Cache instance with an initialized Badger DB
func NewCache(path string) (*Cache, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil // Disable Badger's default logger

	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &Cache{
		Path: path,
		db:   db,
	}, nil
}

// Close closes the Badger DB
func (c *Cache) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (*types.YTMediaInfo, error) {
	var mediaInfo types.YTMediaInfo
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &mediaInfo)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &mediaInfo, nil
}

// Set stores an item in the cache
func (c *Cache) Set(key string, mediaInfo types.YTMediaInfo) error {
	value, err := json.Marshal(mediaInfo)
	if err != nil {
		return err
	}

	return c.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), value).WithTTL(24 * time.Hour) // Set TTL to 24 hours
		return txn.SetEntry(entry)
	})
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// LoadCache is kept for backward compatibility but does nothing
// as Badger DB doesn't need to be loaded from disk
func (c *Cache) LoadCache() error {
	log.Println("Using Badger DB for caching, no need to load cache")
	return nil
}

// SaveCache is kept for backward compatibility but does nothing
// as Badger DB automatically persists data
func (c *Cache) SaveCache() error {
	log.Println("Using Badger DB for caching, no need to save cache")
	return nil
}
