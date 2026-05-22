package app

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/kha7iq/pvu/internal/k8s"
	"github.com/kha7iq/pvu/internal/models"
	"github.com/kha7iq/pvu/internal/tui"
	"github.com/kha7iq/pvu/internal/ui"
)

func Run(opts models.Options) error {
	clientset, err := k8s.NewClientset(opts.Kubeconfig)
	if err != nil {
		return fmt.Errorf("create kubernetes client: %w", err)
	}

	namespace := opts.Namespace
	if namespace == "" {
		namespace = k8s.GetCurrentNamespace(opts.Kubeconfig)
	}

	if opts.PodName != "" {
		return runCLI(context.Background(), clientset, namespace, opts)
	}

	return runTUI(clientset, namespace)
}

func runCLI(ctx context.Context, clientset kubernetes.Interface, namespace string, opts models.Options) error {
	view, err := k8s.GetPodVolumeView(ctx, clientset, namespace, opts.PodName)
	if err != nil {
		return fmt.Errorf("load pod %q in namespace %q: %w", opts.PodName, namespace, err)
	}

	fmt.Println(ui.RenderPodVolumeView(view))
	return nil
}

func runTUI(clientset kubernetes.Interface, namespace string) error {
	if err := tui.Run(clientset, namespace); err != nil {
		return fmt.Errorf("run interactive view in namespace %q: %w", namespace, err)
	}
	return nil
}
