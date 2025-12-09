package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

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

// isFzfAvailable checks if fzf is installed and available in PATH
func isFzfAvailable() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// selectSecretWithFzf launches fzf with the provided secret names and returns the selected secret
func selectSecretWithFzf(secretNames []string) (string, error) {
	if len(secretNames) == 0 {
		return "", errors.New("no secrets found in namespace")
	}

	// Prepare input for fzf
	input := bytes.NewBufferString("")
	for _, name := range secretNames {
		input.WriteString(name + "\n")
	}

	// Execute fzf
	cmd := exec.Command("fzf", "--height=60%", "--reverse", "--prompt=Select secret: ")
	cmd.Stdin = input
	cmd.Stderr = nil // Suppress fzf stderr output

	output, err := cmd.Output()
	if err != nil {
		// User cancelled or fzf failed
		return "", errors.New("secret selection cancelled")
	}

	// Trim newline from output
	selected := string(bytes.TrimSpace(output))
	if selected == "" {
		return "", errors.New("no secret selected")
	}

	return selected, nil
}

func writeSecretToFile(secretData map[string][]byte, file *os.File) error {
	for k, v := range secretData {
		_, err := file.WriteString(fmt.Sprintf("%s: %s\n", k, string(v)))
		if err != nil {
			return err
		}
	}
	return nil
}
