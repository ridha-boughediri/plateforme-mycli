package libs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ridha-boughediri/plateforme-mycli/configs"
)

const aliasFileName = "alias.conf"

func GetAliasFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)
	aliasFilePath := filepath.Join(exeDir, aliasFileName)

	return aliasFilePath, nil
}

func FindAlias(alias string) (string, error) {
	aliasFile, err := GetAliasFilePath()
	if err != nil {
		return "", err
	}

	var aliases []configs.AliasConfig

	if _, err := os.Stat(aliasFile); err != nil {
		return "", fmt.Errorf("alias file not found")
	}

	fileData, err := os.ReadFile(aliasFile)
	if err != nil {
		return "", fmt.Errorf("could not read alias file: %v", err)
	}
	if err := json.Unmarshal(fileData, &aliases); err != nil {
		return "", fmt.Errorf("could not unmarshal alias file: %v", err)
	}

	for _, a := range aliases {
		if a.Alias == alias {
			return a.Remote, nil
		}
	}

	return "", fmt.Errorf("alias not found")
}

func UrlParts(url string) ([]string, error) {
	parts := strings.SplitN(url, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("alias and bucket should be in the correct format alias/bucket")
	}

	return parts, nil
}
