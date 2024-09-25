package commands

import (
	"fmt"
	"os"
	"plateforme-mycli/storage"

	"github.com/spf13/cobra"
)

// ObjectCmd est la commande racine pour les opérations sur les objets dans un bucket
var ObjectCmd = &cobra.Command{
	Use:   "object",
	Short: "Gérer les objets dans les buckets S3",
}

// Commande pour uploader un objet
var uploadObjectCmd = &cobra.Command{
	Use:   "upload [bucket-name] [object-name] [file-path]",
	Short: "Uploader un objet dans un bucket",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		objectName := args[1]
		filePath := args[2]
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la lecture du fichier : %v\n", err)
			os.Exit(1)
		}
		err = storage.UploadObject(bucketName, objectName, fileContent, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de l'upload de l'objet : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Objet '%s' uploadé avec succès dans le bucket '%s'\n", objectName, bucketName)
	},
}

// Commande pour télécharger un objet
var downloadObjectCmd = &cobra.Command{
	Use:   "download [bucket-name] [object-name] [file-path]",
	Short: "Télécharger un objet depuis un bucket",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		objectName := args[1]
		filePath := args[2]
		objectContent, err := storage.DownloadObject(bucketName, objectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors du téléchargement de l'objet : %v\n", err)
			os.Exit(1)
		}
		err = os.WriteFile(filePath, objectContent, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de l'écriture du fichier : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Objet '%s' téléchargé avec succès depuis le bucket '%s' dans '%s'\n", objectName, bucketName, filePath)
	},
}

// Commande pour supprimer un objet
var deleteObjectCmd = &cobra.Command{
	Use:   "delete [bucket-name] [object-name]",
	Short: "Supprimer un objet d'un bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		objectName := args[1]
		err := storage.DeleteObject(bucketName, objectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la suppression de l'objet : %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Objet '%s' supprimé avec succès du bucket '%s'\n", objectName, bucketName)
	},
}

// Commande pour lister les objets dans un bucket
var listObjectsCmd = &cobra.Command{
	Use:   "list [bucket-name]",
	Short: "Lister tous les objets dans un bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		err := storage.ListObjects(bucketName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la liste des objets : %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Ajouter les sous-commandes à la commande object
	ObjectCmd.AddCommand(uploadObjectCmd)
	ObjectCmd.AddCommand(downloadObjectCmd)
	ObjectCmd.AddCommand(deleteObjectCmd)
	ObjectCmd.AddCommand(listObjectsCmd)
}
