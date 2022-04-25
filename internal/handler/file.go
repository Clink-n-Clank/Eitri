package handler

import (
	"fmt"
	"io/ioutil"
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
