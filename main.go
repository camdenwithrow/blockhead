package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	configName     = "blockhead"
	configFileName = "blockhead.conf"
)

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".config", configName, configFileName), nil
}

func editConfigFile() error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Ensure the config directory exists
	configDir := filepath.Dir(configFilePath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	// Ensure the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		file, err := os.Create(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
		file.Close()
	}

	cmd := exec.Command("vi", configFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: blockhead <block|unblock|edit>")
		return
	}

	action := os.Args[1]

	if action == "edit" {
		err := editConfigFile()
		if err != nil {
			fmt.Printf("Failed to edit config file: %v\n", err)
		}
		return
	}
}
