package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ridha-boughediri/plateforme-mycli/libs"
)

func AddObject(url, localPath string) error {

	parts := strings.SplitN(url, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("alias and bucket should be in the correct format alias/bucket")
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

func DownloadObject(alias, url, localPath string) error {

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("error with provided alias: %v", err)
	}

	fullURL := remote + url

	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download object, server responded with status: %d", resp.StatusCode)
	}

	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("error creating local file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing file content: %v", err)
	}

	return nil
}

func DeleteObject(alias, url string) error {

	remote, err := libs.FindAlias(alias)
	if err != nil {
		return fmt.Errorf("error with provided alias: %v", err)
	}

	fullURL := remote + url

	req, err := http.NewRequest("DELETE", fullURL, nil)
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
		return fmt.Errorf("failed to delete object, server responded with status: %d", resp.StatusCode)
	}

	return nil
}
