package s3

import (
	"fmt"
	"net/http"
)

func ListBuckets() {
	resp, err := http.Get("http://your-s3-url/buckets")
	if err != nil {
		fmt.Println("Erreur lors de la récupération des buckets:", err)
		return
	}
	defer resp.Body.Close()

	// Traitez la réponse
}
