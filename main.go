package main

import (
	"log"

	"plateforme-mycli/commands"
	"plateforme-mycli/config"
	"plateforme-mycli/storage"
)

func main() {
	// Charger la configuration
	config.LoadEnv()

	// Initialiser le client MinIO
	err := storage.InitMinioClient()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation du client MinIO : %v", err)
	}

	// Exécuter les commandes CLI
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}
