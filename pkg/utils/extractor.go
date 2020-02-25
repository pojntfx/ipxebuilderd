package utils

//go:generate mkdir -p vendor/ipxe/dist
//go:generate rm -rf vendor/ipxe/src
//go:generate git clone https://github.com/pojntfx/ipxe.git vendor/ipxe/src
//go:generate arc -overwrite archive vendor/ipxe/dist/ipxe.tar.gz vendor/ipxe/src
//go:generate statik -src vendor/ipxe/dist

import (
	"os"

	"github.com/mholt/archiver"
	"github.com/rakyll/statik/fs"

	_ "github.com/pojntfx/ipxebuilderd/pkg/utils/statik" // The embedded assets.
)

// Extractor manages extraction of bundled assets.
type Extractor struct {
	BasePath, ArchivePath, ArchiveOutPath string
}

// Extract extracts the bundled iPXE source code.
func (e *Extractor) Extract() error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	data, err := fs.ReadFile(statikFS, "/ipxe.tar.gz")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(e.ArchivePath, 0777); err != nil {
		return err
	}

	archiveOut, err := os.Create(e.ArchiveOutPath)
	if err != nil {
		return err
	}

	if _, err := archiveOut.Write(data); err != nil {
		return err
	}

	return archiver.Unarchive(e.ArchiveOutPath, e.BasePath)
}
