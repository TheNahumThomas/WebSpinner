package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Functino to create log file and corresponding, timestamped filname
func CreateLogFile(wd string, projectName string) (*os.File, error) {

	logFile, err := os.Create(filepath.Join(wd, fmt.Sprintf("%s_%s.log", projectName, time.Now().Format("2006-01-02_15-04-05"))))
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return nil, err
	}

	return logFile, nil
}

// SetupLogger configures the logger to send log messages to the console via stdout as well as the created log file through the io.MultiWriter package
func SetupLogger(logFile *os.File) {
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime)
}
