package handlers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ridha-boughediri/plateforme-mycli/dtos"
	"github.com/ridha-boughediri/plateforme-mycli/libs"
)

func AddObject(url, localPath string) error {

	parts, err := libs.UrlParts(url)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	alias := parts[0]
	bucket := parts[1]

	fileInfo, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("object does not exist at path %s", localPath)
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fullURL := fmt.Sprintf("%s/%s/%s", remote, bucket, fileInfo.Name())

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file content %v", err)
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("error creating HTTP request %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("File-Name", filepath.Base(fileInfo.Name()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload object server responded with %s \nStatus: %d", resp.Status, resp.StatusCode)
	}

	fmt.Println("Upload successful.")
	return nil
}

func DownloadObject(url, localPath string) error {

	parts, err := libs.UrlParts(url)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	alias := parts[0]
	bucket := parts[1]

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("error with provided alias: %v", err)
	}

	fullURL := fmt.Sprintf("%s/%s", remote, bucket)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download object server responded with %s \nStatus: %d", resp.Status, resp.StatusCode)
	}

	fileInfo, err := os.Stat(localPath)
	if err == nil && fileInfo.IsDir() {
		fileName := filepath.Base(bucket)
		localPath = filepath.Join(localPath, fileName)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("error creating file at path %s: %v", localPath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func DeleteObject(url string) error {

	parts, err := libs.UrlParts(url)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	alias := parts[0]
	bucket := parts[1]

	parts2, err := libs.UrlParts(bucket)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}

	bucketName := parts2[0]
	objectName := parts2[1]

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("error with provided alias: %v", err)
	}

	fullURL := fmt.Sprintf("%s/%s/", remote, bucketName)

	println(fullURL)

	deleteReq := dtos.Delete{
		Quiet: false,
		Object: struct {
			Key string `xml:"Key"`
		}{
			Key: objectName,
		},
	}

	xmlPayload, err := xml.Marshal(deleteReq)
	if err != nil {
		return fmt.Errorf("error marshalling XML: %v", err)
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(xmlPayload))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete object server responded with %s \nStatus: %d", resp.Status, resp.StatusCode)
	}

	return nil
}
