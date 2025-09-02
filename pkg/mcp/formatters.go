package mcp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ResourceFormatter struct{}

func NewResourceFormatter() *ResourceFormatter {
	return &ResourceFormatter{}
}

// FormatPodForAI creates an AI-optimized representation of Pod information
func (f *ResourceFormatter) FormatPodForAI(podData string) (string, error) {
	var pod map[string]interface{}
	if err := json.Unmarshal([]byte(podData), &pod); err != nil {
		return "", err
	}

	summary := &strings.Builder{}
	summary.WriteString("# Pod Summary:\n\n")

	// Basic information
	summary.WriteString(fmt.Sprintf("**Name**: %s\n", pod["name"]))
	summary.WriteString(fmt.Sprintf("**Namespace**: %s\n", pod["namespace"]))
	summary.WriteString(fmt.Sprintf("**Status**: %s\n", pod["status"]))
	summary.WriteString(fmt.Sprintf("**Node**: %s\n", pod["node"]))

	if restarts, ok := pod["restarts"].(float64); ok && restarts > 0 {
		summary.WriteString(fmt.Sprintf("**⚠️ Restarts**: %.0f\n", restarts))
	}

	// Creation time
	if createdAt, ok := pod["createdAt"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			summary.WriteString(fmt.Sprintf("**Created At**: %s\n", formatDuration(time.Since(t))))
		}
	}

	summary.WriteString("\n## Containers:\n\n")

	// Container Information
	if containers, ok := pod["containers"].([]interface{}); ok {
		for _, c := range containers {
			if c, ok := c.(map[string]interface{}); ok {
				name := c["name"].(string)
				image := c["image"].(string)
				ready := c["ready"].(bool)
				state := c["state"].(string)

				status := "🟢 Ready"
				if !ready {
					status = "🔴 Not Ready"
				}

				summary.WriteString(fmt.Sprintf("- **%s**: %s\n", name, status))
				summary.WriteString(fmt.Sprintf("  - Image: %s\n", image))
				summary.WriteString(fmt.Sprintf("  - State: %s\n", state))

				if restarts, ok := c["restartCount"].(float64); ok && restarts > 0 {
					summary.WriteString(fmt.Sprintf("  - ⚠️ Restarts: %.0f\n", restarts))
				}
			}
		}
	}

	// Condition
	if conditions, ok := pod["conditions"].([]interface{}); ok && len(conditions) > 0 {
		summary.WriteString("\n## Conditions:\n")
		for _, cond := range conditions {
			if condStr, ok := cond.(string); ok {
				summary.WriteString(fmt.Sprintf("- %s\n", condStr))
			}
		}
	}

	// Labels
	if labels, ok := pod["labels"].(map[string]interface{}); ok && len(labels) > 0 {
		summary.WriteString("\n## Labels:\n")
		for k, v := range labels {
			summary.WriteString(fmt.Sprintf("- %s: %s\n", k, v))
		}
	}

	summary.WriteString("\n---\n")
	summary.WriteString("*Use this information to understand the pod's current state and troubleshoot any issues.*")

	return summary.String(), nil
}

