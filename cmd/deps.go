package cmd

import (
	"fmt"
	"os/exec"
)

// dependencyStatus checks whether dependencies are already installed and calls the appropriate functions to install them if they are not.

func dependencyStatus() {
	// Check if the node_modules directory exists
	nodeStatus := exec.Command("node", "-v")
	output, err := nodeStatus.Output()
	if err != nil {
		fmt.Println("Node is not installed")
		return
	}
	version := string(output)
	if version[2] < '2' {
		fmt.Printf("Outdated Node version detected (%s)\n", version)
	}
}
