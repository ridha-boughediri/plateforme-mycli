package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"plateforme-mycli/utils"

	"github.com/spf13/cobra"
)

// Fonction pour créer un bucket
func createBucket(bucketName string) {
	// Valider le nom du bucket avant de continuer
	if !utils.IsValidBucketName(bucketName) {
		log.Fatalf("Nom de bucket invalide : %s.", bucketName)
	}

	// Construction de l'URL pour la requête de création du bucket
	url := fmt.Sprintf("http://localhost:8080/bucket/%s", bucketName)

	// Création de la requête HTTP PUT pour créer le bucket
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la requête : %v", err)
	}

	// Client HTTP pour exécuter la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de la requête : %v", err)
	}
	defer resp.Body.Close()

	// Lecture du corps de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
	}

	// Gestion des statuts HTTP
	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Bucket créé avec succès :", bucketName)
	case http.StatusConflict:
		log.Printf("Le bucket %s existe déjà.", bucketName)
	default:
		log.Printf("Erreur lors de la création du bucket, statut : %s, réponse : %s", resp.Status, string(body))
	}
}

// Déclaration de la commande "create"
var CreateCmd = &cobra.Command{
	Use:   "create [bucketName]",
	Short: "Create a new bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createBucket(args[0])
	},
}
