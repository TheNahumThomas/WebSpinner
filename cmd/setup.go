package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	internals "webspinner/internal"
)

// BuildProject is the entry point for the project setup process
func BuildProject(fs internals.FileSystem, cr internals.Commander, fh internals.FileHandler, tech string, projectName string) error {

	techStatus := DependencyStatus(cr, fh, tech)
	switch techStatus {
	case 0:
		log.Println("Dependency Successfully Installed")
	case 1:
		log.Println("Dependency Found")
	case -1:
		return errors.New("dependency not found")
	}

	wd, err := CreateDirectory(fs, projectName)
	if err != nil {
		return err
	}

	err = PopulateProject(fs, cr, tech, wd)
	if err != nil {
		return err
	}

	return nil
}

// createDirectory spins up a new directory for the project with the desired name at the root of the file system
func CreateDirectory(fs internals.FileSystem, projectName string) (string, error) {
	var root string
	if runtime.GOOS == "windows" {
		root = os.Getenv("SystemDrive") + "\\"
	} else {
		root = "/"
	}

	wd := filepath.Join(root, projectName)
	err := fs.Mkdir(wd, 0755)
	if err != nil {
		return "", err
	}

	log.Printf("Project %s Directory created successfully\n", projectName)
	return wd, nil
}

// populateProject populates the created directory by calling the config function for the selected webapp technology
func PopulateProject(fs internals.FileSystem, cr internals.Commander, tech string, wd string) error {

	err := fs.Chdir(wd)
	if err != nil {
		return err
	}

	_, err = cr.RunCommand("git", "init")
	if err != nil {
		log.Println("error initializing git repository:", err)
	}

	switch tech {
	case "node":
		return NodeConfig(fs, cr, wd)
	case "python":
		return PyConfig(fs, cr, wd)
	case "wordpress":
		return WpConfig(fs, cr, wd)
	case "php":
		return PhpConfig(fs, cr, wd)
	default:
		return errors.New("unsupported technology")
	}
}

// nodeConfig sets up a Node.js project
func NodeConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
	return RunSetupScript(fs, cr, wd, "nodeSetup")
}

// pyConfig sets up a Python project
func PyConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
	return RunSetupScript(fs, cr, wd, "pySetup")
}

// wpConfig sets up a WordPress project
func WpConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
	pwd, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		return err
	}
	pwd = filepath.Dir(pwd)
	wpCLI := filepath.Join(pwd, "wp-cli.phar")
	newWPCLI := filepath.Join(wd, "wp-cli.phar")
	err = fs.Link(wpCLI, newWPCLI)
	if err != nil {
		log.Println("Error linking wp-cli.phar:", err)
		return err
	}
	return RunSetupScript(fs, cr, wd, "wpSetup")
}

// phpConfig sets up a PHP project
func PhpConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
	return RunSetupScript(fs, cr, wd, "phpSetup")
}

// runSetupScript runs a setup script for the given technology, this used to be several individual but untestable functions
// they have been abstracted into this for the sake of DRY and testability
func RunSetupScript(fs internals.FileSystem, cr internals.Commander, wd, scriptName string) error {
	pwd, err := os.Executable()
	if err != nil {
		return err
	}
	pwd = filepath.Dir(pwd)
	keywordOne := "bash"
	var script, newScript string

	if scriptName == "nodeSetup" && runtime.GOOS == "windows" {
		script = filepath.Join(pwd, "scripts", scriptName+".bat")
		newScript = filepath.Join(wd, scriptName+".bat")
		keywordOne = "cmd"
	} else {
		script = filepath.Join(pwd, "scripts", scriptName+".sh")
		newScript = filepath.Join(wd, scriptName+".sh")
	}

	err = fs.Link(script, newScript)
	if err != nil {
		return err
	}

	log.Println("Setup script link created, running setup script:", newScript)
	if keywordOne == "cmd" {
		_, err = cr.RunCommand("cmd", "/C", filepath.Base(newScript))
		if err != nil {
			log.Println("Error running setup script:", err)
			return err
		}
	}
	_, err = cr.RunCommand("bash", filepath.Base(newScript))
	if err != nil {
		log.Println("Error running setup script:", err)
		return err
	}

	return err

}
