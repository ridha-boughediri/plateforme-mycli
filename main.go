package main

import (
	"log"
	"os"

	"github.com/ridha-boughediri/plateforme-mycli/commands"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "myS3"}

	// Ajout des commandes au rootCmd
	rootCmd.AddCommand(commands.CreateCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.DeleteCmd)
	rootCmd.AddCommand(commands.UploadCmd) // Ajoutez cette ligne pour la commande upload

	// Ajout des commandes supplémentaires
	rootCmd.AddCommand(commands.VersionCmd)
	rootCmd.AddCommand(commands.SyncCmd)
	rootCmd.AddCommand(commands.ShowCmd)
	rootCmd.AddCommand(commands.UnsyncCmd)

	// Exécution de la commande racine
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Erreur lors de l'exécution de la commande : %v", err)
		os.Exit(1)
	}
}
