package cmd

import (
	"fmt"
	"os"
)

func BuildProject() {

}

func createDirectory() {

}

func createFiles(tech string) {

	techStatus := DependencyStatus(tech)
	switch techStatus {
	case 0:
		fmt.Println("Technology Successfully Installed")
	case 1:
		fmt.Println("Technology Found")
	case -1:
		os.Exit(1)
	}
}
