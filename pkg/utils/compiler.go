package utils

import (
	"os"
	"os/exec"
)

// Compiler is a compiler for iPXE
type Compiler struct {
	ExecPath string
}

// Build builds iPXE
func (c *Compiler) Build(embedPath, platform, driver, extension string) (string, error) {
	if err := os.Setenv("EMBED", embedPath); err != nil {
		return "", err
	}

	cmd := exec.Command("make", platform+"/"+driver+"."+extension)
	cmd.Dir = c.ExecPath
	out, err := cmd.CombinedOutput()

	return string(out), err
}
