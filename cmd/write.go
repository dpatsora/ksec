/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		secretKey := args[1]
		newValue := args[2]

		kubeConf := viper.GetString("kubeconfig")
		config, err := clientcmd.BuildConfigFromFlags("", kubeConf)
		if err != nil {
			panic(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		secretsClient := clientset.CoreV1().Secrets(namespace)
		secret, err := secretsClient.Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		if oldValue, ok := secret.Data[secretKey]; ok {
			if string(oldValue) == newValue {
				fmt.Println("Current value match with the desired one")
				return
			}

			if !confirmOverwrite(string(oldValue), newValue) {
				return
			}
		}

		secret.Data[secretKey] = []byte(newValue)
		secretsClient.Update(context.TODO(), secret, metav1.UpdateOptions{})
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
