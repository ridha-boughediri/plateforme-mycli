package commands

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"plateforme-mycli/storage"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(objectCmd)
}

var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "Gérer les objets",
}

var uploadObjectCmd = &cobra.Command{
	Use:   "upload [bucket-name] [file-path]",
	Short: "Uploader un objet",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		filePath := args[1]
		fileName := filepath.Base(filePath)

		ctx := context.Background()
		_, err := storage.MinioClient.FPutObject(ctx, bucketName, fileName, filePath, minio.PutObjectOptions{})
		if err != nil {
			log.Fatalf("Erreur lors de l'upload de l'objet : %v", err)
		}
		fmt.Printf("Fichier '%s' uploadé dans le bucket '%s'.\n", fileName, bucketName)
	},
}

var deleteObjectCmd = &cobra.Command{
	Use:   "delete [bucket-name] [object-name]",
	Short: "Supprimer un objet",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		objectName := args[1]

		ctx := context.Background()
		err := storage.MinioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Fatalf("Erreur lors de la suppression de l'objet : %v", err)
		}
		fmt.Printf("Objet '%s' supprimé du bucket '%s'.\n", objectName, bucketName)
	},
}

// Nouvelle commande pour lister les objets d'un bucket
var listObjectsCmd = &cobra.Command{
	Use:   "list [bucket-name]",
	Short: "Lister les objets d'un bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		ctx := context.Background()

		// Vérifier si le bucket existe
		exists, err := storage.MinioClient.BucketExists(ctx, bucketName)
		if err != nil {
			log.Fatalf("Erreur lors de la vérification du bucket : %v", err)
		}
		if !exists {
			log.Fatalf("Le bucket '%s' n'existe pas.", bucketName)
		}

		// Créer un canal pour recevoir les objets
		objectCh := storage.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Recursive: true,
		})

		fmt.Printf("Liste des objets dans le bucket '%s':\n", bucketName)
		for object := range objectCh {
			if object.Err != nil {
				log.Fatalf("Erreur lors de la liste des objets : %v", object.Err)
			}
			fmt.Println(object.Key)
		}
	},
}

func init() {
	objectCmd.AddCommand(uploadObjectCmd)
	objectCmd.AddCommand(deleteObjectCmd)
	objectCmd.AddCommand(listObjectsCmd) // Ajouter la commande listObjectsCmd
}
