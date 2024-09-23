package handlers

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ridha-boughediri/plateforme-mycli/dtos"
	"github.com/ridha-boughediri/plateforme-mycli/libs"
)

func printBuckets(body []byte) error {
	bucketResp := dtos.ListAllMyBucketsResult{}
	err := xml.Unmarshal(body, &bucketResp)
	if err != nil {
		return fmt.Errorf("error decoding XML: %v", err)
	}

	if len(bucketResp.Buckets) == 0 {
		fmt.Println("No buckets found.")
	} else {
		for _, bucket := range bucketResp.Buckets {
			creationDate := bucket.CreationDate.Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] - %s\n", creationDate, bucket.Name)
		}
	}
	return nil
}

func printObjects(body []byte) error {
	objectResp := dtos.ListBucketResult{}
	err := xml.Unmarshal(body, &objectResp)
	if err != nil {
		return fmt.Errorf("error decoding XML: %v", err)
	}

	if len(objectResp.Contents) == 0 {
		fmt.Println("No objects found.")
	} else {
		for _, object := range objectResp.Contents {
			lastModified := object.LastModified.Format("2006-01-02 15:04:05")
			fmt.Printf("[%s] %dB - %s\n", lastModified, object.Size, object.Key)
		}
	}
	return nil
}

func AddBucket(url string) error {

	parts, err := libs.UrlParts(url)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	alias := parts[0]
	bucket := parts[1]

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fullURL := fmt.Sprintf("%s/%s/", remote, bucket)

	req, err := http.NewRequest("PUT", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating bucket: %s", resp.Status)
	}

	return nil
}

func ListBuckets(url string) error {

	alias := ""
	bucketName := ""

	if strings.Contains(url, "/") {
		parts, err := libs.UrlParts(url)
		if err != nil {
			return fmt.Errorf("error parsing URL: %v", err)
		}

		alias = parts[0]
		bucketName = parts[1]
	} else {
		alias = url
	}

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fullURL := ""

	if bucketName == "" {
		fullURL = fmt.Sprintf("%s/", remote)
	} else {
		fullURL = fmt.Sprintf("%s/%s/", remote, bucketName)
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if bucketName == "" {
		err = printBuckets(body)
		if err != nil {
			return fmt.Errorf("error listing buckets: %v", err)
		}
	} else {
		err = printObjects(body)
		if err != nil {
			return fmt.Errorf("error listing objects: %v", err)
		}
	}
	return err
}

func DeleteBucket(url string) error {

	parts, err := libs.UrlParts(url)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	alias := parts[0]
	bucketName := parts[1]

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fullURL := fmt.Sprintf("%s/%s/", remote, bucketName)

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting bucket: %s", resp.Status)
	}

	return nil
}
