package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// run "go run" command
func RunGo(target string, buildFlags []string) (string, error) {
	s, err := gocmd("run", target, buildFlags)
	if err != nil {
		return "", fmt.Errorf("taurus_go/cmd run error:\n%s", err)
	}
	return s, nil
}

// run "go list" command
func List(target string, buildFlags []string) error {
	_, err := gocmd("list", target, buildFlags)
	return err
}

func gocmd(command, target string, buildFlags []string) (string, error) {
	args := []string{command}
	args = append(args, buildFlags...)
	args = append(args, target)
	cmd := exec.Command("go", args...)
	stderr := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}
	return stdout.String(), nil
}
