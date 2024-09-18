package main

import (
	"log"
	"os"
	"plateforme-mycli/cmd"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "myS3"}

	// Ajout des commandes au rootCmd
	rootCmd.AddCommand(cmd.CreateCmd)
	rootCmd.AddCommand(cmd.ListCmd)
	rootCmd.AddCommand(cmd.DeleteCmd)
	rootCmd.AddCommand(cmd.UploadCmd) // Ajoutez cette ligne pour la commande upload

	// Exécution de la commande racine
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Erreur lors de l'exécution de la commande : %v", err)
		os.Exit(1)
	}
}
