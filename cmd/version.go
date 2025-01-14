package cmd

import (
	"first/internal/defines"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the current app version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		version := `
App Version: %s
Model Version: %s
Go Version: 1.18.8
Mysql Version: 5.7
`
		fmt.Printf(version, defines.Version, defines.DBVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
