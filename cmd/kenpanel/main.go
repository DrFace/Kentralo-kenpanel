package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DrFace/Kentralo-kenpanel/core/config"
	"github.com/DrFace/Kentralo-kenpanel/core/events"
	"github.com/DrFace/Kentralo-kenpanel/core/jobs"
	"github.com/DrFace/Kentralo-kenpanel/core/policy"
)

func main() {
	configPath := flag.String("config", "/etc/kenpanel/kenpanel.conf", "Path to KenPanel configuration file")
	flag.Parse()

	fmt.Printf("============================================================\n")
	fmt.Printf(" KenPanel Unified Control Plane v%s (%s)\n", config.Version, config.GitCommit)
	fmt.Printf(" One panel. Your servers. Your rules.\n")
	fmt.Printf("============================================================\n")

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("[WARN] Using default configuration: %v\n", err)
		cfg = config.DefaultConfig()
	}

	// Initialize core subsystems
	auditLedger := events.NewAuditLedger([]byte(cfg.JWTSecret))
	policyEngine := policy.NewEngine()
	jobQueue := jobs.NewQueue(100)

	_ = auditLedger
	_ = policyEngine
	_ = jobQueue

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","version":"%s"}`, config.Version)
	})

	mux.HandleFunc("/api/v1/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"%s","git_commit":"%s","build_date":"%s"}`,
			config.Version, config.GitCommit, config.BuildDate)
	})

	server := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("[INFO] KenPanel Control Plane listening on https://%s\n", cfg.ListenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] HTTP server failed: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down KenPanel Control Plane gracefully...")
}
