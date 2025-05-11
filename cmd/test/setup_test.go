package cmd

import (
	"errors"
	"os"
	"testing"
)

type MockFileSystem struct {
	MkdirErr error
}

func (mfs MockFileSystem) Mkdir(name string, perm os.FileMode) error {
	return mfs.MkdirErr
}

func (mfs MockFileSystem) Chdir(dir string) error {
	return nil
}

func (mfs MockFileSystem) Link(oldname, newname string) error {
	return nil
}

func TestCreateDirectory(t *testing.T) {
	mockFS := MockFileSystem{}

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
