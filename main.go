package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/pojntfx/ipxebuilderd/pkg/workers"
	uuid "github.com/satori/go.uuid"
)

var (
	platform  string
	driver    string
	extension string
	script    string
)

const (
	totalCompileSteps = 2247
)

func main() {
	flag.StringVar(&platform, "platform", "bin-x86_64-efi", "The platform to build iPXE for")
	flag.StringVar(&driver, "driver", "ipxe", "The driver to build iPXE for")
	flag.StringVar(&extension, "extension", "efi", "The extension to build iPXE for")
	flag.StringVar(&script, "script", `#!ipxe

autoboot`, "The script to embed")

	flag.Parse()

	builder := workers.Builder{
		BasePath: filepath.Join(os.TempDir(), "ipxebuilderd", uuid.NewV4().String()),
	}

	if err := builder.Extract(); err != nil {
		log.Fatal(err)
	}

	stdoutChan, stderrChan, doneChan, errChan := make(chan string), make(chan string), make(chan string), make(chan error)

	go builder.Build(script, platform, driver, extension, stdoutChan, stderrChan, doneChan, errChan)

	log.Println("Building iPXE")

	bar := pb.Full.Start(totalCompileSteps)

	var outPath string

Main:
	for {
		select {
		case <-stdoutChan:
			bar.Increment()
		case err := <-stderrChan: // This needs to be handled manually; `ar` prints to stdout for some reason
			if !strings.Contains(err, "ar:") {
				bar.Finish()
				log.Fatal("Fatal Error:", err)
				break Main
			}
			bar.Increment()
		case outPath = <-doneChan:
			bar.Finish()
			break Main
		case err := <-errChan:
			bar.Finish()
			log.Fatal("Fatal Error:", err)
		}
	}

	log.Println("iPXE out path:", outPath)
}
