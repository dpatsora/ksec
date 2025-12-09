package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var backupFlag bool

// writeCmd represents the write command
var editCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e"},
	Short:   "Edit secret data in your default editor",
	Long: `Write key/value pair to secret data

	To edit "db-pass" secret data, located in "core" namespace, command will be:
	ksec edit db-pass -n core

	If secret name is not provided and fzf is installed, you can select a secret interactively:
	ksec edit -n core

	WARNING: This will overwrite the existing secret data with the content you provide in the editor.
`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if secret name is provided
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

		// create a temp file and store secret data in it

		file, err := os.CreateTemp("", "ksec-edit-*.yaml")
		if err != nil {
			return err
		}
		// 2. Schedule the cleanup (Crucial!)
		defer func() {
			// Attempt to remove the file
			if rErr := os.Remove(file.Name()); rErr != nil {
				log.Printf("Warning: Could not remove temporary file %s: %v", file.Name(), rErr)
			}
		}()

		err = writeSecretToFile(secret.Data, file)
		if err != nil {
			return err
		}

		// close the file
		file.Close()

		// Get file modification time before opening editor
		fileInfoBefore, err := os.Stat(file.Name())
		if err != nil {
			return err
		}
		modTimeBefore := fileInfoBefore.ModTime()

		// open the file in editor
		err = openEditor(file.Name())
		if err != nil {
			return err
		}

		// Get file modification time after closing editor
		fileInfoAfter, err := os.Stat(file.Name())
		if err != nil {
			return err
		}
		modTimeAfter := fileInfoAfter.ModTime()

		// Check if file was modified
		if modTimeBefore.Equal(modTimeAfter) {
			fmt.Println("No changes detected. Secret not updated.")
			return nil
		}

		// Backup secret if flag is set
		if backupFlag {
			backupPath, err := backupSecret(secretName, namespace, secret.Data)
			if err != nil {
				return fmt.Errorf("backup failed: %v", err)
			}
			fmt.Printf("Backup created: %s\n", backupPath)
		}

		// read the file content and update the secret data
		updatedData, err := os.ReadFile(file.Name())
		if err != nil {
			return err
		}
		newSecretData := make(map[string][]byte)
		err = parseSecretData(updatedData, newSecretData)
		if err != nil {
			return err
		}

		fmt.Printf("Updating secret %s in namespace %s\n", secretName, namespace)
		secret.Data = newSecretData
		_, err = secretsClient.Update(context.TODO(), secret, metav1.UpdateOptions{})

		return err

	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolVarP(&backupFlag, "backup", "b", false, "backup secret before updating")
}

func openEditor(fileName string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // default to vim if EDITOR is not secret
	}

	cmd := exec.Command(editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func parseSecretData(fileContent []byte, secretData map[string][]byte) error {
	lines := bytes.Split(fileContent, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line in secret data: %s", line)
		}
		key := string(bytes.TrimSpace(parts[0]))
		value := string(bytes.TrimSpace(parts[1]))
		secretData[key] = []byte(value)
	}
	return nil
}

func backupSecret(secretName, namespace string, secretData map[string][]byte) (string, error) {
	// Create backup filename with timestamp in current directory
	timestamp := time.Now().Format("20060102-150405")
	backupFileName := fmt.Sprintf("ksec-backup-%s-%s-%s.yaml", secretName, namespace, timestamp)

	// Create backup file in current directory
	backupFile, err := os.Create(backupFileName)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %v", err)
	}
	defer backupFile.Close()

	err = writeSecretToFile(secretData, backupFile)
	if err != nil {
		return "", fmt.Errorf("failed to write backup: %v", err)
	}

	return backupFileName, nil
}
