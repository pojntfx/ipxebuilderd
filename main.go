//go:generate mkdir -p vendor/ipxe/dist
//go:generate rm -rf vendor/ipxe/src
//go:generate git clone https://github.com/pojntfx/ipxe.git vendor/ipxe/src
//go:generate arc -overwrite archive vendor/ipxe/dist/ipxe.tar.gz vendor/ipxe/src
//go:generate statik -src vendor/ipxe/dist

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	"github.com/rakyll/statik/fs"
	uuid "github.com/satori/go.uuid"

	_ "github.com/pojntfx/ipxebuilderd/statik"
)

var (
	basePath       = filepath.Join(os.TempDir(), "ipxebuilderd", uuid.NewV4().String())
	archivePath    = filepath.Join(basePath, "archive")
	execPath       = filepath.Join(basePath, "src", "src")
	embedPath      = filepath.Join(execPath, "main.ipxe")
	archiveOutPath = filepath.Join(archivePath, "ipxe.tar.gz")

	platform  = "bin-x86_64-efi"
	driver    = "ipxe"
	extension = "efi"

	embedScript = `#!ipxe

autoboot`
)

func main() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	data, err := fs.ReadFile(statikFS, "/ipxe.tar.gz")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(archivePath, 0777); err != nil {
		log.Fatal(err)
	}

	archiveOut, err := os.Create(archiveOutPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := archiveOut.Write(data); err != nil {
		log.Fatal(err)
	}

	if err := archiver.Unarchive(archiveOutPath, basePath); err != nil {
		log.Fatal(err)
	}

	embedFile, err := os.Create(embedPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := embedFile.Write([]byte(embedScript)); err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv("EMBED", embedPath); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("make", platform+"/"+driver+"."+extension)
	cmd.Dir = execPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(string(out), err)
	}

	log.Println(string(out))
}
