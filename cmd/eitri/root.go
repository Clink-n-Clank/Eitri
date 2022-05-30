package eitri

import (
	"github.com/Clink-n-Clank/Eitri/internal/commands/generate/mock"
	"github.com/Clink-n-Clank/Eitri/internal/commands/generate/proto"
	"github.com/Clink-n-Clank/Eitri/internal/commands/generate/wire"
	"github.com/Clink-n-Clank/Eitri/internal/commands/setup"
	"github.com/spf13/cobra"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "eitri",
		Short: "eitri Command-line toolsets.",
		Long:  `eitri Command-line toolsets that help maintain code base and be productive.`,
	}
)

// ExecCommandTool eitri commands
func ExecCommandTool() error {
	// Setup commands
	rootCmd.AddCommand(setup.InstallToolsCMD)

	// Generate commands
	rootCmd.AddCommand(mock.GenerateMockCMD)
	rootCmd.AddCommand(wire.GenerateWireBinCMD)
	rootCmd.AddCommand(proto.GenerateEnvoyTranscodingCMD)

	return rootCmd.Execute()
}
