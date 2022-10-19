package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

type Command struct {
	Cmd    string
	Params map[string]string
}

func (comm *Command) Execute() error {
	scmd := comm.Cmd
	for k, v := range comm.Params {
		scmd += fmt.Sprintf(" --%s %s ", k, v)
	}
	cmd := exec.Command("bash", "-c", scmd)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}
	fmt.Println(string(stdout.Bytes()))
	return nil
}
