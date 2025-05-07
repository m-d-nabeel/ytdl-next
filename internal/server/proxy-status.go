package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	dlapi "github.com/m-d-nabeel/ytdl-web/internal/dl-api"
)

// ProxyStatus represents information about proxy configuration and status
type ProxyStatus struct {
	TotalProxies      int      `json:"total_proxies"`
	ActiveProxies     int      `json:"active_proxies"`
	LastRequestTime   string   `json:"last_request_time,omitempty"`
	GlobalRateGap     string   `json:"global_rate_gap"`
	MinRequestGap     string   `json:"min_request_gap"`
	ProxyTypes        []string `json:"proxy_types,omitempty"`
	UserAgentCount    int      `json:"user_agent_count"`
	RateLimitingInfo  string   `json:"rate_limiting_info"`
	RecentErrorsCount int      `json:"recent_errors_count"`
	Healthy           bool     `json:"healthy"`
}

// handleProxyStatus returns information about the current proxy configuration
func (s *Server) handleProxyStatus(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for admin authorization if needed
	// This is a simple check - you might want to implement proper auth
	authKey := r.URL.Query().Get("auth_key")
	requiresAuthKey := len(authKey) > 0 // true if auth_key is provided

	// Get proxy manager from dlapi package
	pm := dlapi.GetProxyManager()

	// Get basic public status (non-sensitive info)
	status := ProxyStatus{
		TotalProxies:      pm.GetTotalProxies(),
		ActiveProxies:     pm.GetActiveProxies(),
		LastRequestTime:   pm.GetLastRequestTime().Format(time.RFC3339),
		GlobalRateGap:     pm.GetGlobalRateGap().String(),
		MinRequestGap:     pm.GetMinRequestGap().String(),
		UserAgentCount:    pm.GetUserAgentCount(),
		RateLimitingInfo:  "Active",
		RecentErrorsCount: 0, // Placeholder for future error tracking
		Healthy:           true,
	}

	// Add sensitive information only if authorized
	if requiresAuthKey {
		// For now, we're not implementing real auth, just demonstrating the concept
		// In a real application, you should validate the auth key securely
		proxyTypes := pm.GetProxyTypes()
		status.ProxyTypes = proxyTypes

		// More detailed info could be added here
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding proxy status: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
