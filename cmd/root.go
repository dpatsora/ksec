/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	kubeconfig string
	namespace  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ksec",
	Short: "Ksec is a CLI tool to work with k8s secrets",
	Long: `Ksec is a CLI tool to work with k8s secrets

Kubeconfig can be passed by flag --kubeconfig or by "KUBECONFIG" env  (flag takes precedence if provided)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "resource k8s namespace")
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "path to k8s configuration file")
	err := viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	if err != nil {
		panic(err)
	}
}

// initConfig reads ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
