package handler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Clink-n-Clank/Eitri/internal/log"
)

const (
	folderEntry = "cmd"
	fileEntry   = "root.go"

	FindErrorFileGeneral FindErrorType = iota
	FindErrorFileNotFound
	FindErrorFileDuplicate
)

// FindErrorType code
type FindErrorType byte

// FindError is an error type specific for search
type FindError struct {
	FindErrorType
	// Err message
	Err error
}

// Error to string
func (r *FindError) Error() string {
	return r.Err.Error()
}

// IsNotFound file
func (r FindError) IsNotFound() bool {
	return r.FindErrorType == FindErrorFileNotFound
}

// IsDuplicateFile file
func (r FindError) IsDuplicateFile() bool {
	return r.FindErrorType == FindErrorFileDuplicate
}

// FindFilePath looks for fileName in root and all it's subdirectories, returns path of fileName if found.
func FindFilePath(root, fileName string, excludeNames []string) (string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, n := range excludeNames {
			if strings.Contains(path, n) {
				return nil
			}
		}

		if info.Name() == fileName {
			files = append(files, path)
		}

		return nil
	})

	if len(files) == 0 {
		return "", &FindError{
			FindErrorType: FindErrorFileNotFound,
			Err:           fmt.Errorf("expected to find %s file in dir %s", fileName, root),
		}
	}
	if err != nil {
		return "", &FindError{
			FindErrorType: FindErrorFileGeneral,
			Err:           fmt.Errorf(fmt.Sprintf("failed to find %s in path: %s. Error: %v", fileName, root, err)),
		}
	}
	if len(files) != 1 {
		return "", &FindError{
			FindErrorType: FindErrorFileDuplicate,
			Err:           fmt.Errorf("expected to find exactly one %s file in dir %s, but found %d", fileName, root, len(files)),
		}
	}

	return files[0], nil
}

// FindCurrentFolder where command will be executed
func FindCurrentFolder() (wd string, hasErr bool) {
	wd, wdErr := os.Getwd()
	hasErr = wdErr != nil

	if hasErr {
		log.PrintError(wdErr.Error())

		return "", hasErr
	}

	return wd, hasErr
}

// FindEntryPointPath where main file and DI file located
func FindEntryPointPath() (string, bool) {
	wd, hasErr := FindCurrentFolder()
	if hasErr {
		return "", true
	}

	cmdPath := filepath.Join(wd, folderEntry)
	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		log.PrintError(fmt.Sprintf("cmd folder (%s) not exist or not found", cmdPath))

		return "", true
	}

	mainPath, err := FindFilePath(cmdPath, fileEntry, []string{})
	if hasErr {
		log.PrintError(err.Error())
		return "", true
	}

	return filepath.Dir(mainPath), false
}
