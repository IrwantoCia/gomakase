package command

import (
	"os"
	"os/exec"
)

type Command interface {
	GoModTidy() error
	ChangeFolder(path string) error
	NPMInstall() error
	GoFmt() error
}

type command struct {
}

func NewCommand() Command {
	return &command{}
}

func (c *command) GoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *command) ChangeFolder(path string) error {
	return os.Chdir(path)
}

func (c *command) NPMInstall() error {
	cmd := exec.Command("npm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *command) GoFmt() error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
