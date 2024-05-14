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
		return nil, fmt.Errorf("error reading config file: %v", err)
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

func runAsAdmin(args []string) error {
    cmdPath, err := os.Executable()
    if err != nil {
        return fmt.Errorf("failed to get executable path: %v", err)
    }

    cmd := fmt.Sprintf(`"%s" %s`, cmdPath, strings.Join(args, " "))
    osascript := fmt.Sprintf(`do shell script "%s" with administrator privileges`, cmd)
    _, err = exec.Command("osascript", "-e", osascript).Output()
    return err
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

func unblockWebsites(websitesToUnblock []string) error {
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
		fmt.Printf("Failed to read config file: %v\n", err)
		return
	}

	if len(os.Args) > 2 && os.Args[2] == "elevated" {
		switch action {
		case "block":
			err = blockWebsites(websitesToBlock)
		case "unblock":
			err = unblockWebsites(websitesToBlock)
		default:
			fmt.Println("Unknown command: ", action)
			return
		}
		if err != nil {
			fmt.Printf("Failed to %s websites: %v\n", action, err)
		} else {
			fmt.Printf("Websites %sed successfully.\n", action)
		}
	} else {
		err := runAsAdmin(append([]string{action}, "elevated"))
		if err != nil {
			fmt.Printf("Failed to run as admin: %v\n", err)
		}
	}

}
