package dto

// UploadObjectRequest représente les données pour uploader un objet
type UploadObjectRequest struct {
	BucketName  string `json:"bucket_name"`
	ObjectName  string `json:"object_name"`
	FileContent []byte `json:"file_content"`
}

// DownloadObjectResponse représente la réponse après téléchargement d'un objet
type DownloadObjectResponse struct {
	ObjectContent []byte `json:"object_content"`
}
