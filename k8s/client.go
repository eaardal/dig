package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"log"
)

// Client creates a Kubernetes clientset based on the specified context and namespace.
func Client(kubeContext string, namespace string) (*kubernetes.Clientset, error) {
	// Create a loading rules object to load the kubeconfig
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	// Specify overrides for context and namespace
	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: kubeContext,
		Context: clientcmdapi.Context{
			Namespace: namespace, // You can set the namespace here, or pass it in separately.
		},
	}

	// Create the client config
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		log.Fatalf("failed to load Kubernetes config: %v", err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to initialize Kubernetes clientset: %v", err)
	}

	return clientset, nil
}

//func Client() (*kubernetes.Clientset, error) {
//	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), nil).ClientConfig()
//	if err != nil {
//		log.Fatalf("failed to load Kubernetes config: %v", err)
//	}
//
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		log.Fatalf("failed to initialize Kubernetes clientset: %v", err)
//	}
//
//	return clientset, nil
//}
