package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pojntfx/ipxebuilderd/pkg/utils"
	uuid "github.com/satori/go.uuid"
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
	extractor := utils.Extractor{
		BasePath:       basePath,
		ArchivePath:    archivePath,
		ArchiveOutPath: archiveOutPath,
	}

	if err := extractor.Extract(); err != nil {
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
