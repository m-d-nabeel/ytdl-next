package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCache_loadCache(t *testing.T) {
	tests := []struct {
		name       string
		initial    string // Initial content of the cache file
		cachePath  string // Path for the cache file
		expected   map[string]interface{}
		setupError bool
		wantErr    bool
	}{
		{
			name:      "Valid cache file",
			initial:   `{"Path":"test.json","Data":{"key1":"value1","key2":42}}`,
			cachePath: "test_valid.json",
			expected:  map[string]interface{}{"key1": "value1", "key2": float64(42)},
			wantErr:   false,
		},
		{
			name:      "Empty cache file",
			initial:   `{}`,
			cachePath: "test_empty.json",
			expected:  map[string]interface{}{},
			wantErr:   false,
		},
		{
			name:      "Malformed JSON",
			initial:   `{"Path":"test.json","Data":{key1:"value1"}}`,
			cachePath: "test_malformed.json",
			wantErr:   true,
		},
		{
			name:      "Non-existent file",
			cachePath: "non_existent.json",
			wantErr:   false,
		},
		{
			name:       "File access error",
			cachePath:  "/root/protected.json",
			setupError: true,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Cache

			// Setup the environment
			if !tt.setupError {
				dir := t.TempDir()
				cacheFilePath := filepath.Join(dir, tt.cachePath)
				tt.cachePath = cacheFilePath
				if tt.initial != "" {
					if err := os.WriteFile(tt.cachePath, []byte(tt.initial), 0o644); err != nil {
						t.Fatalf("Failed to create test file: %v", err)
					}
				}
			}

			c.Path = tt.cachePath
			err := c.loadCache()

			if (err != nil) != tt.wantErr {
				t.Errorf("loadCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !tt.wantErr && tt.expected != nil {
				for key, value := range tt.expected {
					if c.Data[key] != value {
						t.Errorf("loadCache() expected %v, got %v", value, c.Data[key])
					}
				}
			}
		})
	}
}

func TestCache_saveCache(t *testing.T) {
	tests := []struct {
		name      string
		cachePath string // Path for the cache file
		data      map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "Valid save",
			cachePath: "test_valid_save.json",
			data:      map[string]interface{}{"key1": "value1", "key2": "42"},
			wantErr:   false,
		},
		{
			name:      "Empty data",
			cachePath: "test_empty_save.json",
			data:      map[string]interface{}{},
			wantErr:   false,
		},
		{
			name:      "Invalid path",
			cachePath: "/root/protected_save.json",
			data:      map[string]interface{}{"key": "value"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Cache
			dir := t.TempDir()
			cacheFilePath := filepath.Join(dir, tt.cachePath)
			tt.cachePath = cacheFilePath

			c.Path = tt.cachePath
			c.Data = tt.data

			err := c.saveCache()

			if (err != nil) != tt.wantErr {
				t.Errorf("saveCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the file exists
				if _, err := os.Stat(tt.cachePath); os.IsNotExist(err) {
					t.Fatalf("Cache file was not created: %v", tt.cachePath)
				}

				// Verify the file contents
				content, err := os.ReadFile(tt.cachePath)
				if err != nil {
					t.Fatalf("Failed to read cache file: %v", err)
				}
				t.Logf("Cache file content: %s", string(content))

				// Decode and compare
				var result Cache
				if err := json.Unmarshal(content, &result); err != nil {
					t.Fatalf("Failed to unmarshal cache file: %v", err)
				}
				t.Logf("Decoded data: %+v", result.Data)

				if !reflect.DeepEqual(result.Data, tt.data) {
					t.Errorf("saveCache() expected %v, got %v", tt.data, result.Data)
				}
			}
		})
	}
}
