/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	Run: func(cmd *cobra.Command, args []string) {
		kubeConf := viper.GetString("kubeconfig")
		config, err := clientcmd.BuildConfigFromFlags("", kubeConf)
		if err != nil {
			panic(err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		secretName := args[0]
		secretsClient := clientset.CoreV1().Secrets(namespace)
		secret, err := secretsClient.Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		for k, v := range secret.Data {
			fmt.Printf("%s: %s\n", k, string(v))
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
