/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:     "read [secret-name]",
	Aliases: []string{"r"},
	Short:   "Retrieve secret data in human-readable format",
	Long: `Retrieve secret data in human-readable format

To retrieve "db-pass" secret data, located in "core" namespace, command will be:
ksec read db-pass -n core

If secret name is not provided and fzf is installed, you can select a secret interactively:
ksec read -n core
`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clientSet, err := getKubernetesClient()
		if err != nil {
			return err
		}

		secretsClient := clientSet.CoreV1().Secrets(namespace)

		var secretName string
		if len(args) == 0 {
			// No secret name provided, try to use fzf
			if !isFzfAvailable() {
				return fmt.Errorf("secret name is required when fzf is not installed")
			}

			// List all secrets in the namespace
			secretList, err := secretsClient.List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return err
			}

			if len(secretList.Items) == 0 {
				return fmt.Errorf("no secrets found in namespace %s", namespace)
			}

			// Extract secret names
			secretNames := make([]string, len(secretList.Items))
			for i, secret := range secretList.Items {
				secretNames[i] = secret.Name
			}

			// Launch fzf for selection
			selectedSecret, err := selectSecretWithFzf(secretNames)
			if err != nil {
				return err
			}

			secretName = selectedSecret
		} else {
			secretName = args[0]
		}

		secret, err := secretsClient.Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		for k, v := range secret.Data {
			fmt.Printf("%s: %s\n", k, string(v))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