// FormatDeploymentForAI creates an AI-optimized view of deployment information
func (f *ResourceFormatter) FormatDeploymentForAI(deploymentData string) (string, error) {
	var deployment map[string]interface{}
	if err := json.Unmarshal([]byte(deploymentData), &deployment); err != nil {
		return "", err
	}

	summary := &strings.Builder{}
	summary.WriteString("# Deployment Summary:\n\n")

	// Basic information
	summary.WriteString(fmt.Sprintf("**Name**: %s\n", deployment["name"]))
	summary.WriteString(fmt.Sprintf("**Namespace**: %s\n", deployment["namespace"]))
	summary.WriteString(fmt.Sprintf("**Strategy**: %s\n", deployment["strategy"]))

	// Replicas Status
	total := deployment["totalReplicas"].(float64)
	ready := deployment["readyReplicas"].(float64)
	updated := deployment["updatedReplicas"].(float64)

	healthStatus := "🟢 Healthy"
	if ready < total {
		healthStatus = "🟠 Scaling"
	}
	if ready == 0 {
		healthStatus = "🔴 Unhealthy"
	}

	summary.WriteString(fmt.Sprintf("**Status**: %s\n", healthStatus))
	summary.WriteString(fmt.Sprintf("**Replicas**: %d total, %d ready, %d updated\n", int(total), int(ready), int(updated)))

	// Progress indicator
	if total > 0 {
		percentage := (ready / total) * 100
		summary.WriteString(fmt.Sprintf("**Progress**: %.0f%%\n", percentage))
	}

	// Creation time
	if createdAt, ok := deployment["createdAt"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			age := time.Since(t)
			summary.WriteString(fmt.Sprintf("**Created At**: %s\n", formatDuration(age)))
		}
	}

	// Seletor
	if selector, ok := deployment["selector"].(map[string]interface{}); ok && len(selector) > 0 {
		summary.WriteString("\n## Selector:\n")
		for key, value := range selector {
			summary.WriteString(fmt.Sprintf("- %s: %s\n", key, value))
		}
	}

	// Conditions
	if conditions, ok := deployment["conditions"].([]interface{}); ok && len(conditions) > 0 {
		summary.WriteString("\n## Conditions:\n")
		for _, cond := range conditions {
			if condStr, ok := cond.(string); ok {
				summary.WriteString(fmt.Sprintf("- %s\n", condStr))
			}
		}
	}

	// Recommendations
	summary.WriteString("\n## AI Assitant Notes\n\n")
	if ready < total {
		summary.WriteString("⚠️ **Action Needed**: Some replicas are not ready. Check pod status and logs.\n")
	}
	summary.WriteString("\n---\n")
	if ready == 0 {
		summary.WriteString("⚠️ **Action Needed**: Some replicas are not ready. Check pod status and logs.\n")
	}
	if ready == total {
		summary.WriteString("✅ **Status**: Deployment is healthy and all replicas are ready.\n")
	}

	return summary.String(), nil
}

// FormatServiceForAI creates an AI-optimized view of service information
func (f *ResourceFormatter) FormatServiceForAI(serviceData string) (string, error) {
	var service map[string]interface{}
	if err := json.Unmarshal([]byte(serviceData), &service); err != nil {
		return "", err
	}

	summary := &strings.Builder{}
	summary.WriteString("# Service Summary:\n\n")

	// Basic information
	summary.WriteString(fmt.Sprintf("**Name**: %s\n", service["name"]))
	summary.WriteString(fmt.Sprintf("**Namespace**: %s\n", service["namespace"]))
	summary.WriteString(fmt.Sprintf("**Type**: %s\n", service["type"]))

	// Cluster IP
	if clusterIP, ok := service["clusterIP"].(string); ok && clusterIP != "" {
		summary.WriteString(fmt.Sprintf("**Cluster IP**: %s\n", clusterIP))
	}

	// Ports
	if ports, ok := service["ports"].([]interface{}); ok && len(ports) > 0 {
		summary.WriteString("\n## Ports:\n")
		for _, p := range ports {
			if portMap, ok := p.(map[string]interface{}); ok {
				name := ""
				if n, ok := portMap["name"].(string); ok && n != "" {
					name = fmt.Sprintf("%s: ", n)
				}
				summary.WriteString(fmt.Sprintf("- %sPort %d -> Target Port %d (%s)\n", name, int(portMap["port"].(float64)), int(portMap["targetPort"].(float64)), portMap["protocol"].(string)))
			}
		}
	}

	// Selector
	if selector, ok := service["selector"].(map[string]interface{}); ok && len(selector) > 0 {
		summary.WriteString("\n## Selector:\n")
		summary.WriteString("This service routes traffic to pods with these labels:\n")
		for key, value := range selector {
			summary.WriteString(fmt.Sprintf("- %s: %s\n", key, value))
		}
	}

	// Service Type specific info
	serviceType := service["type"].(string)
	summary.WriteString("\n## Connectivity:\n")

	switch serviceType {
	case "ClusterIP":
		summary.WriteString("🔒 **Internal Access Only**: This service is only accessible within the cluster.\n")
	case "NodePort":
		summary.WriteString("🌐 **External Access**: This service is accessible from outside the cluster via node IPs.\n")
	case "LoadBalancer":
		summary.WriteString("⚖️ **Load Balancer**: This service has an external load balancer.\n")
	case "ExternalName":
		summary.WriteString("🔗 **External Name**: This service maps to an external DNS name.\n")
	}

	summary.WriteString("\n---\n")
	summary.WriteString("*Use this information to understand the service's configuration and connectivity.*")

	return summary.String(), nil
}

// Helper function to format duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	days := d.Hours() / 24
	return fmt.Sprintf("%.1fd", days)
}
