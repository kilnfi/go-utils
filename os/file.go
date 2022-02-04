package os

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(sourcePath, destPath string) error {
	_, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	err = os.MkdirAll(filepath.Dir(destPath), 0o700)
	if err != nil {
		return err
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}
