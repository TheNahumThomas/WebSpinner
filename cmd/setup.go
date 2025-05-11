package cmd

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type FileSystem interface {
	Mkdir(name string, perm os.FileMode) error
	Chdir(directory string) error
	Link(oldfile, newfile string) error
}

type Commander interface {
	RunCommand(command string, kwargs ...string) error
}

type OSFileSys struct{}

func (fs OSFileSys) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (fs OSFileSys) Chdir(directory string) error {
	return os.Chdir(directory)
}

func (fs OSFileSys) Link(oldfile, newfile string) error {
	return os.Link(oldfile, newfile)
}

type ExecCommander struct{}

func (ec ExecCommander) RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer() // This redirects command output (such as setup script output) to the console while also logging it

	return cmd.Run()
}

// BuildProject is the entry point for the project setup process
func BuildProject(tech string, projectName string) error {

	fs := OSFileSys{}
	cr := ExecCommander{}

	techStatus := DependencyStatus(tech)
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
func CreateDirectory(fs FileSystem, projectName string) (string, error) {
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
func PopulateProject(fs FileSystem, cr Commander, tech string, wd string) error {
	err := fs.Chdir(wd)
	if err != nil {
		return err
	}

	err = InitializeGitRepo()
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

// initializeGitRepo initializes a git repository in the working directory
func InitializeGitRepo() error {
	cmd := exec.Command("git", "init")
	return cmd.Run()
}

// nodeConfig sets up a Node.js project
func NodeConfig(fs FileSystem, cr Commander, wd string) error {
	return RunSetupScript(fs, cr, wd, "nodeSetup")
}

// pyConfig sets up a Python project
func PyConfig(fs FileSystem, cr Commander, wd string) error {
	return RunSetupScript(fs, cr, wd, "pySetup")
}

// wpConfig sets up a WordPress project
func WpConfig(fs FileSystem, cr Commander, wd string) error {
	pwd, err := os.Executable()
	if err != nil {
		log.Println("Error getting executable path:", err)
		return err
	}
	pwd = filepath.Dir(pwd)
	wpCLI := filepath.Join(pwd, "wp-cli.phar")
	newWPCLI := filepath.Join(wd, "wp-cli.phar")
	err = LinkFile(fs, wpCLI, newWPCLI)
	if err != nil {
		log.Println("Error linking wp-cli.phar:", err)
		return err
	}
	return RunSetupScript(fs, cr, wd, "wpSetup")
}

// phpConfig sets up a PHP project
func PhpConfig(fs FileSystem, cr Commander, wd string) error {
	return RunSetupScript(fs, cr, wd, "phpSetup")
}

// runSetupScript runs a setup script for the given technology, this used to be several individual but untestable functions
// they have been abstracted into this for the sake of DRY and testability
func RunSetupScript(fs FileSystem, cr Commander, wd, scriptName string) error {
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

	err = LinkFile(fs, script, newScript)
	if err != nil {
		return err
	}

	log.Println("Setup script link created, running setup script:", newScript)
	if keywordOne == "cmd" {
		return cr.RunCommand("cmd", "/C", filepath.Base(newScript))
	}
	return cr.RunCommand("bash", filepath.Base(newScript))

}

// file link creates a hard link between two files
// this used to be part of several functions but has also been abstracted out for the sake of DRY and testability
func LinkFile(fs FileSystem, src, dest string) error {
	return fs.Link(src, dest)
}
