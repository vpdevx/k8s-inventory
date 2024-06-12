package k8sconfig

import (
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	// This is required to auth to cloud providers (i.e. GKE)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Kube struct {
	Client kubernetes.Interface
}

var kubeClient *Kube
var once sync.Once

// GetConfigInstance returns a Pluto Kubernetes interface based on the current configuration
func GetConfigInstance(kubeContext string) (*Kube, error) {
	var err error
	var client kubernetes.Interface
	var kubeConfig *rest.Config

	kubeConfig, err = GetConfig(kubeContext)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		if kubeClient == nil {
			client, err = GetKubeClient(kubeConfig)

			kubeClient = &Kube{
				Client: client,
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

// GetConfig returns the current kube config with a specific context
func GetConfig(kubeContext string) (*rest.Config, error) {
	if kubeContext != "" {
		klog.V(3).Infof("using kube context: %s", kubeContext)
	}

	kubeConfig, err := config.GetConfigWithContext(kubeContext)
	if err != nil {
		return nil, err
	}
	return kubeConfig, nil
}

// GetKubeClient returns a Kubernetes.Interface based on the current configuration
func GetKubeClient(kubeConfig *rest.Config) (kubernetes.Interface, error) {
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
