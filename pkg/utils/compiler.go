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
func (c *Compiler) Build(embedPath, platform, driver, extension string) (string, string, error) {
	path := platform + "/" + driver + "." + extension

	if err := os.Setenv("EMBED", embedPath); err != nil {
		return "", "", err
	}

	cmd := exec.Command("make", path)
	cmd.Dir = c.ExecPath
	out, err := cmd.CombinedOutput()

	return string(out), path, err
}
