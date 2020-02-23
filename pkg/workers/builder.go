package workers

import (
	"path/filepath"

	"github.com/pojntfx/ipxebuilderd/pkg/utils"
)

// Builder is a worker that can build iPXEs
type Builder struct {
	BasePath string
}

// Extract extracts the iPXE source code
func (b *Builder) Extract() error {
	archivePath := filepath.Join(b.BasePath, "archive")

	extractor := utils.Extractor{
		BasePath:       b.BasePath,
		ArchivePath:    archivePath,
		ArchiveOutPath: filepath.Join(archivePath, "ipxe.tar.gz"),
	}

	return extractor.Extract()
}

// Build buils an iPXE
func (b *Builder) Build(script, platform, driver, extension string, stdoutChan, stderrChan chan string, doneChan chan string, errChan chan error) {
	execPath := filepath.Join(b.BasePath, "src", "src")
	embedPath := filepath.Join(execPath, "main.ipxe")

	embedder := utils.Embedder{
		File: embedPath,
	}

	if err := embedder.Init(); err != nil {
		errChan <- err

		return
	}

	if err := embedder.Write(script); err != nil {
		errChan <- err

		return
	}

	compiler := utils.Compiler{
		ExecPath: execPath,
	}

	compiler.Build(embedPath, platform, driver, extension, stdoutChan, stderrChan, doneChan, errChan)

	return
}
