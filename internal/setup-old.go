package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// BuildProject is the entry point for the project setup process
func BuildProject(tech string, projectName string) {

	dependencyCall(tech, projectName)

}

// This function is private and finds the root directory for the user's file system and creates a new directory with the project name
func createDirectory(projectName string) string {

	var root string
	if runtime.GOOS == "windows" {
		root = os.Getenv("SystemDrive")
		root += "\\"
	} else {
		root = "/"
	}
	directoryErr := os.Mkdir(filepath.Join(root, projectName), 0755)
	if directoryErr != nil {
		log.Println("error creating project directory")
		log.Println(directoryErr)
		os.Exit(1)
	}

	log.Printf("Project %s Directory created successfully, attempting to populate\n", projectName)
	wd := filepath.Join(root, projectName)
	return wd

}

// This function is also private and takes the created working directory, initialises a git repository and launches the config process for the selected technology
func populateProject(tech string, wd string) {

	pwd, err := os.Getwd()
	if err != nil {
		log.Println("error determining working directory")
		os.Exit(1)
	}
	if pwd != wd {
		os.Chdir(wd)
	}

	cmd := exec.Command("git", "init")
	gitErr := cmd.Run()
	if gitErr != nil {
		log.Println("error initializing git repository - leaving for now")
	}

	switch tech {
	case "node":
		nodeConfig(wd)

	case "python":

		pyConfig(wd)

	case "wordpress":

		wpConfig(wd)
	case "php":

		phpConfig(wd)

	}

}

// dependency call is designed to pull in the pre-requisite dependencies for the selected webapp technology and return this to the createDirectory and populateProject functions
func dependencyCall(tech string, projectName string) {

	techStatus := DependencyStatus(tech)
	switch techStatus {
	case 0:
		log.Println("Technology Successfully Installed")
		wd := createDirectory(projectName)
		populateProject(tech, wd)
	case 1:
		log.Println("Technology Found")
		wd := createDirectory(projectName)
		populateProject(tech, wd)
	case -1:
		os.Exit(1)
	}
}

// The below private functions are called by the populateProject function and are designed to run the corresponding setup scripts for the selected technology.
// The scripts are located in the scripts directory and mostly consist of bash files with WordPress using the WP-CLI tool to build the project and Node using Batch on Windows.

func nodeConfig(wd string) {

	log.Println("Node.js Configuration Begun")

	pwd, err := os.Executable() // gets path to the WebSpinner executable
	if err != nil {
		log.Println(err)
	}
	pwd = filepath.Dir(pwd)

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" { // check for unix based OS

		script := filepath.Join(pwd, "scripts", "nodeSetup.sh")
		newScript := filepath.Join(wd, "nodeSetup.sh") // Creates new script file in the project directory
		err = os.Link(script, newScript)               // Creates a hard link to the script file (essentially copying its contents) so it can be run in the project directory
		if err != nil {
			log.Println(err)
		}
		log.Println("Setup script link created, running setup script")
		os.Chdir(wd)
		cmd := exec.Command("bash", "nodeSetup.sh")
		cmd.Dir = wd
		_, err = cmd.CombinedOutput()

	} else { // run batch file if on Windows

		script := filepath.Join(pwd, "scripts", "nodeSetup.bat")
		newScript := filepath.Join(wd, "nodeSetup.bat")
		err = os.Link(script, newScript)
		if err != nil {
			log.Println(err)
		}
		log.Println("Setup script link created, running setup script")
		os.Chdir(wd)
		cmd := exec.Command("cmd", "/C", "nodeSetup.bat")
		cmd.Dir = wd
		_, err = cmd.CombinedOutput()
	}

	if err != nil {
		log.Println("Error:", err)
	}

	log.Println("Node.js Configuration Complete")

}

func pyConfig(wd string) {

	log.Println("Python Configuration Begun")

	pwd, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	pwd = filepath.Dir(pwd)

	script := filepath.Join(pwd, "scripts", "pySetup.sh")
	newScript := filepath.Join(wd, "pySetup.sh")
	err = os.Link(script, newScript)
	if err != nil {
		log.Println(err)
	}
	log.Println("Setup script link created, running setup script")
	os.Chdir(wd)
	cmd := exec.Command("bash", "pySetup.sh")
	cmd.Dir = wd

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("Python Configuration Complete")
}

func wpConfig(wd string) {

	log.Println("Wordpress Configuration Begun")

	pwd, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	pwd = filepath.Dir(pwd)

	script := filepath.Join(pwd, "scripts", "wpSetup.sh")
	newScript := filepath.Join(wd, "wpSetup.sh")
	wpCli := filepath.Join(pwd, "wp-cli.phar")
	newWpCli := filepath.Join(wd, "wp-cli.phar")

	err = os.Link(wpCli, newWpCli) // Links the installed wordpress command line into the new project directory for use in the setup script so we don't have to globally install it
	if err != nil {
		log.Println(err)
	}
	err = os.Link(script, newScript)
	if err != nil {
		log.Println(err)
	}
	log.Println("Setup script link created, running setup script")
	os.Chdir(wd)
	cmd := exec.Command("bash", "wpSetup.sh")
	cmd.Dir = wd
	_, err = cmd.CombinedOutput()

	if err != nil {
		log.Println("Error:", err)
	}

	log.Println("Wordpress Configuration Complete")
}

func phpConfig(wd string) {

	log.Println("PHP Configuration Begun")

	pwd, err := os.Executable()
	if err != nil {
		log.Println(err)
	}
	pwd = filepath.Dir(pwd)

	script := filepath.Join(pwd, "scripts", "phpSetup.sh")
	newScript := filepath.Join(wd, "phpSetup.sh")
	err = os.Link(script, newScript)
	if err != nil {
		log.Println(err)
	}
	log.Println("Setup script link created, running setup script")
	os.Chdir(wd)
	cmd := exec.Command("bash", "phpSetup.sh")
	cmd.Dir = wd
	_, err = cmd.CombinedOutput()

	if err != nil {
		log.Println("Error:", err)
	}

	log.Println("PHP Configuration Complete")
}
