package cmd

import (
	"first/internal/app/migrate"
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "upgrade model",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start migrating...")
		migrate.Migrate()
		fmt.Println("migrate successfully")
	},
}

func init() {
	migrateCmd.DisableSuggestions = true
	rootCmd.AddCommand(migrateCmd)
}