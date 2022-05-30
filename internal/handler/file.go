package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const permissionFile = 0644

// CreateFile ...
func CreateFile(path, name, content string) error {
	if err := ioutil.WriteFile(filepath.Join(path, name), []byte(content), permissionFile); err != nil {
		return fmt.Errorf("unable to write file: %v", err)
	}

	return nil
}

// OverwriteToFile ...
func OverwriteToFile(filePath, content string) error {
	f, fWriteErr := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, permissionFile)
	if fWriteErr != nil {
		return fWriteErr
	}

	if _, fWriteErr = f.WriteString(content); fWriteErr != nil {
		return fmt.Errorf("unable to write file: %v", fWriteErr)
	}

	if fCloseErr := f.Close(); fCloseErr != nil {
		return fmt.Errorf("unable to write file: %v", fCloseErr)
	}

	return nil
}
