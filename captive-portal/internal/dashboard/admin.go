package handlers

import (
	"captive-portal/internal/integrations"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (h *PortalHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Check if user is admin
	session, err := h.getSession(r)
	if err != nil || !session.User.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Get active users from MikroTik
	activeUsers, err := h.mikrotik.GetActiveUsers()
	if err != nil {
		http.Error(w, "Failed to get active users", http.StatusInternalServerError)
		return
	}
	
	// Get system statistics
	stats := map[string]interface{}{
		"active_users": len(activeUsers),
		"total_users":  h.getTotalUsers(),
		"uptime":       time.Since(h.startTime).String(),
	}
	
	// Render admin dashboard
	data := map[string]interface{}{
		"ActiveUsers": activeUsers,
		"Stats":       stats,
		"Branding":    h.cfg.Branding,
	}
	
	h.renderTemplate(w, "admin_dashboard.html", data)
}

func (h *PortalHandler) APIHandler(w http.ResponseWriter, r *http.Request) {
	// Extract API endpoint from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid API request", http.StatusBadRequest)
		return
	}
	
	endpoint := pathParts[2]
	
	switch endpoint {
	case "users":
		h.handleUsersAPI(w, r)
	case "stats":
		h.handleStatsAPI(w, r)
	case "vouchers":
		h.handleVouchersAPI(w, r)
	default:
		http.Error(w, "Unknown API endpoint", http.StatusNotFound)
	}
}

func (h *PortalHandler) handleUsersAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := h.getSession(r)
	if err != nil || !session.User.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	switch r.Method {
	case "GET":
		// Get user list
		users := h.getUsers()
		json.NewEncoder(w).Encode(users)
	case "POST":
		// Create new user
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		if err := h.createUser(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *PortalHandler) handleStatsAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := h.getSession(r)
	if err != nil || !session.User.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get time range from query parameters
	timeRange := r.URL.Query().Get("range")
	if timeRange == "" {
		timeRange = "24h"
	}
	
	stats := h.getStatistics(timeRange)
	json.NewEncoder(w).Encode(stats)
}

func (h *PortalHandler) handleVouchersAPI(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	session, err := h.getSession(r)
	if err != nil || !session.User.IsAdmin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	switch r.Method {
	case "GET":
		// Get voucher list
		vouchers := h.getVouchers()
		json.NewEncoder(w).Encode(vouchers)
	case "POST":
		// Generate new vouchers
		var request struct {
			Count    int     `json:"count"`
			Duration string  `json:"duration"`
			Value    float64 `json:"value"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		
		vouchers, err := h.generateVouchers(request.Count, request.Duration, request.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(vouchers)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}