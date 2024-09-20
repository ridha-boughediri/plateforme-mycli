package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/ridha-boughediri/plateforme-mycli/libs"
)

func runCommand(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestCreateAliasCommand(t *testing.T) {
	output, err := runCommand("sync", "testalias", "http://localhost:8080", "testuser", "testpassword")
	if err != nil {
		t.Fatalf("Failed to create alias: %v", err)
	}

	if !strings.Contains(output, "Alias saved successfully") {
		t.Errorf("Expected alias creation message, but got: %s", output)
	}
}

func TestCheckAliasFilePath(t *testing.T) {
	aliasFilePath, err := libs.GetAliasFilePath()
	if err != nil {
		t.Fatalf("Failed to get alias file path: %v", err)
	}
	fmt.Println("Alias file path:", aliasFilePath)
}

func TestCreateBucketCommand(t *testing.T) {
	output, err := runCommand("ba", "testalias/testbucket")
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}

	if !strings.Contains(output, "Bucket created successfully") {
		t.Errorf("Expected bucket creation message, but got: %s", output)
	}
}
