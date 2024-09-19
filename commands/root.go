package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "Une application CLI pour gérer les buckets et objets",
}

// Execute exécute la commande racine
func Execute() error {
	return rootCmd.Execute()
}
