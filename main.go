package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	embedder := utils.Embedder{
		File: embedPath,
	}

	if err := embedder.Init(); err != nil {
		log.Fatal(err)
	}

	if err := embedder.Write(embedScript); err != nil {
		log.Fatal(err)
	}

	compiler := utils.Compiler{
		ExecPath: execPath,
	}

	out, outPath, err := compiler.Build(embedPath, platform, driver, extension)
	if err != nil {
		log.Fatal(out, outPath, err)
	}

	outFile, err := os.Open(outPath)
	if err != nil {
		log.Fatal(err)
	}

	product, err := ioutil.ReadAll(outFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(product)
}
