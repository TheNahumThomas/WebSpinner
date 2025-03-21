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

}
