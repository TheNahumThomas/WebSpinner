package main

import (
	"flag"
	"fmt"
	"webspinner/cmd"
)

func main() {

	c := cmd.Config{}
	c.Setup()

	flag.Parse()

	fmt.Println(c.GetMessage())

	cmd.DependencyStatus("git")
	cmd.DependencyStatus("sqlite3")
	cmd.DependencyStatus("curl")
	cmd.BuildProject(c.Tech, c.Name)

}
