package cmd

import (
	"errors"
	"testing"
)

func TestDependencyStatus(t *testing.T) {
	mockCR := &MockCommander{}
	mockfh := &MockFileHandler{}

	// Test installed dependency
	status := DependencyStatus(mockCR, mockfh, "node")
	if status != 1 {
		t.Fatalf("Expected status 1 for installed dependency, got %d", status)
	}

	// Test missing dependency
	mockCR.RunCommandErr = errors.New("commandfailed")
	status = DependencyStatus(mockCR, mockfh, "unsupportedTech")
	if status != -1 {
		t.Fatalf("Expected status -1 for missing dependency, got %d", status)
	}
	mockCR.RunCommandErr = nil
}
