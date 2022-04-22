package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFilePath(t *testing.T) {
	stubFileName := "stub.txt"
	duplicateFileName := "duplicate.txt"
	relativeRootPath := "../.."
	testCases := []struct{
		name string
		root string
		fileName string
		expectedPath string
		expectedErr bool
	}{
		{
			name: "should return file path",
			root: relativeRootPath,
			fileName: stubFileName,
			expectedPath: fmt.Sprintf("%s/assets/tests/files/search/stub.txt", relativeRootPath),
		},
		{
			name: "should return duplicate error",
			root: relativeRootPath,
			fileName: duplicateFileName,
			expectedErr: true,
		},
		{
			name: "should return path error",
			root: "non_existing_path",
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			folder, hasErr := FindCurrentFolder()
			if hasErr {
				fmt.Println("there was an error")
			}
			fmt.Printf("This is current folder: %s", folder)
			path, err := FindFilePath(tc.root, tc.fileName)
			fmt.Printf("this is the path: %s", path)
			assert.Equal(t, tc.expectedPath, path)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
