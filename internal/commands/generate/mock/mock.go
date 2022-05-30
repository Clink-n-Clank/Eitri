package mock

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/Clink-n-Clank/Eitri/internal/handler"
	"github.com/Clink-n-Clank/Eitri/internal/log"
	"github.com/spf13/cobra"
)

const (
	mockGenCfgFile                  = "gomockhandler.json"
	mockGenFlagConfigPath           = "config"
	mockGenFlagConfigPathDefaultVal = ""
)

var (
	// GenerateMockCMD represents the new command.
	GenerateMockCMD = &cobra.Command{
		Use:     "mock",
		Short:   "Generate Mock: Generated mock code from config file: " + mockGenCfgFile,
		Long:    "Generate Mock: Will create and use mock manager to generate mock files from configuration " + mockGenCfgFile,
		Example: "Eitri mock \n Eitri mock --config cfg/" + mockGenCfgFile,
		Run:     runMockGen,
	}

	// templateMockConfig for mock management
	templateMockConfig = `{
	"mocks": {
		"MyMockFile": {
			"checksum": "AUTO_GENERATED",
			"source_checksum": "AUTO_GENERATED",
			"mode": "SOURCE_MODE or REFLECT_MODE",
			"source_mode_runner": {
				"source": "where/is/my_file.go",
				"interfaces": "option_if_i_need_it",
				"destination": "internal/generated/mocks/my_file.go"
			}
		}
	}
}`
)

func init() {
	GenerateMockCMD.PersistentFlags().String(
		mockGenFlagConfigPath,
		mockGenFlagConfigPathDefaultVal,
		"Custom location where mock config "+mockGenCfgFile+" is located",
	)
}

func runMockGen(c *cobra.Command, args []string) {
	if pErr := c.ParseFlags(args); pErr != nil {
		log.PrintError(fmt.Sprintf("failed to parse command flags, error: %s", pErr.Error()))
		return
	}

	// Try resolve location of mock config
	var basePathToConfig string
	cPath, isPathErr := handler.FindCurrentFolder()
	if isPathErr {
		return
	}

	if v, _ := c.Flags().GetString(mockGenFlagConfigPath); v != mockGenFlagConfigPathDefaultVal {
		basePathToConfig = path.Clean(path.Join(cPath, v))
	}
	if len(basePathToConfig) == 0 {
		v, mockCfgErr := handler.FindFilePath(cPath, mockGenCfgFile, []string{"vendor", "cmd"})
		if resolveError(mockCfgErr, mockGenCfgFile, cPath) {
			return
		}

		basePathToConfig = path.Clean(v)
	}

	log.PrintInfo("Managing project mocks...")

	cmd := exec.Command("gomockhandler", fmt.Sprintf("-config=%s", basePathToConfig), "mockgen")
	cmd.Dir = cPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.PrintError(err.Error())
	}

	log.PrintInfo("Done...")
}

func resolveError(mockCfgErr error, mockGenCfgFile string, rootFolder string) bool {
	findErr, isMockErrorType := mockCfgErr.(*handler.FindError)
	if !isMockErrorType && mockCfgErr != nil {
		log.PrintError(fmt.Sprintf("unable to find %s file, error: %s", mockGenCfgFile, mockCfgErr.Error()))

		return true
	}

	if isMockErrorType {
		if findErr.IsNotFound() {
			if cFile := handler.CreateFile(rootFolder, mockGenCfgFile, templateMockConfig); cFile != nil {
				log.PrintInfo(fmt.Sprintf(
					"unable to create mock configuration file, path: (%s) file (%s)",
					rootFolder,
					mockGenCfgFile,
				))
			}
		} else {
			log.PrintError(fmt.Sprintf("unable to find %s file, error: %s", mockGenCfgFile, mockCfgErr.Error()))

			return true
		}
	}

	return false
}
