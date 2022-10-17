package cmd

import (
	"errors"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubernetesClient() (*kubernetes.Clientset, error) {
	kubeConf := viper.GetString("kubeconfig")
	if kubeConf == "" {
		return nil, errors.New("Kubeconfig is not provided")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConf)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}
