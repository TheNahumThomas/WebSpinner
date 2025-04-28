package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"webspinner/cmd"
)

func main() {

	c := cmd.Config{}
	c.Setup()

	flag.Parse()

	fmt.Println(c.GetMessage())

	ld, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error determining home directory")
		os.Exit(1)
	}
	logFile, err := cmd.CreateLogFile(ld, c.Name)
	if err != nil {
		fmt.Println("error creating log file")
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.SetupLogger(logFile)
	log.Println("Logging Started")

	cmd.DependencyStatus("git")
	cmd.DependencyStatus("sqlite3")
	cmd.DependencyStatus("curl")
	cmd.BuildProject(c.Tech, c.Name)

}
