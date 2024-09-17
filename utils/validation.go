package utils

import "regexp"

// Fonction pour valider le nom d'un bucket
func IsValidBucketName(bucketName string) bool {
	// Les noms de bucket doivent respecter les conventions S3 (lettres minuscules, chiffres, points et tirets)
	var bucketNameRegex = regexp.MustCompile(`^[a-z0-9.-]{3,63}$`)
	return bucketNameRegex.MatchString(bucketName)
}
