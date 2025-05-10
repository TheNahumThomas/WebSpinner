package cmd

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// BuildProject is the entry point for the project setup process
func BuildProject(tech string, projectName string) error {
	techStatus := DependencyStatus(tech)
	if techStatus == -1 {
		return errors.New("dependency not found")
	}

	wd, err := createDirectory(projectName)
	if err != nil {
		return err
	}

	err = populateProject(tech, wd)
	if err != nil {
		return err
	}

	return nil
}

// createDirectory creates a new directory for the project
func createDirectory(projectName string) (string, error) {
	var root string
	if runtime.GOOS == "windows" {
		root = os.Getenv("SystemDrive") + "\\"
	} else {
		root = "/"
	}

	wd := filepath.Join(root, projectName)
	err := os.Mkdir(wd, 0755)
	if err != nil {
		return "", err
	}

	log.Printf("Project %s Directory created successfully\n", projectName)
	return wd, nil
}

// populateProject initializes the project with the selected technology
func populateProject(tech string, wd string) error {
	err := os.Chdir(wd)
	if err != nil {
		return err
	}

	err = initializeGitRepo()
	if err != nil {
		log.Println("error initializing git repository:", err)
	}

	switch tech {
	case "node":
		return nodeConfig(wd)
	case "python":
		return pyConfig(wd)
	case "wordpress":
		return wpConfig(wd)
	case "php":
		return phpConfig(wd)
	default:
		return errors.New("unsupported technology")
	}
}

// initializeGitRepo initializes a git repository in the current directory
func initializeGitRepo() error {
	cmd := exec.Command("git", "init")
	return cmd.Run()
}

// nodeConfig sets up a Node.js project
func nodeConfig(wd string) error {
	return runSetupScript(wd, "nodeSetup")
}

// pyConfig sets up a Python project
func pyConfig(wd string) error {
	return runSetupScript(wd, "pySetup")
}

// wpConfig sets up a WordPress project
func wpConfig(wd string) error {
	err := linkFile("wp-cli.phar", wd)
	if err != nil {
		return err
	}
	return runSetupScript(wd, "wpSetup")
}

// phpConfig sets up a PHP project
func phpConfig(wd string) error {
	return runSetupScript(wd, "phpSetup")
}

// runSetupScript runs a setup script for the given technology
func runSetupScript(wd, scriptName string) error {
	pwd, err := os.Executable()
	if err != nil {
		return err
	}
	pwd = filepath.Dir(pwd)
	newScript := filepath.Join(wd, scriptName)
	if scriptName == "nodeSetup" && runtime.GOOS == "windows" {
		script := filepath.Join(pwd, "scripts", scriptName+".bat")
		newScript = filepath.Join(wd, scriptName+".bat")
		err = linkFile(script, newScript)
		if err != nil {
			return err
		}
	} else {
		script := filepath.Join(pwd, "scripts", scriptName+".sh")
		newScript = filepath.Join(wd, scriptName+".sh")
		err = linkFile(script, newScript)
		if err != nil {
			return err
		}
	}
	cmd := exec.Command("bash", newScript)
	cmd.Dir = wd
	_, err = cmd.CombinedOutput()
	return err
}

// linkFile creates a hard link between two files
func linkFile(src, dest string) error {
	return os.Link(src, dest)
}
