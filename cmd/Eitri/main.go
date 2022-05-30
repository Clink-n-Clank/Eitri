package main

import (
	"os"

	"github.com/Clink-n-Clank/Eitri/internal/commands/generate/mock"
	"github.com/Clink-n-Clank/Eitri/internal/commands/generate/proto"
	"github.com/Clink-n-Clank/Eitri/internal/commands/generate/wire"
	"github.com/Clink-n-Clank/Eitri/internal/commands/setup"
	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "Eitri",
		Short: "Eitri Command-line toolsets.",
		Long:  `Eitri Command-line toolsets that help maintain code base and be productive.`,
	}
)

// init contains list of commands of ScanX framework
func init() {
	// Setup commands
	rootCmd.AddCommand(setup.InstallToolsCMD)

	// Generate commands
	rootCmd.AddCommand(mock.GenerateMockCMD)
	rootCmd.AddCommand(wire.GenerateWireBinCMD)
	rootCmd.AddCommand(proto.GenerateEnvoyTranscodingCMD)
}

func main() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
