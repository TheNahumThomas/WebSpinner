package cmd

import (
	"encoding/json"
	"fmt"
	"io"
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

	// Check if the node_modules directory exists
	status := exec.Command("%s", "--version", dependency)
	output, err := status.Output()
	if err != nil {
		fmt.Printf("%s is not installed", dependency)
		return getDependencies(dependency)
	}

	fmt.Printf("%s is installed - version: %s", dependency, output)
	// status code returns 1 if dependency is already installed, 0 if dependency is installed successfully, -1 if dependency installation fails
	return 1
}

func getDependencies(dependency string) int {

	// Get the user's operating system
	userOs := runtime.GOOS

	// Open the dependency_id.json file
	jsonFile, JsonfileErr := os.Open("dependency_id.json")
	if JsonfileErr != nil {
		fmt.Println("Error reading dependency names from file")
		return -1
	}
	defer jsonFile.Close()

	// io reader converts value to byte array
	byteValue, _ := io.ReadAll(jsonFile)
	var packageList dependencies
	// JSON structure "unmarshalled" into struct of type packageList.dependencies
	JsonErr := json.Unmarshal(byteValue, &packageList)
	if JsonErr != nil {
		fmt.Println("Error reading dependency names from file")
		return -1
	}

	// match user os to find package installation moniker
	if userOs == "windows" || userOs == "darwin" || userOs == "linux" {
		jsonObject := fmt.Sprintf("%s_%s", dependency, userOs)
		fmt.Println(packageList.Dependencies)
		fmt.Println(jsonObject)
		dependencyId, exists := packageList.findDependency(jsonObject)
		if !exists {
			fmt.Println("Dependency not found in list")
			return -1
		}
		fmt.Println(dependencyId)
		// call installDependency function to install the dependency and return status code
		return installDependency(userOs, dependencyId)
	} else {
		fmt.Println("Unsupported OS")
		return -1
	}

}

func installDependency(userOs string, dependency string) int {

	switch userOs {
	// Windows installation using winget package manager
	case "windows":
		fmt.Println("Installing dependency using winget, the installer may ask for administrative permissions")
		cmd := exec.Command("winget", "install", "-e", "--id", dependency, "--silent", "--accept-package-agreements", "--accept-source-agreements")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error installing %s", dependency)
			return -1
		}
	// Linux (debian, ubuntu) installation using apt package manager with Super User permissions
	case "linux":
		cmd := exec.Command("sudo", "apt", "install", dependency)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error installing %s", dependency)
			return -1
		}
	// MacOS installation using Homebrew package manager
	case "darwin":
		cmd := exec.Command("brew", "install", dependency)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error installing %s", dependency)
			return -1
		}
	default:
		fmt.Printf("Command Line Arguments Couldn't be Established")
		return -1
	}

	fmt.Printf("%s installed successfully", dependency)
	return 0

}
