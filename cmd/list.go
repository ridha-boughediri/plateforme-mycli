package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// Fonction pour lister les buckets
func listBuckets() {
	url := "http://localhost:8080/"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des buckets : %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("List Buckets Response:", string(body))
	} else {
		log.Printf("Erreur lors de la récupération des buckets, statut : %s, réponse : %s", resp.Status, string(body))
	}
}

// Déclaration de la commande "list"
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all buckets",
	Run: func(cmd *cobra.Command, args []string) {
		listBuckets()
	},
}
