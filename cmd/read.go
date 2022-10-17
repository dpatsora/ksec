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
	Use:     "read",
	Aliases: []string{"r"},
	Short:   "Retrieve secret data in human-readable format",
	Long: `Retrieve secret data in human-readable format

To retrieve "db-pass" secret data, located in "core" namespace, command will be:
ksec read db-pass -n core
`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		clientSet, err := getKubernetesClient()
		if err != nil {
			return err
		}

		secretName := args[0]
		secretsClient := clientSet.CoreV1().Secrets(namespace)
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
