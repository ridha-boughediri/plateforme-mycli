package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ridha-boughediri/plateforme-mycli/configs"
	"github.com/ridha-boughediri/plateforme-mycli/libs"
)

func checkRemote(remote string) error {
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}

	resp, err := client.Get(remote)
	if err != nil {
		return fmt.Errorf("could not reach the remote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote responded with status: %d", resp.StatusCode)
	}

	return nil
}

func SaveAlias(aliasConfig configs.AliasConfig) error {

	if err := checkRemote(aliasConfig.Remote); err != nil {
		return fmt.Errorf("unable to connect to server: %v", err)
	}

	aliasFile, err := libs.GetAliasFilePath()
	if err != nil {
		return err
	}

	var aliases []configs.AliasConfig
	if _, err := os.Stat(aliasFile); err == nil {
		fileData, err := os.ReadFile(aliasFile)
		if err != nil {
			return fmt.Errorf("could not read alias file: %v", err)
		}
		json.Unmarshal(fileData, &aliases)

		for _, a := range aliases {
			if a.Alias == aliasConfig.Alias {
				return fmt.Errorf("alias already exists")
			}
		}
	}

	aliases = append(aliases, aliasConfig)

	fileData, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal aliases: %v", err)
	}
	if err := os.WriteFile(aliasFile, fileData, 0644); err != nil {
		return fmt.Errorf("could not write alias file: %v", err)
	}

	return nil
}

func ShowAliases() ([]configs.AliasConfig, error) {
	aliasFile, err := libs.GetAliasFilePath()
	if err != nil {
		return nil, err
	}

	var aliases []configs.AliasConfig

	if _, err := os.Stat(aliasFile); err == nil {
		fileData, err := os.ReadFile(aliasFile)
		if err != nil {
			return nil, fmt.Errorf("could not read alias file: %v", err)
		}
		json.Unmarshal(fileData, &aliases)
	} else {
		return nil, fmt.Errorf("alias file not found")
	}

	return aliases, nil
}

func DeleteAlias(alias string) error {
	aliasFile, err := libs.GetAliasFilePath()
	if err != nil {
		return err
	}

	var aliases []configs.AliasConfig

	if _, err := os.Stat(aliasFile); err == nil {
		fileData, err := os.ReadFile(aliasFile)
		if err != nil {
			return fmt.Errorf("could not read alias file: %v", err)
		}
		json.Unmarshal(fileData, &aliases)
	} else {
		return fmt.Errorf("alias file not found")
	}

	updatedAliases := []configs.AliasConfig{}
	for _, a := range aliases {
		if a.Alias != alias {
			updatedAliases = append(updatedAliases, a)
		}
	}

	fileData, err := json.MarshalIndent(updatedAliases, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal aliases: %v", err)
	}
	if err := os.WriteFile(aliasFile, fileData, 0644); err != nil {
		return fmt.Errorf("could not write alias file: %v", err)
	}

	return nil
}
