package cmd

import (
	"github.com/spf13/cobra"
	"first/internal/app/install"

)
var installCmd = &cobra.Command{
	Use: "install",
	Short: "install and init model",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		install.Install()
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}