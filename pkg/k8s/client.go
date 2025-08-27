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
