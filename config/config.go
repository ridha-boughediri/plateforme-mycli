package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv charge les variables d'environnement depuis le fichier .env
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur de chargement du fichier .env : %v", err)
	}
}
