package main

import (
	"fmt"
	"os"
	"plateforme-mycli/commands"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "mycli"}

	// Ajoutez les commandes pour les buckets et les objets
	rootCmd.AddCommand(commands.BucketCmd)
	rootCmd.AddCommand(commands.ObjectCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
