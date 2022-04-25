package setup

import (
	"github.com/Clink-n-Clank/Eitri/internal/handler"
	"github.com/Clink-n-Clank/Eitri/internal/log"
	"github.com/spf13/cobra"
)

// InstallToolsCMD represents the new command.
var InstallToolsCMD = &cobra.Command{
	Use:   "tools",
	Short: "Install tools: Will install depended tools for work",
	Run:   runInstallTools,
}

func runInstallTools(_ *cobra.Command, _ []string) {
	log.PrintInfo("Installing and Updating depended on tools for ScanX...")

	installWire()
	installMockHandler()
}

func installWire() {
	// wire is automated initialization of dependency injection created @Google https://go.dev/blog/wire
	const wire = "github.com/google/wire/cmd/wire@latest"

	// Install or Update command
	if err := handler.Install(wire); err != nil {
		log.PrintError(err.Error())
		return
	}
}

func installMockHandler() {
	const (
		mockGen     = "github.com/golang/mock/mockgen@latest"
		mockManager = "github.com/sanposhiho/gomockhandler@latest"
	)

	// Install or Update command
	if err := handler.Install(mockGen); err != nil {
		log.PrintError(err.Error())
		return
	}

	// Install or Update command
	if err := handler.Install(mockManager); err != nil {
		log.PrintError(err.Error())
		return
	}
}
