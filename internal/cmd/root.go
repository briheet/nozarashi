package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/briheet/nozarashi/internal/cmd/system"
	"github.com/spf13/cobra"
)

// This is the main Execute function of the application.
// Everything flows from here.
func Execute(ctx context.Context) int {
	var profile bool
	var cpuProfile *os.File

	// This is the base command for nozarashi
	// All entry points, subcommands and such go through this
	rootCmd := &cobra.Command{
		Use:                   "nozarashi",
		Short:                 "Nozarashi is a cli based application for orchestrating apple container.",
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			printAscii()
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			file, err := os.Create("cpu.pprof")
			if err != nil {
				return err
			}
			if err := pprof.StartCPUProfile(file); err != nil {
				return errors.Join(err, file.Close())
			}

			cpuProfile = file
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			pprof.StopCPUProfile()
			cpuCloseErr := cpuProfile.Close()

			memoryProfile, err := os.Create("mem.pprof")
			if err != nil {
				return errors.Join(cpuCloseErr, err)
			}

			runtime.GC()
			return errors.Join(
				cpuCloseErr,
				pprof.WriteHeapProfile(memoryProfile),
				memoryProfile.Close(),
			)
		},
	}

	// Profile any command and write both files to the current directory.
	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "record CPU and memory profiles")

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
	rootCmd.AddCommand(system.SystemCmd())

	// This is the create command for preparing project resources
	rootCmd.AddCommand(CreateCmd())

	// This is the up command for starting the containers
	rootCmd.AddCommand(UpCmd())

	// This is the down command for stopping the containers
	rootCmd.AddCommand(DownCmd())

	// This is the build command for preparing service images
	rootCmd.AddCommand(BuildCmd())

	// This is the logs command for getting logs of containers
	rootCmd.AddCommand(LogsCmd())

	// This is the ps command for listing running containers
	rootCmd.AddCommand(PsCmd())

	// This is the exec command for running commands in service containers
	rootCmd.AddCommand(ExecCmd())

	// This is the destroy command for deleting all project runtime objects
	rootCmd.AddCommand(DestroyCmd())

	// This is the restart command for restarting service containers
	rootCmd.AddCommand(RestartCmd())

	// This is the inspect command for displaying service container details
	rootCmd.AddCommand(InspectCmd())

	// This is the tui command for viewing running containers
	rootCmd.AddCommand(TuiCmd())

	// Execute and return if any error
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Printf("Error: %v", err)
		return -1
	}
	return 0
}
