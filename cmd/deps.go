package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	internals "webspinner/internal"
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

func DependencyStatus(cr internals.Commander, fh internals.FileHandler, dependency string) int {

	// Check if the dependency is already installed
	output, err := cr.RunCommand(dependency, "--version")
	outputString := string(output)
	if err != nil && outputString == "" {
		log.Printf("%s is not installed \n", dependency)
		return getDependencies(cr, fh, dependency)
	}
	// If the dependency is already installed, prints the version number and returns status code 1
	log.Printf("%s is installed - version: %s \n", dependency, output)
	return 1
	// status code returns 1 if dependency is already installed, 0 if dependency is installed successfully, -1 if dependency installation fails
}

func getDependencies(cr internals.Commander, fh internals.FileHandler, dependency string) int {

	// Get the user's operating system
	userOs := runtime.GOOS

	// Open the dependency_id.json file
	jsonFile, JsonfileErr := fh.Open("dependency_id.json")
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
		return installDependency(cr, fh, userOs, dependencyId)
	} else {
		log.Println("Unsupported OS")
		return -1
	}

}

func installDependency(cr internals.Commander, fh internals.FileHandler, userOs string, dependency string) int {

	if dependency != "Automattic.Wordpress" {
		// switch statement to install dependencies based on user's package manager
		switch userOs {
		case "windows":
			// installs dependency with silent flag and accepts package/source agreements
			log.Println("Attempting to install dependency using winget, please follow any prompts")
			_, err := cr.RunCommand("winget", "install", "-e", "--id", dependency, "--silent", "--accept-package-agreements", "--accept-source-agreements")
			if err != nil {
				log.Printf("Error installing %s: %v\n", dependency, err)
				return -1
			}

		case "linux":
			log.Println("Attempting to install dependency using apt, please follow any prompts")
			_, err := cr.RunCommand("sudo", "apt", "install", "-y", dependency)
			if err != nil {
				log.Printf("Error installing %s: %v\n", dependency, err)
				return -1
			}

		case "darwin":
			log.Println("Attempting to install dependency using homebrew, please follow any prompts")
			_, err := cr.RunCommand("brew", "install", dependency)
			if err != nil {
				log.Printf("Error installing %s: %v\n", dependency, err)
				return -1
			}

		default:
			log.Printf("Unsupported OS: %s \n", userOs)
			return -1
		}

	} else {
		// installs wordpress using curl
		getDependencies(cr, fh, "php")
		log.Println("Attempting to install Wordpress using curl, please follow any prompts")
		url := "https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar"
		output := "wp-cli.phar"
		switch userOs {
		case "windows":
			_, err := cr.RunCommand("powershell", "-Command", "curl", "-o", output, url)
			if err != nil {
				log.Printf("Error installing %s: %v\n", dependency, err)
				return -1
			}

		default:
			_, err := cr.RunCommand("curl", "-o", output, url)
			if err != nil {
				log.Printf("Error installing %s: %v\n", dependency, err)
				return -1
			}

		}
	}

	log.Printf("%s installed successfully \n", dependency)
	return 0

}
