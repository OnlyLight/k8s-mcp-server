# Building MCP server kubernetes go

## 1. Project Structure

```
k8s-mcp-server/
├── cmd/server/          # Main application entry point
├── pkg/
│   ├── mcp/            # MCP protocol implementation
│   ├── k8s/            # Kubernetes client wrapper
│   └── types/          # Shared type definitions
├── internal/
│   ├── config/         # Configuration management
│   └── logging/        # Logging setup
├── go.mod
├── go.sum
├── Makefile            # Build and test scripts
└── scripts/            # Helper scripts
```

### 1.1. Dependencies

```
# Real MCP Go library from mark3labs
go get github.com/mark3labs/mcp-go@v0.36.0

# Kubernetes client
go get k8s.io/client-go@v0.31.2
go get k8s.io/api@v0.31.2
go get k8s.io/apimachinery@v0.31.2

# Logging and utilities
go get github.com/sirupsen/logrus@v1.9.3
go get gopkg.in/yaml.v3@v3.0.1
```

### 1.2. Configuration Setup
```internal/config/config.go```
**What this code does:** Creates a configuration management system that defines how our MCP server will connect to Kubernetes and what settings it will use.

**Data Flow:** Environment variable → File reading → YAML parsing → Config struct → Return to caller

### 1.3. Logging Setup
```internal/logging/logger.go```

**What this code does:** Creates a sophisticated logging system that tracks both MCP protocol operations and Kubernetes API calls with structured logging and proper error handling.

**Data transformation:** String parameters → Structured log fields → Formatted output → stdout

## 2. Kubernetes

### 2.1. Kubernetes Types
```pkg/k8s/types.go```
**What this code does:** Defines simplified data structures that extract only the essential information from complex Kubernetes objects, making them easier for AI models to understand and process.

**Data transformation:** Complex Kubernetes API objects → Simplified structs → JSON serializable data → AI-friendly format

### 2.2. Kubernetes Client Implementation
```pkg/k8s/client.go```
**What this code does:**
- Implements the core Kubernetes client functionality with connection handling, health checking, and pod listing capabilities that transform complex Kubernetes API responses into our simplified data structures.

- Implement service and deployment listing functionality with data transformation similar to pod listing, extracting networking and scaling information respectively.

### 2.3. AI-Formatters
```pkg/mcp/formatters.go```
**What this code does:** Create formatters that transform raw Kubernetes JSON into AI-friendly markdown output.
