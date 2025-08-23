package main

import (
	"captive-portal/internal/config"
	"captive-portal/internal/handlers"
	"captive-portal/internal/integrations"
	"log"
	"net/http"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize MikroTik integration
	mikrotik, err := integrations.NewMikroTikIntegration(cfg.MikroTik)
	if err != nil {
		log.Fatal("Failed to initialize MikroTik:", err)
	}
	defer mikrotik.Close()

	// Initialize PHPNuxBill integration
	nuxbill := integrations.NewPHPNuxBillIntegration(cfg.PHPNuxBill)

	// Initialize handlers
	handler := handlers.NewPortalHandler(cfg, mikrotik, nuxbill)

	// Setup routes
	http.HandleFunc("/", handler.RedirectToPortal)
	http.HandleFunc("/login", handler.ShowLogin)
	http.HandleFunc("/auth", handler.Authenticate)
	http.HandleFunc("/voucher", handler.VoucherAuth)
	http.HandleFunc("/social/", handler.SocialAuth)
	http.HandleFunc("/terms", handler.AcceptTerms)
	http.HandleFunc("/portal", handler.ShowPortal)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/admin/", handler.AdminDashboard)
	http.HandleFunc("/api/", handler.APIHandler)
	
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	server := &http.Server{
		Addr:         cfg.Server.Address,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Server starting on %s", cfg.Server.Address)
	log.Fatal(server.ListenAndServe())
}