package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/m-d-nabeel/ytdl-web/internal/types"
)

type Cache struct {
	Path   string
	db     *badger.DB
	stopGC chan struct{} // Channel to signal garbage collection to stop
}

// NewCache creates a new Cache instance with an initialized Badger DB
func NewCache(path string) (*Cache, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil // Disable Badger's default logger

	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	cache := &Cache{
		Path:   path,
		db:     db,
		stopGC: make(chan struct{}),
	}

	// Start garbage collection in a goroutine
	go cache.runGarbageCollection()

	return cache, nil
}

// runGarbageCollection periodically runs garbage collection on the Badger DB
// to actually remove expired items from the database
func (c *Cache) runGarbageCollection() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Run garbage collection
			err := c.db.RunValueLogGC(0.5) // Run GC if we can reclaim at least 50% of a file
			if err != nil && err != badger.ErrNoRewrite {
				// It's normal to get ErrNoRewrite if there's not enough garbage to collect
				log.Printf("Badger value log GC: %v", err)
			}
		case <-c.stopGC:
			return // Stop the garbage collection routine
		}
	}
}

// Close closes the Badger DB
func (c *Cache) Close() error {
	// Signal the garbage collection routine to stop
	close(c.stopGC)

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
		entry := badger.NewEntry([]byte(key), value).WithTTL(5 * time.Minute) // Set TTL to 24 hours
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
