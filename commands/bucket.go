package commands

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"plateforme-mycli/storage" // Assurez-vous que ce chemin correspond bien à votre projet

	"github.com/spf13/cobra"
)

var s3Client *http.Client

// HMAC-SHA256 pour la signature (réutilisé)
func signHMACSHA256(key, data string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Créer la signature V4 (réutilisé)
func createSignatureV4(method, bucketName, objectName, region, accessKey, secretKey string, t time.Time) (http.Header, error) {
	service := "s3"
	endpoint := os.Getenv("Endpoint")
	host := fmt.Sprintf("%s.%s", bucketName, endpoint)
	amzDate := t.UTC().Format("20060102T150405Z")
	dateStamp := t.UTC().Format("20060102")
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, service)

	// Créer une chaîne de signature (Canonical Request)
	canonicalURI := fmt.Sprintf("/%s", objectName)
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-date:%s\n", host, amzDate)
	signedHeaders := "host;x-amz-date"
	payloadHash := sha256Hex("")
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", method, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)

	// String to sign
	algorithm := "AWS4-HMAC-SHA256"
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", algorithm, amzDate, credentialScope, sha256Hex(canonicalRequest))

	// Clé de signature
	kSecret := "AWS4" + secretKey
	kDate := signHMACSHA256(kSecret, dateStamp)
	kRegion := signHMACSHA256(string(kDate), region)
	kService := signHMACSHA256(string(kRegion), service)
	kSigning := signHMACSHA256(string(kService), "aws4_request")
	signature := hex.EncodeToString(signHMACSHA256(string(kSigning), stringToSign))

	// Créer l'en-tête Authorization
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", algorithm, accessKey, credentialScope, signedHeaders, signature)

	headers := http.Header{}
	headers.Set("Authorization", authorizationHeader)
	headers.Set("x-amz-date", amzDate)
	return headers, nil
}

// SHA256 pour une chaîne vide ou des données (réutilisé)
func sha256Hex(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// UploadObject pour uploader un objet dans un bucket
func UploadObject(bucketName, objectName string, fileContent []byte) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("PUT", bucketName, objectName, region, accessKeyID, secretAccessKey, t)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL et la requête HTTP
	url := fmt.Sprintf("https://%s.%s/%s", bucketName, os.Getenv("Endpoint"), objectName)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header = headers

	// Envoyer la requête
	resp, err := s3Client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec de l'upload de l'objet : %s, réponse : %s", resp.Status, string(body))
	}

	fmt.Println("Objet uploadé avec succès :", objectName)
	return nil
}

// DownloadObject pour télécharger un objet d'un bucket
func DownloadObject(bucketName, objectName string) ([]byte, error) {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("GET", bucketName, objectName, region, accessKeyID, secretAccessKey, t)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL et la requête HTTP
	url := fmt.Sprintf("https://%s.%s/%s", bucketName, os.Getenv("Endpoint"), objectName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header = headers

	// Envoyer la requête
	resp, err := s3Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'envoi de la requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("échec du téléchargement de l'objet : %s, réponse : %s", resp.Status, string(body))
	}

	// Lire et retourner le contenu de l'objet
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse : %v", err)
	}

	fmt.Println("Objet téléchargé avec succès :", objectName)
	return body, nil
}

// DeleteObject pour supprimer un objet d'un bucket
func DeleteObject(bucketName, objectName string) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("DELETE", bucketName, objectName, region, accessKeyID, secretAccessKey, t)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL et la requête HTTP
	url := fmt.Sprintf("https://%s.%s/%s", bucketName, os.Getenv("Endpoint"), objectName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header = headers

	// Envoyer la requête
	resp, err := s3Client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec de la suppression de l'objet : %s, réponse : %s", resp.Status, string(body))
	}

	fmt.Println("Objet supprimé avec succès :", objectName)
	return nil
}

// BucketCmd est la commande principale pour les opérations sur les buckets
var BucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Gérer les buckets S3",
	Long:  "Commande pour créer, lister et supprimer des buckets sur le serveur S3",
}

// init ajoute les sous-commandes (create, list, delete) à BucketCmd
func init() {
	BucketCmd.AddCommand(&cobra.Command{
		Use:   "create [bucket-name]",
		Short: "Créer un bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bucketName := args[0]
			err := storage.CreateBucket(bucketName)
			if err != nil {
				fmt.Println("Erreur lors de la création du bucket :", err)
				return
			}
			fmt.Println("Bucket créé avec succès :", bucketName)
		},
	})

	BucketCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Lister les buckets",
		Run: func(cmd *cobra.Command, args []string) {
			err := storage.ListBuckets()
			if err != nil {
				fmt.Println("Erreur lors de la récupération des buckets :", err)
				return
			}
		},
	})

	BucketCmd.AddCommand(&cobra.Command{
		Use:   "delete [bucket-name]",
		Short: "Supprimer un bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bucketName := args[0]
			err := storage.DeleteBucket(bucketName)
			if err != nil {
				fmt.Println("Erreur lors de la suppression du bucket :", err)
				return
			}
			fmt.Println("Bucket supprimé avec succès :", bucketName)
		},
	})
}
