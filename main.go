package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	hostsFilePath  = "/etc/hosts"
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

func readWebsitesFromConfig() ([]string, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var websites []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		website := strings.TrimSpace(scanner.Text())
		if website != "" {
			websites = append(websites, website)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return websites, nil
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

func blockWebsites(websitesToBlock []string) error {
	file, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open /etc/hosts: %v", err)
	}
	defer file.Close()

	for _, website := range websitesToBlock {
		_, err := file.WriteString("127.0.0.1 " + website + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to /etc/hosts: %v", err)
		}
	}
	return nil
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

	websitesToBlock, err := readWebsitesFromConfig()
	if err != nil {
		fmt.Println("Failed to read config file: %v\n", err)
		return
	}

}
