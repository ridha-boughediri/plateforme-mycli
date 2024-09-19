package storage

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
)

var s3Client *http.Client

func InitS3Client() {
	s3Client = &http.Client{}
}

// HMAC-SHA256 pour la signature
func signHMACSHA256(key, data string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Créer la signature V4 pour les requêtes S3
func createSignatureV4(method, bucketName, objectName, region, accessKey, secretKey string, t time.Time) (http.Header, error) {
	service := "s3"
	endpoint := os.Getenv("Endpoint")
	host := fmt.Sprintf("%s.%s", bucketName, endpoint)
	amzDate := t.UTC().Format("20060102T150405Z")
	dateStamp := t.UTC().Format("20060102")
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, service)

	// Chaîne canonique (Canonical Request)
	canonicalURI := fmt.Sprintf("/%s", objectName)
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-date:%s\n", host, amzDate)
	signedHeaders := "host;x-amz-date"
	payloadHash := sha256Hex("")
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", method, canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)

	// String to sign
	algorithm := "AWS4-HMAC-SHA256"
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", algorithm, amzDate, credentialScope, sha256Hex(canonicalRequest))

	// Générer la signature
	kSecret := "AWS4" + secretKey
	kDate := signHMACSHA256(kSecret, dateStamp)
	kRegion := signHMACSHA256(string(kDate), region)
	kService := signHMACSHA256(string(kRegion), service)
	kSigning := signHMACSHA256(string(kService), "aws4_request")
	signature := hex.EncodeToString(signHMACSHA256(string(kSigning), stringToSign))

	// En-tête Authorization
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s", algorithm, accessKey, credentialScope, signedHeaders, signature)

	headers := http.Header{}
	headers.Set("Authorization", authorizationHeader)
	headers.Set("x-amz-date", amzDate)
	return headers, nil
}

// SHA256 pour une chaîne vide ou des données
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

// InitS3Client initialise un client HTTP pour interagir avec le serveur S3
// var s3Client *http.Client

// func InitS3Client() {
// 	s3Client = &http.Client{}
// }

// CreateBucket pour créer un bucket en envoyant une requête PUT
func CreateBucket(bucketName string) error {
	endpoint := os.Getenv("Endpoint")
	url := fmt.Sprintf("%s/%s", endpoint, bucketName)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête: %v", err)
	}

	resp, err := s3Client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("échec de la création du bucket, statut: %s", resp.Status)
	}

	fmt.Println("Bucket créé avec succès :", bucketName)
	return nil
}

// ListBuckets pour lister les buckets en envoyant une requête GET
func ListBuckets() error {
	endpoint := os.Getenv("Endpoint")
	url := fmt.Sprintf("%s/", endpoint)

	resp, err := s3Client.Get(url)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de la liste des buckets: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("échec lors de la récupération des buckets, statut: %s", resp.Status)
	}

	fmt.Println("Liste des buckets récupérée avec succès.")
	return nil
}

// DeleteBucket pour supprimer un bucket en envoyant une requête DELETE
func DeleteBucket(bucketName string) error {
	endpoint := os.Getenv("Endpoint")
	url := fmt.Sprintf("%s/%s", endpoint, bucketName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête: %v", err)
	}

	resp, err := s3Client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("échec de la suppression du bucket, statut: %s", resp.Status)
	}

	fmt.Println("Bucket supprimé avec succès :", bucketName)
	return nil
}
