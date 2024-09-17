// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"regexp"

// 	"github.com/spf13/cobra"
// )

// // Fonction pour valider le nom d'un bucket
// func isValidBucketName(bucketName string) bool {
// 	// Les noms de bucket doivent respecter les conventions S3 (lettres minuscules, chiffres, points et tirets)
// 	var bucketNameRegex = regexp.MustCompile(`^[a-z0-9.-]{3,63}$`)
// 	return bucketNameRegex.MatchString(bucketName)
// }

// // Fonction pour créer un bucket
// func createBucket(bucketName string) {
// 	// Valider le nom du bucket avant de continuer
// 	if !isValidBucketName(bucketName) {
// 		log.Fatalf("Nom de bucket invalide : %s. Les noms doivent contenir uniquement des lettres minuscules, chiffres, points, et tirets (3-63 caractères).", bucketName)
// 	}

// 	// Construction de l'URL pour la requête de création du bucket
// 	url := fmt.Sprintf("http://localhost:8080/bucket/%s", bucketName)

// 	// Création de la requête HTTP PUT pour créer le bucket
// 	req, err := http.NewRequest("PUT", url, nil)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de la création de la requête : %v", err)
// 	}

// 	// Client HTTP pour exécuter la requête
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de l'envoi de la requête : %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Lecture du corps de la réponse
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
// 	}

// 	// Gestion des statuts HTTP pour des retours explicites
// 	switch resp.StatusCode {
// 	case http.StatusOK:
// 		fmt.Println("Bucket créé avec succès :", bucketName)
// 	case http.StatusConflict:
// 		log.Printf("Le bucket %s existe déjà.", bucketName)
// 	case http.StatusBadRequest:
// 		log.Printf("Requête incorrecte pour la création du bucket : %s", bucketName)
// 	default:
// 		log.Printf("Erreur lors de la création du bucket, statut : %s, réponse : %s", resp.Status, string(body))
// 	}
// }

// // Fonction pour lister les buckets
// func listBuckets() {
// 	url := "http://localhost:8080/"
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de la récupération des buckets : %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
// 	}

// 	// Vérifier si la requête a réussi
// 	if resp.StatusCode == http.StatusOK {
// 		fmt.Println("List Buckets Response:", string(body))
// 	} else {
// 		log.Printf("Erreur lors de la récupération des buckets, statut : %s, réponse : %s", resp.Status, string(body))
// 	}
// }

// // Fonction pour supprimer un bucket
// func deleteBucket(bucketName string) {
// 	// Valider le nom du bucket avant de continuer
// 	if !isValidBucketName(bucketName) {
// 		log.Fatalf("Nom de bucket invalide : %s.", bucketName)
// 	}

// 	url := fmt.Sprintf("http://localhost:8080/bucket/%s", bucketName)
// 	req, err := http.NewRequest("DELETE", url, nil)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de la création de la requête de suppression : %v", err)
// 	}
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de l'envoi de la requête de suppression : %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalf("Erreur lors de la lecture de la réponse : %v", err)
// 	}

// 	// Gestion des statuts HTTP
// 	switch resp.StatusCode {
// 	case http.StatusOK:
// 		fmt.Println("Bucket supprimé avec succès :", bucketName)
// 	case http.StatusNotFound:
// 		log.Printf("Le bucket %s n'existe pas.", bucketName)
// 	default:
// 		log.Printf("Erreur lors de la suppression du bucket, statut : %s, réponse : %s", resp.Status, string(body))
// 	}
// }

// // Fonction principale
// func main() {
// 	// Déclaration de la commande racine pour la CLI
// 	var rootCmd = &cobra.Command{Use: "myS3"}

// 	// Commande pour créer un bucket
// 	var cmdCreate = &cobra.Command{
// 		Use:   "create [bucketName]",
// 		Short: "Create a new bucket",
// 		Args:  cobra.ExactArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			createBucket(args[0])
// 		},
// 	}

// 	// Commande pour lister tous les buckets
// 	var cmdList = &cobra.Command{
// 		Use:   "list",
// 		Short: "List all buckets",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			listBuckets()
// 		},
// 	}

// 	// Commande pour supprimer un bucket
// 	var cmdDelete = &cobra.Command{
// 		Use:   "delete [bucketName]",
// 		Short: "Delete a bucket",
// 		Args:  cobra.ExactArgs(1),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			deleteBucket(args[0])
// 		},
// 	}

// 	// Ajout des commandes à la commande racine
// 	rootCmd.AddCommand(cmdCreate, cmdList, cmdDelete)

// 	// Exécution de la commande racine
// 	if err := rootCmd.Execute(); err != nil {
// 		log.Fatalf("Erreur lors de l'exécution de la commande : %v", err)
// 		os.Exit(1)
// 	}
// }

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
