package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Attendu: list-buckets, create-bucket, upload-file, etc.")
		return
	}

	switch os.Args[1] {
	case "list-buckets":
		listBuckets()
	case "create-bucket":
		createBucket()
	// Ajoutez ici les autres commandes
	default:
		fmt.Println("Commande inconnue:", os.Args[1])
	}
}

func listBuckets() {
	// Appel à l'API S3 pour lister les buckets
	fmt.Println("Liste des buckets...")
}

func createBucket() {
	// Appel à l'API S3 pour créer un bucket
	fmt.Println("Création d'un bucket...")
}
