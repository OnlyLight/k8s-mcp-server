package server

import (
	"context"
	"log"
	"onlylight/k8s-mcp-server/internal/config"
	"onlylight/k8s-mcp-server/internal/logging"
	"onlylight/k8s-mcp-server/pkg/k8s"
	"onlylight/k8s-mcp-server/pkg/mcp"
	"os"
	"os/signal"
	"syscall"
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

	// Create MCP server
	mcpServer := mcp.NewServer(cfg, k8sClient)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start MCP server
	serverErrChan := make(chan error, 1)
	go func() {
		if err := mcpServer.Start(ctx); err != nil {
			serverErrChan <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-sigChan:
		logger.Infof("Received signal: %v. Shutting down...", sig)
		cancel()
	case err := <-serverErrChan:
		logger.Errorf("MCP server error: %v", err)
		cancel()
	}

	logger.Info("Server exited gracefully")
}
