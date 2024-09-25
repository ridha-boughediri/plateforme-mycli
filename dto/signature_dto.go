package dto

// SignatureV4Request contient les informations nécessaires à la création d'une signature AWS V4
type SignatureV4Request struct {
	Method     string
	BucketName string
	ObjectName string
	Region     string
	AccessKey  string
	SecretKey  string
	Timestamp  string
}

// SignatureV4Response représente la réponse contenant l'en-tête signé
type SignatureV4Response struct {
	Headers map[string]string
}
