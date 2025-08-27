package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	clientset *kubernetes.Clientset
	logger    *logrus.Logger
}

func NewClient(configPath string, logger *logrus.Logger) (*Client, error) {
	config, err := buildConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return &Client{
		clientset: clientset,
		logger:    logger,
	}, nil
}

func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("Kubernetes cluster not reachable: %w", err)
	}

	return nil
}

func (c *Client) ListPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var podInfos []PodInfo
	for _, pod := range pods.Items {
		podInfos = append(podInfos, PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Phase:     string(pod.Status.Phase),
			Node:      pod.Spec.NodeName,
			Labels:    pod.Labels,
			CreatedAt: pod.CreationTimestamp.Time,
			Restarts:  getPodRestartCount(&pod),
		})
	}

	return podInfos, nil
}

func (c *Client) ListServices(ctx context.Context, namespace string) ([]ServiceInfo, error) {
	services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in namespace %s: %w", namespace, err)
	}

	var serviceInfos []ServiceInfo
	for _, svc := range services.Items {
		var ports []ServicePort
		for _, port := range svc.Spec.Ports {
			ports = append(ports, ServicePort{
				Name:       port.Name,
				Port:       port.Port,
				TargetPort: port.TargetPort.String(),
				Protocol:   string(port.Protocol),
			})
		}

		serviceInfos = append(serviceInfos, ServiceInfo{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     ports,
			Labels:    svc.Labels,
			CreatedAt: svc.CreationTimestamp.Time,
		})
	}

	return serviceInfos, nil
}

func (c *Client) ListDeployments(ctx context.Context, namespace string) ([]DeploymentInfo, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	var deploymentInfos []DeploymentInfo
	for _, deploy := range deployments.Items {
		strategy := "RollingUpdate"
		if deploy.Spec.Strategy.Type == appsv1.RecreateDeploymentStrategyType {
			strategy = "Recreate"
		}

		deploymentInfos = append(deploymentInfos, DeploymentInfo{
			Name:            deploy.Name,
			Namespace:       deploy.Namespace,
			TotalReplicas:   *deploy.Spec.Replicas,
			ReadyReplicas:   deploy.Status.ReadyReplicas,
			UpdatedReplicas: deploy.Status.UpdatedReplicas,
			Labels:          deploy.Labels,
			CreatedAt:       deploy.CreationTimestamp.Time,
			Strategy:        strategy,
		})
	}

	return deploymentInfos, nil
}

func buildConfig(configPath string) (*rest.Config, error) {
	// try in-cluster config first
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	// Fallback to kubeconfig file
	if configPath == "" {
		if home := homedir.HomeDir(); home != "" {
			configPath = filepath.Join(home, ".kube", "config")
		}
	}

	return clientcmd.BuildConfigFromFlags("", configPath)
}

func getPodRestartCount(pod *corev1.Pod) int32 {
	var restarts int32
	for _, status := range pod.Status.ContainerStatuses {
		restarts += status.RestartCount
	}
	return restarts
}

func getPodConditions(pod *corev1.Pod) []string {
	var conditions []string
	for _, condition := range pod.Status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			conditions = append(conditions, string(condition.Type))
		}
	}

	return conditions
}

func getDeploymentConditions(deployment *appsv1.Deployment) []string {
	var conditions []string
	for _, condition := range deployment.Status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			conditions = append(conditions, fmt.Sprintf("%s: %s", condition.Type, condition.Message))
		}
	}

	return conditions
}
