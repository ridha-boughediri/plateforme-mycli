package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

var s3Client *http.Client

// InitS3Client initialise le client HTTP pour interagir avec le serveur S3
func InitS3Client() {
	s3Client = &http.Client{}
}

// --- Utilitaires de signature et hashage ---

// Génère une signature SHA-256 pour AWS
func sha256Hex(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

// Génère un HMAC-SHA256 pour AWS
func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// Génère la clé de signature AWS V4
func getSignatureKey(key, dateStamp, regionName, serviceName string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+key), dateStamp)
	kRegion := hmacSHA256(kDate, regionName)
	kService := hmacSHA256(kRegion, serviceName)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}

// Crée une signature AWS V4 pour les requêtes HTTP
func createSignatureV4(method, bucketName, objectName, region, accessKey, secretKey string, t time.Time, payloadHash, contentType string) (http.Header, error) {
	service := "s3"
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Extraire le host de l'endpoint
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid endpoint URL: %v", err)
	}
	host := endpointURL.Host

	amzDate := t.UTC().Format("20060102T150405Z")
	dateStamp := t.UTC().Format("20060102")
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, region, service)

	// Construire l'URI canonique
	canonicalURI := "/"
	if bucketName != "" {
		canonicalURI = fmt.Sprintf("/%s", bucketName)
		if objectName != "" {
			canonicalURI = fmt.Sprintf("/%s/%s", bucketName, objectName)
		}
	}

	// Construire les en-têtes canoniques
	headersToSign := map[string]string{
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}
	if contentType != "" {
		headersToSign["content-type"] = contentType
	}

	// Trier les clés des en-têtes
	headerKeys := make([]string, 0, len(headersToSign))
	for k := range headersToSign {
		headerKeys = append(headerKeys, strings.ToLower(k))
	}
	sort.Strings(headerKeys)

	// Construire les en-têtes canoniques et les en-têtes signés
	var canonicalHeaders strings.Builder
	var signedHeaders strings.Builder
	for i, k := range headerKeys {
		v := headersToSign[k]
		canonicalHeaders.WriteString(fmt.Sprintf("%s:%s\n", strings.ToLower(k), strings.TrimSpace(v)))
		if i > 0 {
			signedHeaders.WriteString(";")
		}
		signedHeaders.WriteString(strings.ToLower(k))
	}

	// Construire la requête canonique
	canonicalRequest := fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s", method, canonicalURI, canonicalHeaders.String(), signedHeaders.String(), payloadHash)

	// Construire la chaîne à signer
	algorithm := "AWS4-HMAC-SHA256"
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", algorithm, amzDate, credentialScope, sha256Hex(canonicalRequest))

	// Calculer la signature
	signingKey := getSignatureKey(secretKey, dateStamp, region, service)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	// Construire l'en-tête Authorization
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, accessKey, credentialScope, signedHeaders.String(), signature)

	// Préparer les en-têtes HTTP
	headers := http.Header{}
	headers.Set("Authorization", authorizationHeader)
	headers.Set("x-amz-date", amzDate)
	headers.Set("x-amz-content-sha256", payloadHash)
	headers.Set("Host", host)
	if contentType != "" {
		headers.Set("Content-Type", contentType)
	}

	return headers, nil
}

// --- Gestion des buckets ---

// Structure pour parser la réponse XML des buckets
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets Buckets  `xml:"Buckets"`
}

type Buckets struct {
	Bucket []Bucket `xml:"Bucket"`
}

type Bucket struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

// Fonction pour lister les buckets
func ListBuckets() error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Payload vide pour une requête GET sans corps
	payloadHash := sha256Hex("")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("GET", "", "", region, accessKeyID, secretAccessKey, t, payloadHash, "")
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL pour lister les buckets
	url := fmt.Sprintf("%s/", endpoint)

	// Créer la requête HTTP GET
	req, err := http.NewRequest("GET", url, nil)
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

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec lors de la récupération des buckets, statut : %s, réponse : %s", resp.Status, string(body))
	}

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du corps de la réponse : %v", err)
	}

	// Parser le XML
	var result ListAllMyBucketsResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("erreur lors du parsing du XML : %v", err)
	}

	// Afficher les buckets de manière lisible
	fmt.Println("Liste des buckets récupérée avec succès :")
	for _, bucket := range result.Buckets.Bucket {
		fmt.Printf("- %s (créé le %s)\n", bucket.Name, bucket.CreationDate)
	}

	return nil
}

