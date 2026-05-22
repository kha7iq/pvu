package k8s

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func BuildConfig(kubeconfigFlag string) (*rest.Config, error) {
	var cfg *rest.Config
	var err error

	if kubeconfigFlag != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfigFlag)
		if err != nil {
			return nil, fmt.Errorf("build kube config from flag %q: %w", kubeconfigFlag, err)
		}
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

			if os.Getenv("KUBECONFIG") == "" {
				home, homeErr := os.UserHomeDir()
				if homeErr != nil {
					return nil, fmt.Errorf("get user home dir: %w", homeErr)
				}
				loadingRules.ExplicitPath = filepath.Join(home, ".kube", "config")
			}

			overrides := &clientcmd.ConfigOverrides{}
			clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)

			cfg, err = clientConfig.ClientConfig()
			if err != nil {
				return nil, fmt.Errorf("load kube config: %w", err)
			}
		}
	}

	cfg.Timeout = 30 * time.Second
	cfg.Dial = (&net.Dialer{
		Timeout: 30 * time.Second,
	}).DialContext

	return cfg, nil
}

func NewClientset(kubeconfigFlag string) (*kubernetes.Clientset, error) {
	cfg, err := BuildConfig(kubeconfigFlag)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	return clientset, nil
}
