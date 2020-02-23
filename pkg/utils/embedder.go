package utils

import (
	"os"
)

// Embedder manages the script that is to be embedded
type Embedder struct {
	File string
}

// Init creates the file to be embedded if it does not yet exist
func (e *Embedder) Init() error {
	file, _ := os.Stat(e.File)
	if file == nil {
		if _, err := os.Create(e.File); err != nil {
			return err
		}
	}

	return nil
}

// Write writes to the file that is to be embedded
func (e *Embedder) Write(content string) error {
	embedFile, err := os.OpenFile(e.File, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	if _, err := embedFile.Write([]byte(content)); err != nil {
		return err
	}

	return embedFile.Close()
}
