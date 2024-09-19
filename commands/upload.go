package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func uploadObject(bucketName string, objectName string) {
	// Construire l'URL de requête pour uploader l'objet
	url := fmt.Sprintf("http://localhost:8080/%s/%s", bucketName, objectName)

	// Ouvrir le fichier à uploader
	file, err := os.Open(objectName)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	// Créer une requête PUT
	req, err := http.NewRequest("PUT", url, file)
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

	// Vérification du statut de la réponse
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Objet uploadé avec succès :", objectName)
	} else {
		log.Printf("Erreur lors de l'upload de l'objet, statut : %d", resp.StatusCode)
	}
}

// Déclaration de la commande "upload"
var UploadCmd = &cobra.Command{
	Use:   "upload [bucketName] [objectName]",
	Short: "Upload an object to a bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		uploadObject(args[0], args[1])
	},
}
