package handlers

import (
	"captive-portal/internal/auth"
	"captive-portal/internal/integrations"
	"captive-portal/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AuthHandler struct {
	cfg       *config.Config
	mikrotik  *integrations.MikroTikIntegration
	nuxbill   *integrations.PHPNuxBillIntegration
	localAuth *auth.LocalAuthenticator
	radiusAuth *auth.RadiusAuthenticator
	socialAuth *auth.SocialAuthenticator
	voucherAuth *auth.VoucherAuthenticator
}

func NewAuthHandler(cfg *config.Config, mikrotik *integrations.MikroTikIntegration, nuxbill *integrations.PHPNuxBillIntegration) *AuthHandler {
	return &AuthHandler{
		cfg:       cfg,
		mikrotik:  mikrotik,
		nuxbill:   nuxbill,
		localAuth: auth.NewLocalAuthenticator(cfg.Database),
		radiusAuth: auth.NewRadiusAuthenticator(cfg.Radius),
		socialAuth: auth.NewSocialAuthenticator(cfg.SocialAuth),
		voucherAuth: auth.NewVoucherAuthenticator(nuxbill, mikrotik),
	}
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	username := r.FormValue("username")
	password := r.FormValue("password")
	
	// Try local authentication first
	user, err := h.localAuth.Authenticate(username, password)
	if err != nil {
		// Fall back to RADIUS authentication
		user, err = h.radiusAuth.Authenticate(username, password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}
	
	// Create session
	session, err := models.NewSession(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	})
	
	// Redirect to portal or terms acceptance if needed
	if h.cfg.Features.TermsConditionsEnabled && !user.TermsAccepted {
		http.Redirect(w, r, "/terms", http.StatusFound)
		return
	}
	
	http.Redirect(w, r, "/portal", http.StatusFound)
}

func (h *AuthHandler) VoucherAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	if !h.cfg.Features.VoucherEnabled {
		http.Error(w, "Voucher authentication is disabled", http.StatusForbidden)
		return
	}
	
	voucherCode := r.FormValue("voucher_code")
	
	user, err := h.voucherAuth.Authenticate(voucherCode)
	if err != nil {
		http.Error(w, "Invalid voucher", http.StatusUnauthorized)
		return
	}
	
	// Create session
	session, err := models.NewSession(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	})
	
	// Redirect to portal or terms acceptance if needed
	if h.cfg.Features.TermsConditionsEnabled && !user.TermsAccepted {
		http.Redirect(w, r, "/terms", http.StatusFound)
		return
	}
	
	http.Redirect(w, r, "/portal", http.StatusFound)
}

func (h *AuthHandler) SocialAuth(w http.ResponseWriter, r *http.Request) {
	// Extract provider from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid social auth request", http.StatusBadRequest)
		return
	}
	
	provider := pathParts[2]
	
	// Handle different social auth providers
	switch provider {
	case "facebook":
		h.handleFacebookAuth(w, r)
	case "google":
		h.handleGoogleAuth(w, r)
	default:
		http.Error(w, "Unsupported social provider", http.StatusBadRequest)
	}
}

func (h *AuthHandler) handleFacebookAuth(w http.ResponseWriter, r *http.Request) {
	// Implementation for Facebook OAuth2 flow
	// This would redirect to Facebook, then handle the callback
}

func (h *AuthHandler) handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	// Implementation for Google OAuth2 flow
	// This would redirect to Google, then handle the callback
}