// Fonction pour créer un bucket
func CreateBucket(bucketName string) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Préparer le corps de la requête si nécessaire
	var requestBody string
	var contentType string
	if region != "" && region != "us-east-1" {
		requestBody = fmt.Sprintf(
			`<CreateBucketConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><LocationConstraint>%s</LocationConstraint></CreateBucketConfiguration>`,
			region)
		contentType = "application/xml"
	}

	// Calculer le hash du payload
	payloadHash := sha256Hex(requestBody)

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("PUT", bucketName, "", region, accessKeyID, secretAccessKey, t, payloadHash, contentType)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL en utilisant le style d'adressage path-style
	canonicalURI := "/"
	if bucketName != "" {
		canonicalURI = fmt.Sprintf("/%s", bucketName)
	}
	url := fmt.Sprintf("%s%s", endpoint, canonicalURI)

	// Créer la requête HTTP
	req, err := http.NewRequest("PUT", url, strings.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header = headers
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(requestBody)))
	}

	// Envoyer la requête
	resp, err := s3Client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec de la création du bucket, statut : %s, réponse : %s", resp.Status, string(body))
	}

	fmt.Println("Bucket créé avec succès :", bucketName)
	return nil
}

// Fonction pour supprimer un bucket
func DeleteBucket(bucketName string) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Le payload est vide pour DeleteBucket
	payloadHash := sha256Hex("")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("DELETE", bucketName, "", region, accessKeyID, secretAccessKey, t, payloadHash, "")
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL
	url := fmt.Sprintf("%s/%s", endpoint, bucketName)

	// Créer la requête HTTP
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

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec de la suppression du bucket, statut : %s, réponse : %s", resp.Status, string(body))
	}

	fmt.Println("Bucket supprimé avec succès :", bucketName)
	return nil
}

// --- Gestion des objets ---

// Structure pour parser la réponse XML de la liste des objets
type ListBucketResult struct {
	XMLName xml.Name   `xml:"ListBucketResult"`
	Objects []S3Object `xml:"Contents"`
}

type S3Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int64  `xml:"Size"`
}

// Fonction pour lister les objets dans un bucket
func ListObjects(bucketName string) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Payload vide pour une requête GET sans corps
	payloadHash := sha256Hex("")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("GET", bucketName, "", region, accessKeyID, secretAccessKey, t, payloadHash, "")
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL pour lister les objets dans un bucket
	url := fmt.Sprintf("%s/%s", endpoint, bucketName)

	// Créer la requête HTTP GET
	req, err := http.NewRequest("GET", url, nil)
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

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec lors de la récupération des objets, statut : %s, réponse : %s", resp.Status, string(body))
	}

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du corps de la réponse : %v", err)
	}

	// Parser le XML
	var result ListBucketResult
	err = xml.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("erreur lors du parsing du XML : %v", err)
	}

	// Afficher les objets
	fmt.Println("Liste des objets récupérée avec succès :")
	for _, object := range result.Objects {
		fmt.Printf("- %s (taille : %d octets, modifié le %s)\n", object.Key, object.Size, object.LastModified)
	}

	return nil
}

// Fonction pour uploader un objet
func UploadObject(bucketName, objectName string, fileContent []byte, contentType string) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Calculer le hash du payload
	payloadHash := sha256Hex(string(fileContent))

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("PUT", bucketName, objectName, region, accessKeyID, secretAccessKey, t, payloadHash, contentType)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName)

	// Créer la requête HTTP
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header = headers
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Envoyer la requête
	resp, err := s3Client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec de l'upload de l'objet, statut : %s, réponse : %s", resp.Status, string(body))
	}

	fmt.Println("Objet uploadé avec succès :", objectName)
	return nil
}

// Fonction pour télécharger un objet
func DownloadObject(bucketName, objectName string) ([]byte, error) {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Le payload est vide pour les requêtes GET
	payloadHash := sha256Hex("")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("GET", bucketName, objectName, region, accessKeyID, secretAccessKey, t, payloadHash, "")
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName)

	// Créer la requête HTTP
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
		return nil, fmt.Errorf("échec du téléchargement de l'objet, statut : %s, réponse : %s", resp.Status, string(body))
	}

	// Lire le contenu de l'objet
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du corps de la réponse : %v", err)
	}

	fmt.Println("Objet téléchargé avec succès :", objectName)
	return body, nil
}

// Fonction pour supprimer un objet
func DeleteObject(bucketName, objectName string) error {
	region := os.Getenv("Region")
	accessKeyID := os.Getenv("Access_Key")
	secretAccessKey := os.Getenv("Secret_Key")
	endpoint := os.Getenv("Endpoint")
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	// Le payload est vide pour les requêtes DELETE
	payloadHash := sha256Hex("")

	// Créer la signature AWS V4
	t := time.Now()
	headers, err := createSignatureV4("DELETE", bucketName, objectName, region, accessKeyID, secretAccessKey, t, payloadHash, "")
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la signature : %v", err)
	}

	// Construire l'URL
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName)

	// Créer la requête HTTP
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

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("échec de la suppression de l'objet, statut : %s, réponse : %s", resp.Status, string(body))
	}

	fmt.Println("Objet supprimé avec succès :", objectName)
	return nil
}
