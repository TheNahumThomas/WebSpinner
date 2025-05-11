package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"webspinner/cmd"
	internals "webspinner/internal"
)

func main() {

	fs := internals.OSFileSys{}
	cr := internals.ExecCommander{}

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
	log.Printf("Logging Started, Log file created at: %s", ld)

	cmd.DependencyStatus(cr, "git")
	cmd.DependencyStatus(cr, "sqlite3")
	cmd.DependencyStatus(cr, "curl")
	cmd.BuildProject(fs, cr, c.Tech, c.Name)

}
