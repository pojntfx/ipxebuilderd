package utils

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

// StatusCommandWorker is a command that returns a status
type StatusCommandWorker struct {
	Command []string
}

// Read reads from a reader
func (c *StatusCommandWorker) Read(reader io.Reader, outChan chan string, errChan chan error) {
	bufStdout := bufio.NewReader(reader)

	for {
		line, _, err := bufStdout.ReadLine()
		if err != nil {
			errChan <- err

			return
		}

		outChan <- string(line)
	}
}

// Run runs the command
func (c *StatusCommandWorker) Run(stdoutChan, stderrChan chan string, errChan chan error) {
	cmd := exec.Command(c.Command[0], c.Command[1:]...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	go c.Read(stdout, stdoutChan, errChan)
	go c.Read(stderr, stderrChan, errChan)

	errChan <- cmd.Run()

	close(stdoutChan)
	close(stderrChan)
	close(errChan)
}
