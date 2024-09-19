package commands

import (
	"context"
	"fmt"
	"log"

	"plateforme-mycli/storage"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(bucketCmd)
}

var bucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Gérer les buckets",
}

var createBucketCmd = &cobra.Command{
	Use:   "create [bucket-name]",
	Short: "Créer un bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		ctx := context.Background()
		err := storage.MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Erreur lors de la création du bucket : %v", err)
		}
		fmt.Printf("Bucket '%s' créé avec succès.\n", bucketName)
	},
}

var listBucketsCmd = &cobra.Command{
	Use:   "list",
	Short: "Lister les buckets",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		buckets, err := storage.MinioClient.ListBuckets(ctx)
		if err != nil {
			log.Fatalf("Erreur lors de la récupération des buckets : %v", err)
		}
		fmt.Println("Liste des buckets :")
		for _, bucket := range buckets {
			fmt.Println(bucket.Name)
		}
	},
}

var deleteBucketCmd = &cobra.Command{
	Use:   "delete [bucket-name]",
	Short: "Supprimer un bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		ctx := context.Background()
		err := storage.MinioClient.RemoveBucket(ctx, bucketName)
		if err != nil {
			log.Fatalf("Erreur lors de la suppression du bucket : %v", err)
		}
		fmt.Printf("Bucket '%s' supprimé avec succès.\n", bucketName)
	},
}

func init() {
	bucketCmd.AddCommand(createBucketCmd)
	bucketCmd.AddCommand(listBucketsCmd)
	bucketCmd.AddCommand(deleteBucketCmd)
}
