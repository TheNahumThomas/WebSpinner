package internals

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

type FileSystem interface {
	Mkdir(name string, perm os.FileMode) error
	Chdir(directory string) error
	Link(oldfile, newfile string) error
}

type Commander interface {
	RunCommand(command string, kwargs ...string) ([]byte, error)
}

type FileHandler interface {
	Open(name string) (io.ReadCloser, error)
}

type OSFileSys struct{}

func (fs OSFileSys) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (fs OSFileSys) Chdir(directory string) error {
	return os.Chdir(directory)
}

func (fs OSFileSys) Link(oldfile, newfile string) error {
	return os.Link(oldfile, newfile)
}

// OSFileHandler implements the FileHandler interface using os.Open
type OSFileHandler struct{}

func (fh OSFileHandler) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

type ExecCommander struct{}

func (ec ExecCommander) RunCommand(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)

	var outputBuffer bytes.Buffer

	Writer := io.Writer(&outputBuffer)
	cmd.Stdout = Writer
	cmd.Stderr = Writer // This redirects command output (such as setup script output) to the console while also logging it

	err := cmd.Run()

	log.Println(outputBuffer.String())

	return outputBuffer.Bytes(), err
}
