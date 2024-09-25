package main

import (
	"fmt"
	"log"
	"os"
	"plateforme-mycli/commands"
	"plateforme-mycli/storage"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env : %v", err)
	}

	// Initialiser le client HTTP S3
	storage.InitS3Client()

	var rootCmd = &cobra.Command{Use: "clis3"}

	// Ajouter les commandes pour les buckets et les objets
	rootCmd.AddCommand(commands.BucketCmd)
	rootCmd.AddCommand(commands.ObjectCmd)

	// Exécuter la commande principale
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
