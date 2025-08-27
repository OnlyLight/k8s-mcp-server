package server

import (
	"context"
	"log"
	"onlylight/k8s-mcp-server/internal/config"
	"onlylight/k8s-mcp-server/internal/logging"
	"onlylight/k8s-mcp-server/pkg/k8s"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logging.NewLogger("info", "text")

	// Initialize Kubernetes client
	k8sClient, err := k8s.NewClient(cfg.K8s.ConfigPath, logger.Logger)
	if err != nil {
		logger.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Test Kubernetes connection
	ctx := context.Background()
	if err := k8sClient.HealthCheck(ctx); err != nil {
		logger.Fatalf("Kubernetes health check failed: %v", err)
	}
	logger.Info("Kubernetes connection established successfully")
}
