package k8s

import (
	"k8s.io/client-go/tools/clientcmd"
)

func GetCurrentNamespace(kubeconfigFlag string) string {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfigFlag != "" {
		loadingRules.ExplicitPath = kubeconfigFlag
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{},
	)

	ns, _, err := kubeConfig.Namespace()
	if err != nil || ns == "" {
		return "default"
	}

	return ns
}
