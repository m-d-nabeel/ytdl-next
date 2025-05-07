package dlapi

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

// ProxyConfigFile represents the structure of a proxy configuration file
type ProxyConfigFile struct {
	Proxies           []ProxyConfig `json:"proxies"`
	UserAgents        []string      `json:"user_agents"`
	MinRequestGap     string        `json:"min_request_gap"`
	GlobalRateGap     string        `json:"global_rate_gap"`
	DefaultProxy      string        `json:"default_proxy"`
	RetryAttempts     int           `json:"retry_attempts"`
	ConnectionTimeout string        `json:"connection_timeout"`
}

// LoadProxyConfig loads proxy configuration from a JSON file
func LoadProxyConfig(filename string) (*ProxyConfigFile, error) {
	// Find config file in multiple locations
	var configPaths []string

	// First check working directory
	workingDir, err := os.Getwd()
	if err == nil {
		configPaths = append(configPaths, filepath.Join(workingDir, filename))
	}

	// Check home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPaths = append(configPaths, filepath.Join(homeDir, ".ytdl-next", filename))
	}

	// Check /etc directory for system-wide configuration
	configPaths = append(configPaths, filepath.Join("/etc", "ytdl-next", filename))

	// Try each path until we find a valid config file
	var data []byte
	var foundPath string

	for _, path := range configPaths {
		data, err = os.ReadFile(path)
		if err == nil {
			foundPath = path
			break
		}
	}

	// If no config file found, return empty config
	if data == nil {
		log.Printf("No proxy configuration file found. Checked: %v", configPaths)
		return &ProxyConfigFile{
			Proxies:           []ProxyConfig{},
			UserAgents:        []string{},
			MinRequestGap:     "2m",
			GlobalRateGap:     "10s",
			DefaultProxy:      "",
			RetryAttempts:     3,
			ConnectionTimeout: "30s",
		}, nil
	}

	log.Printf("Loading proxy configuration from %s", foundPath)

	var config ProxyConfigFile
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Configure the proxy manager with the loaded config
	pm := GetProxyManager()

	// Set custom user agents if provided
	if len(config.UserAgents) > 0 {
		pm.userAgents = config.UserAgents
	}

	// Add all configured proxies
	for _, proxy := range config.Proxies {
		pm.AddProxy(proxy)
	}

	// Set default proxy
	if config.DefaultProxy != "" {
		pm.SetDefaultProxy(config.DefaultProxy)
	}

	// Set time gaps if provided
	if config.MinRequestGap != "" {
		if duration, err := time.ParseDuration(config.MinRequestGap); err == nil {
			pm.SetMinRequestGap(duration)
		} else {
			log.Printf("Warning: Invalid min_request_gap format: %v", err)
		}
	}

	if config.GlobalRateGap != "" {
		if duration, err := time.ParseDuration(config.GlobalRateGap); err == nil {
			pm.SetGlobalRateGap(duration)
		} else {
			log.Printf("Warning: Invalid global_rate_gap format: %v", err)
		}
	}

	return &config, nil
}

// SaveProxyConfig saves the proxy configuration to a JSON file
func SaveProxyConfig(config *ProxyConfigFile, filename string) error {
	// Determine save location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".ytdl-next")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	configPath := filepath.Join(configDir, filename)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// InitializeDefaultProxyConfig creates a default proxy configuration file if none exists
func InitializeDefaultProxyConfig(filename string) error {
	// Check if config already exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".ytdl-next")
	configPath := filepath.Join(configDir, filename)

	// Don't overwrite existing config
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}

	// Create default config
	defaultConfig := ProxyConfigFile{
		Proxies: []ProxyConfig{
			// Add some free proxies as examples
			// Users should replace these with their own proxies
			{
				URL:       "http://example-proxy.com:8080",
				UserAgent: "",
				Weight:    1,
			},
		},
		UserAgents:        defaultUserAgents(),
		MinRequestGap:     "2m",
		GlobalRateGap:     "10s",
		DefaultProxy:      "",
		RetryAttempts:     3,
		ConnectionTimeout: "30s",
	}

	// Create directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	// Save the default config
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	log.Printf("Creating default proxy configuration at %s", configPath)
	return os.WriteFile(configPath, data, 0644)
}
