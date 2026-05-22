package k8s

import (
	"context"
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/kha7iq/pvu/internal/models"
)

func ListPods(ctx context.Context, clientset kubernetes.Interface, namespace string) ([]models.PodListItem, error) {
	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods in namespace %q: %w", namespace, err)
	}

	items := make([]models.PodListItem, 0, len(podList.Items))
	for _, pod := range podList.Items {
		items = append(items, models.PodListItem{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			NodeName:  pod.Spec.NodeName,
			Phase:     string(pod.Status.Phase),
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return items, nil
}
