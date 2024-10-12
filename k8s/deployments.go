package k8s

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Deployments(ctx context.Context, client *kubernetes.Clientset, namespace string) ([]v1.Deployment, error) {
	deploymentList, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return deploymentList.Items, nil
}
