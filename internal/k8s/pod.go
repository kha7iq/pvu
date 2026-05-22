package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kha7iq/pvu/internal/models"
)

func GetPod(ctx context.Context, clientset kubernetes.Interface, namespace, podName string) (models.PodListItem, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return models.PodListItem{}, fmt.Errorf("get pod %q in namespace %q: %w", podName, namespace, err)
	}

	return models.PodListItem{
		Name:      pod.Name,
		Namespace: pod.Namespace,
		NodeName:  pod.Spec.NodeName,
		Phase:     string(pod.Status.Phase),
	}, nil
}
