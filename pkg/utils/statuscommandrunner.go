package utils

import (
	"bufio"
	"io"
	"os/exec"
)

// StatusCommandRunner is a command that returns a status.
type StatusCommandRunner struct {
}

func (c *StatusCommandRunner) read(reader io.Reader, outChan chan string) {
	bufStdout := bufio.NewReader(reader)

	for {
		line, _, err := bufStdout.ReadLine()
		if err != nil {
			return // Don't handle other errors, this just means EOF.
		}

		outChan <- string(line)
	}
}

// Run runs the command.
func (c *StatusCommandRunner) Run(cmd *exec.Cmd, stdoutChan, stderrChan chan string, errChan chan error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err

		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errChan <- err

		return
	}

	go c.read(stdout, stdoutChan)
	go c.read(stderr, stderrChan)

	if err := cmd.Run(); err != nil {
		errChan <- err
	}
}
