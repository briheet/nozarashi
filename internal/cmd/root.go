package cmd

import (
	"context"
	"log"

	"github.com/briheet/nozarashi/internal/cmd/system"
	"github.com/spf13/cobra"
)

// This is the main Execute function of the application.
// Everything flows from here.
func Execute(ctx context.Context) int {
	// This is the base command for nozarashi
	// All entry points, subcommands and such go through this
	rootCmd := &cobra.Command{
		Use:   "nozarashi",
		Short: "Nozarashi is a cli based application for orchestrating apple container.",
		RunE: func(cmd *cobra.Command, args []string) error {
			printAscii()
			return cmd.Help()
		},
	}

	// Register all top-level commands here.
	//
	// Nested commands should be registered on their parent command.
	// For example:
	//
	//	nozarashi up info
	//
	// `up` is a child of the root command, and `info` is a child of `up`.

	// This is the base system command wrapping over
	// apple container's cli for managing it
	rootCmd.AddCommand(system.SystemCmd(ctx))

	// This is the up command for starting the containers
	rootCmd.AddCommand(UpCmd(ctx))

	// Execute and return if any error
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error: %v", err)
		return -1
	}
	return 0
}
