package dlapi

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

// ProxyConfig represents configuration for a proxy server
type ProxyConfig struct {
	URL       string // Format: protocol://user:pass@host:port or protocol://host:port
	UserAgent string // Custom user agent to use with this proxy
	Weight    int    // Weight for random selection (higher = more likely to be chosen)
	LastUsed  time.Time
}

// ProxyManager handles proxy rotation and rate limiting
type ProxyManager struct {
	proxies        []ProxyConfig
	userAgents     []string
	mutex          sync.Mutex
	minRequestGap  time.Duration // Minimum time between requests to the same proxy
	lastRequestAt  time.Time     // Time of the last request to any proxy
	globalRateGap  time.Duration // Minimum time between any requests (global rate limiting)
	defaultProxy   string        // Default proxy to use if no proxies are available or configured
	defaultTimeout time.Duration // Default timeout for proxy connections
}

// NewProxyManager creates a new ProxyManager
func NewProxyManager(proxies []ProxyConfig, userAgents []string) *ProxyManager {
	// Set default user agents if none provided
	if len(userAgents) == 0 {
		userAgents = defaultUserAgents()
	}

	// Initialize last used time for each proxy to ensure they're initially available
	now := time.Now().Add(-24 * time.Hour) // Set last used to 24 hours ago
	for i := range proxies {
		proxies[i].LastUsed = now
	}

	return &ProxyManager{
		proxies:        proxies,
		userAgents:     userAgents,
		minRequestGap:  2 * time.Minute,  // Don't use the same proxy more than once every 2 minutes
		globalRateGap:  10 * time.Second, // Don't send requests faster than one every 10 seconds
		defaultTimeout: 30 * time.Second,
	}
}

// GetNextProxy returns the next available proxy and user agent
func (pm *ProxyManager) GetNextProxy() (string, string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Apply global rate limiting
	since := time.Since(pm.lastRequestAt)
	if since < pm.globalRateGap {
		sleepTime := pm.globalRateGap - since
		time.Sleep(sleepTime)
	}

	// Update last request time
	pm.lastRequestAt = time.Now()

	// If no proxies configured, return empty string (will use direct connection)
	if len(pm.proxies) == 0 {
		return pm.defaultProxy, randomUserAgent(pm.userAgents)
	}

	// Find available proxy
	now := time.Now()
	var availableProxies []int
	for i, proxy := range pm.proxies {
		if now.Sub(proxy.LastUsed) >= pm.minRequestGap {
			// Add proxy index to available proxies based on its weight
			weight := proxy.Weight
			if weight <= 0 {
				weight = 1
			}
			for j := 0; j < weight; j++ {
				availableProxies = append(availableProxies, i)
			}
		}
	}

	// If no proxies are available, wait for one to become available
	if len(availableProxies) == 0 {
		// Find the proxy that will be available the soonest
		earliestTime := now.Add(24 * time.Hour)
		earliestIndex := 0

		for i, proxy := range pm.proxies {
			availableAt := proxy.LastUsed.Add(pm.minRequestGap)
			if availableAt.Before(earliestTime) {
				earliestTime = availableAt
				earliestIndex = i
			}
		}

		// Wait until that proxy becomes available
		waitTime := earliestTime.Sub(now)
		time.Sleep(waitTime)

		// Update the proxy's last used time
		pm.proxies[earliestIndex].LastUsed = time.Now()

		// Use the user agent specified for this proxy or a random one
		userAgent := pm.proxies[earliestIndex].UserAgent
		if userAgent == "" {
			userAgent = randomUserAgent(pm.userAgents)
		}

		return pm.proxies[earliestIndex].URL, userAgent
	}

	// Select a random proxy from available ones
	selectedIdx := availableProxies[rand.Intn(len(availableProxies))]
	pm.proxies[selectedIdx].LastUsed = now

	// Use the user agent specified for this proxy or a random one
	userAgent := pm.proxies[selectedIdx].UserAgent
	if userAgent == "" {
		userAgent = randomUserAgent(pm.userAgents)
	}

	return pm.proxies[selectedIdx].URL, userAgent
}

// SetDefaultProxy sets the default proxy to use when no proxies are configured
func (pm *ProxyManager) SetDefaultProxy(proxy string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.defaultProxy = proxy
}

// SetMinRequestGap sets the minimum time between requests to the same proxy
func (pm *ProxyManager) SetMinRequestGap(duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.minRequestGap = duration
}

// SetGlobalRateGap sets the minimum time between any requests
func (pm *ProxyManager) SetGlobalRateGap(duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.globalRateGap = duration
}

// AddProxy adds a new proxy to the manager
func (pm *ProxyManager) AddProxy(proxy ProxyConfig) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// Set last used to long ago to make it immediately available
	proxy.LastUsed = time.Now().Add(-24 * time.Hour)
	pm.proxies = append(pm.proxies, proxy)
}

// RemoveProxy removes a proxy by URL
func (pm *ProxyManager) RemoveProxy(url string) bool {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for i, proxy := range pm.proxies {
		if proxy.URL == url {
			// Remove the proxy by swapping with the last element and truncating
			pm.proxies[i] = pm.proxies[len(pm.proxies)-1]
			pm.proxies = pm.proxies[:len(pm.proxies)-1]
			return true
		}
	}

	return false
}

// GetTotalProxies returns the total number of configured proxies
func (pm *ProxyManager) GetTotalProxies() int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return len(pm.proxies)
}

// GetActiveProxies returns the number of currently available proxies
func (pm *ProxyManager) GetActiveProxies() int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	now := time.Now()
	count := 0
	for _, proxy := range pm.proxies {
		if now.Sub(proxy.LastUsed) >= pm.minRequestGap {
			count++
		}
	}
	return count
}

// GetLastRequestTime returns the time of the last request
func (pm *ProxyManager) GetLastRequestTime() time.Time {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.lastRequestAt
}

// GetGlobalRateGap returns the minimum time between any requests
func (pm *ProxyManager) GetGlobalRateGap() time.Duration {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.globalRateGap
}

// GetMinRequestGap returns the minimum time between requests to the same proxy
func (pm *ProxyManager) GetMinRequestGap() time.Duration {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.minRequestGap
}

// GetUserAgentCount returns the number of configured user agents
func (pm *ProxyManager) GetUserAgentCount() int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return len(pm.userAgents)
}

// GetProxyTypes returns the types of proxies in use (http, https, socks5, etc.)
func (pm *ProxyManager) GetProxyTypes() []string {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	typeMap := make(map[string]bool)
	for _, proxy := range pm.proxies {
		// Extract protocol from URL (e.g., "http://example.com" -> "http")
		if len(proxy.URL) > 0 {
			parts := strings.SplitN(proxy.URL, "://", 2)
			if len(parts) > 0 {
				typeMap[parts[0]] = true
			}
		}
	}

	// Convert map keys to slice
	types := make([]string, 0, len(typeMap))
	for t := range typeMap {
		types = append(types, t)
	}

	return types
}

// randomUserAgent returns a random user agent from the provided list
func randomUserAgent(userAgents []string) string {
	if len(userAgents) == 0 {
		// Default user agent if none provided
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	}
	return userAgents[rand.Intn(len(userAgents))]
}

// defaultUserAgents returns a list of common user agents
func defaultUserAgents() []string {
	return []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36 Edg/92.0.902.55",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	}
}
