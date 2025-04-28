package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
)

type dependencyPackage struct {
	DependencyString string `json:"dependency_string"`
}

type dependencies struct {
	Dependencies map[string]dependencyPackage `json:"dependencies"`
}

// findIndex finds the dependency string for a given key in the dependencies map
func (d *dependencies) findDependency(key string) (string, bool) {
	dep, exists := d.Dependencies[key]
	return dep.DependencyString, exists
}

// dependencyStatus checks whether dependencies are already installed and calls the appropriate functions to install them if they are not.

func DependencyStatus(dependency string) int {

	// Check if the dependency is already installed
	status := exec.Command(dependency, "--version")
	output, err := status.Output()
	if err != nil {
		log.Printf("%s is not installed \n", dependency)
		return getDependencies(dependency)
	}
	// If the dependency is already installed, prints the version number and returns status code 1
	log.Printf("%s is installed - version: %s \n", dependency, output)
	return 1
	// status code returns 1 if dependency is already installed, 0 if dependency is installed successfully, -1 if dependency installation fails
}

func getDependencies(dependency string) int {

	// Get the user's operating system
	userOs := runtime.GOOS

	// Open the dependency_id.json file
	jsonFile, JsonfileErr := os.Open("dependency_id.json")
	if JsonfileErr != nil {
		log.Println("Error reading dependency names from file")
		return -1
	}
	defer jsonFile.Close()

	// io reader converts value to byte array
	byteValue, _ := io.ReadAll(jsonFile)
	var packageList dependencies
	// JSON structure "unmarshalled" into struct of type packageList.dependencies
	JsonErr := json.Unmarshal(byteValue, &packageList)
	if JsonErr != nil {
		log.Println("Error reading dependency names from file")
		return -1
	}

	// match user os to find package installation moniker
	if userOs == "windows" || userOs == "darwin" || userOs == "linux" {
		jsonObject := fmt.Sprintf("%s_%s", dependency, userOs)
		dependencyId, exists := packageList.findDependency(jsonObject)
		if !exists {
			log.Println("Dependency not found in list")
			return -1
		}
		// call installDependency function to install the dependency and return status code
		log.Println("Calling fucntion to install dependency")
		return installDependency(userOs, dependencyId)
	} else {
		log.Println("Unsupported OS")
		return -1
	}

}

func installDependency(userOs string, dependency string) int {
	var cmd *exec.Cmd

	if dependency != "Automattic.Wordpress" {
		// switch statement to install dependencies based on user's package manager
		switch userOs {
		case "windows":
			// installs dependency with silent flag and accepts package/source agreements
			log.Println("Attempting to install dependency using winget, please follow any prompts")
			cmd = exec.Command("winget", "install", "-e", "--id", dependency, "--silent", "--accept-package-agreements", "--accept-source-agreements")
		case "linux":
			log.Println("Attempting to install dependency using apt, please follow any prompts")
			cmd = exec.Command("sudo", "apt", "install", "-y", dependency)
		case "darwin":
			log.Println("Attempting to install dependency using homebrew, please follow any prompts")
			cmd = exec.Command("brew", "install", dependency)
		default:
			log.Printf("Unsupported OS: %s \n", userOs)
			return -1
		}
	} else {
		// installs wordpress using curl
		getDependencies("php")
		log.Println("Attempting to install Wordpress using curl, please follow any prompts")
		url := "https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar"
		output := "wp-cli.phar"
		switch userOs {
		case "windows":
			cmd = exec.Command("powershell", "-Command", "curl", "-o", output, url)
		default:
			cmd = exec.Command("curl", "-o", output, url)
		}
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Error installing %s: %v\n", dependency, err)
		return -1
	}

	log.Printf("%s installed successfully \n", dependency)
	return 0
}
