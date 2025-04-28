package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Flags parsed from command line

type Config struct {
	Tech string // Technology to use, no default"
	Name string // Name of project, defaults to "MyNewWebApp"
	O    bool   // Object-Oriented option flag
}

// setup function creates commnad line arguments as type Config
func (c *Config) Setup() {
	flag.StringVar(&c.Tech, "tech", "", "Selected WebApp Technology, must be specified [No Default Option]")
	flag.StringVar(&c.Name, "name", "MyNewWebApp", "Project Name, defaults to 'MyNewWebApp'")
	flag.BoolVar(&c.O, "o", false, "Option to specify Object-Oriented project, False by default")
	flag.Parse()
	// Check if technology is specified

	if c.Tech == "" {
		log.Println("err: no technology specified")
		flag.Usage()
		os.Exit(1)
	}
}

// GetMessage function returns a string repeating user selection
func (c *Config) GetMessage() string {
	projectType := "basic"
	if c.O {
		projectType = "object-oriented"
	}
	msg := fmt.Sprintf("Creating a new %s project with the name %s", projectType, c.Tech)
	return msg
}
