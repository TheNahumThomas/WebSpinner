package cmd

import (
	"flag"
	"fmt"
	"os"
)

// Flags parsed from command line

type Config struct {
	tech string // Technology to use, no default"
	name string // Name of project, defaults to "MyNewWebApp"
	o    bool   // Object-Oriented option flag
}

// setup function creates commnad line arguments as type Config
func (c *Config) Setup() {
	flag.StringVar(&c.tech, "tech", "", "Selected WebApp Technology, must be specified [No Default Option]")
	flag.StringVar(&c.name, "name", "MyNewWebApp", "Project Name, defaults to 'MyNewWebApp'")
	flag.BoolVar(&c.o, "o", false, "Option to specify Object-Oriented project, False by default")
	flag.Parse()
	// Check if technology is specified

	if c.tech == "" {
		fmt.Println("err: no technology specified")
		flag.Usage()
		os.Exit(1)
	}
}

// GetMessage function returns a string repeating user selection
func (c *Config) GetMessage() string {
	projectType := "basic"
	if c.o {
		projectType = "object-oriented"
	}
	msg := fmt.Sprintf("Creating a new %s project with the name %s", projectType, c.name)
	return msg
}
