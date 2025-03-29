package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func BuildProject(tech string, projectName string) {

	dependencyCall(tech, projectName)

}

func createDirectory(projectName string) string {

	/*if err != nil {
		fmt.Println("error determining root directory")
		os.Exit(1)
	}*/

	var root string
	if runtime.GOOS == "windows" {
		root = os.Getenv("SystemDrive")
		root += "\\"
	} else {
		root = "/"
	}
	directoryErr := os.Mkdir(filepath.Join(root, projectName), 0755)
	if directoryErr != nil {
		fmt.Println("error creating project directory")
		fmt.Println(directoryErr)
		os.Exit(1)
	}

	fmt.Printf("Project %s Directory created successfully, attempting to populate\n", projectName)
	wd := filepath.Join(root, projectName)
	return wd

}

func populateProject(tech string, projectName string, wd string) {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error determining working directory")
		os.Exit(1)
	}
	if pwd != wd {
		os.Chdir(wd)
	}

	cmd := exec.Command("git", "init")
	gitErr := cmd.Run()
	if gitErr != nil {
		fmt.Println("error initializing git repository - leaving for now")
	}

	switch tech {
	case "node":
		nodeConfig(wd)

	case "python":

		pyConfig(wd, projectName)

	case "wordpress":

		wpConfig(wd)

	case "php":

		phpConfig(wd)

	}

}

func dependencyCall(tech string, projectName string) {

	techStatus := DependencyStatus(tech)
	switch techStatus {
	case 0:
		fmt.Println("Technology Successfully Installed")
		wd := createDirectory(projectName)
		populateProject(tech, projectName, wd)
	case 1:
		fmt.Println("Technology Found")
		wd := createDirectory(projectName)
		populateProject(tech, projectName, wd)
	case -1:
		os.Exit(1)
	}
}

func nodeConfig(wd string) {

	fmt.Println("Node Configuration Begun")

	cmd := exec.Command("npm", "init", "-y")
	cmd.Dir = wd //  executes in the current directory
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
	}

	cmd = exec.Command("npm", "install", "express")
	cmd.Dir = wd
	_, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Node Configuration Complete")

}

func pyConfig(wd string, projectName string) {

	fmt.Println("Python Configuration Begun")

	pwd, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	pwd = filepath.Dir(pwd)

	script := filepath.Join(pwd, "pySetup.sh")
	newScript := filepath.Join((os.Getenv("SystemDrive") + "\\"), projectName, "pySetup.sh")
	err = os.Link(script, newScript)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("powershell", "-Command", "cd", wd, ";", "bash", "pySetup.sh")
	cmd.Dir = wd

	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Python Configuration Complete")
}

func wpConfig(wd string) {

	fmt.Println("Wordpress Configuration Begun")

	cmd := exec.Command("wp", "core", "download")
	cmd.Dir = wd //  executes in the current directory
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
	}

	cmd = exec.Command("npm", "install", "express")
	cmd.Dir = wd
	_, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Wordpress Configuration Complete")
}

func phpConfig(wd string) {
	fmt.Println("PHP Configuration Begun")
	setupScript := "scripts\\phpSetup.sh"
	cmd := exec.Command("bash", setupScript)
	cmd.Dir = wd //  executes in the current directory
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
	}

	cmd = exec.Command("npm", "install", "express")
	cmd.Dir = wd
	_, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("PHP Configuration Complete")
}
