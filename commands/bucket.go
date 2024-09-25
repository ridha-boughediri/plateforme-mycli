package commands

import (
	"fmt"
	"os"
	"plateforme-mycli/storage"

	"github.com/spf13/cobra"
)

// BucketCmd est la commande racine pour les opérations sur les buckets
var BucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Gérer les buckets S3",
}

// Commande pour créer un bucket
var createBucketCmd = &cobra.Command{
	Use:   "create [bucket-name]",
	Short: "Créer un nouveau bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		err := storage.CreateBucket(bucketName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la création du bucket : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bucket '%s' créé avec succès\n", bucketName)
	},
}

// Commande pour supprimer un bucket
var deleteBucketCmd = &cobra.Command{
	Use:   "delete [bucket-name]",
	Short: "Supprimer un bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		err := storage.DeleteBucket(bucketName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la suppression du bucket : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Bucket '%s' supprimé avec succès\n", bucketName)
	},
}

// Commande pour lister les buckets
var listBucketsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lister tous les buckets",
	Run: func(cmd *cobra.Command, args []string) {
		err := storage.ListBuckets()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la liste des buckets : %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Ajouter les sous-commandes à la commande bucket
	BucketCmd.AddCommand(createBucketCmd)
	BucketCmd.AddCommand(deleteBucketCmd)
	BucketCmd.AddCommand(listBucketsCmd)
}
