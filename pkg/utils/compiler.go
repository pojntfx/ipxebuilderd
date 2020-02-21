package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Compiler is a compiler for iPXE
type Compiler struct {
	ExecPath            string
	statusCommandRunner StatusCommandRunner
}

// Init sets up the compiler
func (c *Compiler) Init() {
	c.statusCommandRunner = StatusCommandRunner{}
}

// Build builds iPXE
func (c *Compiler) Build(embedPath, platform, driver, extension string, stdoutChan, stderrChan chan string, doneChan chan string, errChan chan error) {
	path := platform + "/" + driver + "." + extension

	if err := os.Setenv("EMBED", embedPath); err != nil {
		errChan <- err

		return
	}

	cmd := exec.Command("make", path)
	cmd.Dir = c.ExecPath

	c.statusCommandRunner.Run(cmd, stdoutChan, stderrChan, errChan)

	doneChan <- filepath.Join(c.ExecPath, path)
}
