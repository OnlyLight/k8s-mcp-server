package mcp

import (
	"context"
	"fmt"
	"strings"

	"onlylight/k8s-mcp-server/internal/config"
	"onlylight/k8s-mcp-server/internal/logging"
	"onlylight/k8s-mcp-server/pkg/k8s"
	"onlylight/k8s-mcp-server/pkg/types"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the MCP server
type Server struct {
	config    *config.Config
	k8sClient *k8s.Client
	logger    *logging.Logger
	mcpServer *server.MCPServer
	formatter *ResourceFormatter
}

// NewServer creates a new MCP server instance with proper MCP protocol implementation
func NewServer(cfg *config.Config, k8sClient *k8s.Client) *Server {
	logger := logging.NewLogger("info", "text")

	// Create MCP server
	mcpServer := server.NewMCPServer("k8s-mcp-server", "1.0.0", server.WithResourceCapabilities(true, true))

	s := &Server{
		config:    cfg,
		k8sClient: k8sClient,
		logger:    logger,
		mcpServer: mcpServer,
		formatter: NewResourceFormatter(),
	}

	// Register MCP resources
	s.registerResources()

	return s
}

// Start starts the MCP server with stdio transport
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting Kubernetes MCP Server")

	// Use the convenient ServeStdio function
	if err := server.ServeStdio(s.mcpServer); err != nil {
		s.logger.Errorf("MCP server error: %v", err)
		return fmt.Errorf("MCP server failed: %w", err)
	}

	s.logger.Info("MCP Server stopped")
	return nil
}

// registerResources sets up the MCP resources and their handlers
func (s *Server) registerResources() {
	ctx := context.Background()

	// Register Pod resource
	services, err := s.k8sClient.ListServices(ctx, "")
	if err != nil {
		s.logger.Errorf("Failed to list pods: %v", err)
	} else {
		count := 0
		for _, service := range services {
			if count >= 5 { // limit to 5 pods for demo purposes
				break
			}

			resource := mcp.Resource{
				URI:         fmt.Sprintf("k8s://service/%s/%s", service.Namespace, service.Name),
				Name:        fmt.Sprintf("Service: %s/%s", service.Namespace, service.Name),
				Description: fmt.Sprintf("Kubernetes Service in namespace %s (Type: %s)", service.Namespace, service.Type),
				MIMEType:    "application/json",
			}

			s.mcpServer.AddResource(resource, s.handleResourceRead)
		}
	}
}

func (s *Server) handleResourceRead(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI
	s.logger.Infof("Handling read_resource request for URI: %s", uri)

	if !strings.HasPrefix(uri, "k8s://") {
		return nil, fmt.Errorf("invalid URI format. Expected k8s://<resource-type>/<namespace>/<name>, got: %s", uri)
	}

	// Parse URI: k8s://<resource-type>/<namespace>/<name>
	parts := strings.Split(strings.TrimPrefix(uri, "k8s://"), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid URI format. Expected k8s://<resource-type>/<namespace>/<name>, got %d parts", len(parts))
	}

	resourceType, namespace, name := parts[0], parts[1], parts[2]

	var resourceTypeEnum types.K8sResourceType
	switch resourceType {
	case "pod":
		resourceTypeEnum = types.ResourceTypePod
	case "service":
		resourceTypeEnum = types.ResourceTypeService
	case "deployment":
		resourceTypeEnum = types.ResourceTypeDeployment
	default:
		return nil, fmt.Errorf("unsupported resource type: %s. Supported types: pod, service, deployment", resourceType)
	}

	content, err := s.k8sClient.GetResource(ctx, &types.ResourceIdentifier{
		Type:      resourceTypeEnum,
		Namespace: namespace,
		Name:      name,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get resource %s: %w", uri, err)
	}

	// Format the content using AI-optimized formatters
	var formattedContent string
	var mimeType string

	switch resourceType {
	case "pod":
		formattedContent, err = s.formatter.FormatPodForAI(content)
		if err != nil {
			s.logger.Errorf("Failed to format pod data: %v", err)
			// Fall back to raw JSON
			formattedContent = content
			mimeType = "application/json"
		} else {
			mimeType = "text/markdown"
		}

	case "service":
		formattedContent, err = s.formatter.FormatServiceForAI(content)
		if err != nil {
			s.logger.Errorf("Failed to format service data: %v", err)
			// Fall back to raw JSON
			formattedContent = content
			mimeType = "application/json"
		} else {
			mimeType = "text/markdown"
		}

	case "deployment":
		formattedContent, err = s.formatter.FormatDeploymentForAI(content)
		if err != nil {
			s.logger.Errorf("Failed to format deployment data: %v", err)
			// Fall back to raw JSON
			formattedContent = content
			mimeType = "application/json"
		} else {
			mimeType = "text/markdown"
		}

	default:
		// For unsupported types, return raw JSON
		formattedContent = content
		mimeType = "application/json"
	}

	// Return the formatted resource contents
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      uri,
			MIMEType: mimeType,
			Text:     formattedContent,
		},
	}, nil
}
