package commands

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// Définir les structures pour parser le XML
type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

// Fonction pour parser le XML
func parseXML(xmlData string) {
	// Création d'un décodeur XML personnalisé
	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	var result ListAllMyBucketsResult
	err := decoder.Decode(&result)
	if err != nil {
		log.Fatalf("Erreur lors du parsing XML : %v", err)
	}

	// Afficher les buckets et leurs dates de création
	fmt.Println("Liste des Buckets :")
	for _, bucket := range result.Buckets {
		fmt.Printf("Bucket: %s, Créé le: %s\n", bucket.Name, bucket.CreationDate)
	}
}

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
		// Appel du parser pour afficher les résultats
		parseXML(string(body))
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
