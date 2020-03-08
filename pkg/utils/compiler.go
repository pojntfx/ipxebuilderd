package utils

import (
	"os"
	"path/filepath"

	"github.com/pojntfx/dibs/pkg/utils"
)

// Compiler is a compiler for iPXE.
type Compiler struct {
	ExecPath string
}

// Build builds iPXE.
func (c *Compiler) Build(embedPath, platform, driver, extension string, stdoutChan, stderrChan chan string, doneChan chan string, errChan chan error) {
	path := platform + "/" + driver + "." + extension

	if err := os.Setenv("EMBED", embedPath); err != nil {
		errChan <- err

		return
	}

	cmd := utils.NewManageableCommand("make "+path, c.ExecPath, stdoutChan, stderrChan)

	if err := cmd.Start(); err != nil {
		errChan <- err

		return
	}

	if err := cmd.Wait(); err != nil {
		errChan <- err

		return
	}

	doneChan <- filepath.Join(c.ExecPath, path)
}
