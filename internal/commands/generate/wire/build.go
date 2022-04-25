package wire

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Clink-n-Clank/Eitri/internal/handler"
	"github.com/Clink-n-Clank/Eitri/internal/log"
	"github.com/spf13/cobra"
)

const binName = "executable"

var (
	GenerateWireBinCMD = &cobra.Command{
		Use:   "build",
		Short: "Generate binary: Generates Google Wire DI code and compile it to binary",
		Run:   runCompile,
	}
)

// runCompile of main.go and wire it
func runCompile(_ *cobra.Command, _ []string) {
	entryFolder, hasErrToGetEntryFolder := handler.FindEntryPointPath()
	if hasErrToGetEntryFolder {
		return
	}

	if runWire(entryFolder) {
		runBuild(entryFolder)
	}
}

func runWire(entryFolder string) bool {
	log.PrintInfo("Compiling Google Wire DI...")

	cmd := exec.Command("wire")
	cmd.Dir = entryFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.PrintError(err.Error())

		return false
	}

	return true
}

func runBuild(entryFolder string) {
	binaryPath := fmt.Sprintf("%s/%s", entryFolder, binName)

	log.PrintInfo(fmt.Sprintf("Building binary: %s...", binaryPath))
	defer log.PrintInfo("Completed...")

	// TODO Add build flags later.
	bCMD := exec.Command("go", "build", "-o", binName)
	bCMD.Dir = entryFolder
	bCMD.Stdout = os.Stdout
	bCMD.Stderr = os.Stderr

	if err := bCMD.Run(); err != nil {
		log.PrintError(err.Error())
	}
}
