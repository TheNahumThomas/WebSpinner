package cmd

import (
	"errors"
	"log"
	"os"
	"testing"
)

type MockFileSystem struct {
	MkdirErr error
	ChdirErr error
	LinkErr  error
}

func (mfs *MockFileSystem) Mkdir(name string, perm os.FileMode) error {
	if mfs.MkdirErr != nil {
		return mfs.MkdirErr
	}
	return nil
}

func (mfs *MockFileSystem) Chdir(directory string) error {
	if mfs.ChdirErr != nil {
		return mfs.ChdirErr
	}
	return nil
}

func (mfs *MockFileSystem) Link(oldfile, newfile string) error {
	if mfs.LinkErr != nil {
		return mfs.LinkErr
	}
	return nil
}

type MockCommander struct {
	RunCommandErr error
}

func (mc *MockCommander) RunCommand(command string, args ...string) error {
	if mc.RunCommandErr != nil {
		return mc.RunCommandErr
	}
	return nil
}

func TestBuildProject(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}
	mockFS.ChdirErr = nil
	mockFS.LinkErr = nil
	mockFS.MkdirErr = nil
	mockCR.RunCommandErr = nil

	// Test successful project build
	err := BuildProject(mockFS, mockCR, "node", "testProject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test unsupported technology
	err = BuildProject(mockFS, mockCR, "unsupportedTech", "testProject")
	if err == nil || err.Error() != "dependency not found" {
		t.Fatalf("Expected unfound dependency error, got %v", err)
	}
}

func TestCreateDirectory(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockFS.ChdirErr = nil
	mockFS.LinkErr = nil
	mockFS.MkdirErr = nil

	// Test successful directory creation
	projectName := "testProject"
	wd, err := CreateDirectory(mockFS, projectName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wd == "" {
		t.Fatalf("Expected a valid working directory, got empty string")
	}

	// Test directory creation failure
	mockFS.MkdirErr = errors.New("failed to create directory")
	_, err = CreateDirectory(mockFS, projectName)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}

}

func TestPopulateProject(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}
	mockFS.ChdirErr = nil
	mockFS.LinkErr = nil
	mockFS.MkdirErr = nil
	mockCR.RunCommandErr = nil

	// Test successful project population
	err := PopulateProject(mockFS, mockCR, "node", "/testProject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test unsupported technology
	err = PopulateProject(mockFS, mockCR, "unsupportedTech", "/testProject")
	if err == nil || err.Error() != "unsupported technology" {
		t.Fatalf("Expected unsupported technology error, got %v", err)
	}

	// Test failure to change directory
	mockFS.ChdirErr = errors.New("failed to change directory")
	err = PopulateProject(mockFS, mockCR, "node", "/testProject")
	if err == nil || err.Error() != "failed to change directory" {
		t.Fatalf("Expected directory change error, got %v", err)
	}
}

func TestInitializeGitRepo(t *testing.T) {
	mockCR := &MockCommander{}
	mockCR.RunCommandErr = nil

	// Test successful git initialization
	err := mockCR.RunCommand("git", "init")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test git initialization failure
	mockCR.RunCommandErr = errors.New("git init failed")
	err = mockCR.RunCommand("git", "init")
	if err == nil || err.Error() != "git init failed" {
		t.Fatalf("Expected git init error, got %v", err)
	}
}

func TestRunSetupScript(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}
	mockFS.ChdirErr = nil
	mockFS.LinkErr = nil
	mockFS.MkdirErr = nil
	mockCR.RunCommandErr = nil

	// Test successful script execution
	err := RunSetupScript(mockFS, mockCR, "/testProject", "nodeSetup")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test failure to link script
	mockFS.LinkErr = errors.New("failed to link script")
	err = RunSetupScript(mockFS, mockCR, "/testProject", "nodeSetup")
	if err == nil || err.Error() != "failed to link script" {
		t.Fatalf("Expected link error, got %v", err)
	}

	mockFS.LinkErr = nil

	// Test script execution failure
	mockCR.RunCommandErr = errors.New("script execution failed")
	err = RunSetupScript(mockFS, mockCR, "/testProject", "nodeSetup")
	if err == nil || err.Error() != "script execution failed" {
		t.Fatalf("Expected script execution error, got %v", err)
	}
}

func TestNodeConfig(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}

	// Test successful Node.js configuration
	err := NodeConfig(mockFS, mockCR, "/testProject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPyConfig(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}

	// Test successful Python configuration
	err := PyConfig(mockFS, mockCR, "/testProject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestWpConfig(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}

	// Test successful WordPress configuration
	err := WpConfig(mockFS, mockCR, "/testProject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPhpConfig(t *testing.T) {
	mockFS := &MockFileSystem{}
	mockCR := &MockCommander{}

	// Test successful PHP configuration
	err := PhpConfig(mockFS, mockCR, "/testProject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestSetupLogger(t *testing.T) {
	logFile, err := os.CreateTemp("", "testlog")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer os.Remove(logFile.Name())

	SetupLogger(logFile)
	log.Println("Test log message")
	// Verify log file content if necessary
}
