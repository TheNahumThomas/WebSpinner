## Documentation for WebSpinner

### Overview
WebSpinner is a Go-based CLI tool designed to automate the creation and setup of web application projects using various technologies, including Node.js, Python (Flask), PHP, and WordPress. It handles dependency installation, directory creation, and project initialization with secure configurations.

---

## Project Structure
```
WebSpinner/
├── cmd/
│   ├── setup.go          # Core logic for project setup
│   ├── deps.go           # Dependency management
│   ├── log.go            # Logging utilities
│   ├── utils.go          # Utility functions
│   ├── deps_test.go      # Unit tests for dependency management
│   ├── setup_test.go     # Unit tests for setup logic
│   └── test/             # Additional test files
├── internal/
│   └── interfaces.go     # Interfaces for file system and command execution
├── scripts/
│   ├── nodeSetup.bat     # Node.js setup script for Windows
│   ├── nodeSetup.sh      # Node.js setup script for Unix
│   ├── phpSetup.sh       # PHP setup script
│   ├── pySetup.sh        # Python setup script
│   └── wpSetup.sh        # WordPress setup script
├── dependency_id.json    # Dependency mapping for technologies
├── go.mod                # Go module file
├── main.go               # Entry point for the application
```

---

## Key Components

### 1. `BuildProject` Function
The `BuildProject` function is the main entry point for the project setup process. It performs the following steps:
1. Checks if the required dependencies for the selected technology are installed.
2. Creates a project directory.
3. Populates the directory with the necessary files and configurations.

```go
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
```

---

### 2. Directory Creation
The `CreateDirectory` function creates a new directory for the project at the root of the file system.

```go
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
```

---

### 3. Project Population
The `PopulateProject` function initializes the project directory with the necessary files and configurations based on the selected technology.

```go
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
```

---

### 4. Technology-Specific Configurations
Each supported technology has its own configuration function that sets up the project.

#### Node.js
```go
func NodeConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
    return RunSetupScript(fs, cr, wd, "nodeSetup")
}
```

#### Python
```go
func PyConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
    return RunSetupScript(fs, cr, wd, "pySetup")
}
```

#### WordPress
```go
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
```

#### PHP
```go
func PhpConfig(fs internals.FileSystem, cr internals.Commander, wd string) error {
    return RunSetupScript(fs, cr, wd, "phpSetup")
}
```

---

### 5. Running Setup Scripts
The `RunSetupScript` function executes the appropriate setup script for the selected technology.

```go
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
```

---

## Usage

### Command-Line Arguments
- `-tech`: Specifies the technology to use (e.g., `node`, `python`, `php`, `wordpress`).
- `-name`: Specifies the project name (default: `MyNewWebApp`).
- `-o`: Enables object-oriented project setup (default: `false`). NOTE: This hasn't been implemented as of yet and will do nothing if set

Example:
```sh
webspinner -tech=node -name=MyNodeApp (-o)
```

---

## Supported Technologies
- **Node.js**: Creates a secure Express.js application.
- **Python**: Creates a secure Flask application.
- **PHP**: Creates a basic PHP web application.
- **WordPress**: Sets up a secure WordPress installation.

---

## Logs
Logs are stored in a timestamped file in the user's home directory. The log file records all operations, including dependency checks, directory creation, and script execution.

---
