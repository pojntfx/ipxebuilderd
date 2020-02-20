package utils

import (
	"errors"
	"os"
)

// Embedder manages the script that is to be embedded
type Embedder struct {
	File      string
	embedFile *os.File
}

// Init creates the file to be embedded if it does not yet exist
func (e *Embedder) Init() error {
	file, _ := os.Stat(e.File)
	if file == nil {
		embedFile, err := os.Create(e.File)
		if err != nil {
			return err
		}

		e.embedFile = embedFile
	}

	return nil
}

// Write writes to the file that is to be embedded
func (e *Embedder) Write(content string) error {
	if e.embedFile == nil {
		return errors.New("file not yet set up")
	}

	_, err := e.embedFile.Write([]byte(content))

	return err
}
