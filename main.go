package main

import (
	"flag"
	"fmt"
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

	platform  string
	driver    string
	extension string
	script    string
)

func main() {
	flag.StringVar(&platform, "platform", "bin-x86_64-efi", "The platform to build iPXE for")
	flag.StringVar(&driver, "driver", "ipxe", "The driver to build iPXE for")
	flag.StringVar(&extension, "extension", "efi", "The extension to build iPXE for")
	flag.StringVar(&script, "script", `#!ipxe

autoboot`, "The script to embed")

	flag.Parse()

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

	if err := embedder.Write(script); err != nil {
		log.Fatal(err)
	}

	compiler := utils.Compiler{
		ExecPath: execPath,
	}

	stdoutChan, stderrChan, doneChan, errChan := make(chan string), make(chan string), make(chan string), make(chan error)

	go compiler.Build(embedPath, platform, driver, extension, stdoutChan, stderrChan, doneChan, errChan)

	var outPath string

Main:
	for {
		select {
		case outMsg := <-stdoutChan:
			log.Println("Status:", outMsg)
		case outErr := <-stderrChan:
			log.Println("Warning:", outErr)
		case outPath = <-doneChan:
			break Main
		case err := <-errChan:
			log.Fatal("Error:", err)
		}
	}

	fmt.Println("Done! Out path:", outPath)
}
