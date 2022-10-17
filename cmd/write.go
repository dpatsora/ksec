/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:     "write",
	Aliases: []string{"w"},
	Short:   "Write key/value pair to secret data",
	Long: `Write key/value pair to secret data

To add "USER_PASSWORD: admin123" to "db-pass" secret data, located in "core" namespace, command will be:
ksec write db-pass USER_PASSWORD admin123 -n core
`,
	Args: cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretName := args[0]
		secretKey := args[1]
		newValue := args[2]

		clientSet, err := getKubernetesClient()
		if err != nil {
			return err
		}

		secretsClient := clientSet.CoreV1().Secrets(namespace)
		secret, err := secretsClient.Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if oldValue, ok := secret.Data[secretKey]; ok {
			if string(oldValue) == newValue {
				fmt.Println("Current value match with the desired one")
				return nil
			}

			if !confirmOverwrite(string(oldValue), newValue) {
				return nil
			}
		}

		secret.Data[secretKey] = []byte(newValue)
		_, err = secretsClient.Update(context.TODO(), secret, metav1.UpdateOptions{})

		return err
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
}

func confirmOverwrite(oldValue, newValue string) bool {
	var input string

	fmt.Println("Current value:", oldValue)
	fmt.Println("New value:", newValue)
	fmt.Println()

	for {
		fmt.Printf("Do you want to continue with this operation? [y|n]: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			panic(err)
		}
		input = strings.ToLower(input)

		if input == "y" || input == "yes" {
			return true
		}

		if input == "n" || input == "no" {
			return false
		}

		fmt.Printf("Unrecognized input %s\n", input)
	}
}
