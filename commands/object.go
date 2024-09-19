package commands

import (
	"fmt"
	"io/ioutil"
	"plateforme-mycli/storage"

	"github.com/spf13/cobra"
)

// ObjectCmd est la commande principale pour les opérations sur les objets
var ObjectCmd = &cobra.Command{
	Use:   "object",
	Short: "Gérer les objets S3",
	Long:  "Commande pour uploader, télécharger et supprimer des objets dans un bucket S3",
}

// UploadObjectCmd pour uploader un objet
var UploadObjectCmd = &cobra.Command{
	Use:   "upload [bucket-name] [file-path]",
	Short: "Uploader un objet dans un bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		filePath := args[1]

		// Lire le fichier
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier :", err)
			return
		}

		// Uploader l'objet
		err = storage.UploadObject(bucketName, filePath, content)
		if err != nil {
			fmt.Println("Erreur lors de l'upload de l'objet :", err)
			return
		}
	},
}

// DownloadObjectCmd pour télécharger un objet
var DownloadObjectCmd = &cobra.Command{
	Use:   "download [bucket-name] [object-name]",
	Short: "Télécharger un objet d'un bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		objectName := args[1]

		// Télécharger l'objet
		content, err := storage.DownloadObject(bucketName, objectName)
		if err != nil {
			fmt.Println("Erreur lors du téléchargement de l'objet :", err)
			return
		}

		// Afficher le contenu téléchargé
		fmt.Println("Contenu de l'objet :", string(content))
	},
}

// DeleteObjectCmd pour supprimer un objet
var DeleteObjectCmd = &cobra.Command{
	Use:   "delete [bucket-name] [object-name]",
	Short: "Supprimer un objet d'un bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		objectName := args[1]

		// Supprimer l'objet
		err := storage.DeleteObject(bucketName, objectName)
		if err != nil {
			fmt.Println("Erreur lors de la suppression de l'objet :", err)
			return
		}
	},
}

func init() {
	// Ajouter les sous-commandes au commandement principal "object"
	ObjectCmd.AddCommand(UploadObjectCmd)
	ObjectCmd.AddCommand(DownloadObjectCmd)
	ObjectCmd.AddCommand(DeleteObjectCmd)
}
