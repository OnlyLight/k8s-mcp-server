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
			summary.WriteString(fmt.Sprintf("**Created At**: %s (%s ago)\n", t.Format("2006-01-02 15:04:05"), formatDuration(time.Since(t))))
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

// // FormatDeploymentForAI creates an AI-optimized view of deployment information
// func (f *ResourceFormatter) FormatDeploymentForAI(deploymentData string) (string, error) {}

// // FormatServiceForAI creates an AI-optimized view of service information
// func (f *ResourceFormatter) FormatServiceForAI(serviceData string) (string, error) {}

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
