package main

import (
	"github.com/Clink-n-Clank/Eitri/internal/commands"
	"github.com/spf13/cobra"
	"os"
)

var (
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "eitri",
		Short: "Eitri Command-line toolsets.",
		Long:  `Eitri Command-line toolsets that help maintain code base and be productive.`,
	}
)

// init contains list of commands of ScanX framework
func init() {
	// Setup commands
	rootCmd.AddCommand(commands.InstallToolsCMD)
}

func main() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
