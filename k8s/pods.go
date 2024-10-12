package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPodsForDeployment(ctx context.Context, client *kubernetes.Clientset, namespace string, deploymentName string) ([]v1.Pod, error) {
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", deploymentName),
	}

	podList, err := client.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}
