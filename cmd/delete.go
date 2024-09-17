package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"plateforme-mycli/utils"

	"github.com/spf13/cobra"
)

// Fonction pour supprimer un bucket
func deleteBucket(bucketName string) {
	// Valider le nom du bucket avant de continuer
	if !utils.IsValidBucketName(bucketName) {
		log.Fatalf("Nom de bucket invalide : %s.", bucketName)
	}

	url := fmt.Sprintf("http://localhost:8080/bucket/%s", bucketName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la requête de suppression : %v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de la requête de suppression : %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Bucket supprimé avec succès :", bucketName)
	case http.StatusNotFound:
		log.Printf("Le bucket %s n'existe pas.", bucketName)
	default:
		log.Printf("Erreur lors de la suppression du bucket, statut : %s, réponse : %s", resp.Status, string(body))
	}
}

// Déclaration de la commande "delete"
var DeleteCmd = &cobra.Command{
	Use:   "delete [bucketName]",
	Short: "Delete a bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteBucket(args[0])
	},
}
