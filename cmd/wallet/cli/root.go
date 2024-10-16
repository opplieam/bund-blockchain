package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	accountName string
	accountPath string
)

const (
	keyExtension = ".ecdsa"
)

func init() {
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&accountName, "account", "a", "private.ecdsa", "The account to use.")
	rootCmd.PersistentFlags().StringVarP(&accountPath, "account-path", "p", "conf/accounts/", "Path to the directory with private keys.")
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Simple wallet",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getPrivateKeyPath() string {
	if !strings.HasSuffix(accountName, keyExtension) {
		accountName += keyExtension
	}

	return filepath.Join(accountPath, accountName)
}
