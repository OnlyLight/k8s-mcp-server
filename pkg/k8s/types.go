package k8s

import "time"

// PodInfo represents essential pod information.
type PodInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Status    string            `json:"status"`
	Phase     string            `json:"phase"`
	Node      string            `json:"node"`      // shows which cluster node is running the pod
	Labels    map[string]string `json:"labels"`    // provide metadata for grouping and selection
	CreatedAt time.Time         `json:"createdAt"` // tracks pod age
	Restarts  int32             `json:"restarts"`  // indicates stability issues
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
	Protocol   string `json:"protocol"`
}

// ServiceInfo represents essential service information.
type ServiceInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"`      // shows service exposure (ClusterIP, NodePort, LoadBalancer)
	ClusterIP string            `json:"clusterIP"` // provides internal cluster address
	Pods      []ServicePort     `json:"ports"`     // array defines all exposed endpoints
	Labels    map[string]string `json:"labels"`
	CreatedAt time.Time         `json:"createdAt"`
}

// DeploymentInfo represents essential deployment information.
type DeploymentInfo struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	TotalReplicas   int32             `json:"totalReplicas"`
	ReadyReplicas   int32             `json:"readyReplicas"`
	UpdatedReplicas int32             `json:"updatedReplicas"`
	Labels          map[string]string `json:"labels"`
	CreatedAt       time.Time         `json:"createdAt"`
	Strategy        string            `json:"strategy"` // indicates deployment approach (RollingUpdate vs Recreate)
}

// NamespaceInfo represents essential namespace information.
type NamespaceInfo struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Labels    map[string]string `json:"labels"`
	CreatedAt time.Time         `json:"createdAt"`
}
