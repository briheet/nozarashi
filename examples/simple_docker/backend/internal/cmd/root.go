package cmd

import (
	"context"
	"errors"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {

	var profile bool
	var cpuProfile *os.File
	rootCmd := &cobra.Command{
		Use:   "backend",
		Short: "Simple backend for examples",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			f, perr := os.Create("cpu.pprof")
			if perr != nil {
				return perr
			}

			if err := pprof.StartCPUProfile(f); err != nil {
				return errors.Join(err, f.Close())
			}
			cpuProfile = f
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if !profile {
				return nil
			}

			pprof.StopCPUProfile()
			cpuCloseErr := cpuProfile.Close()

			f, perr := os.Create("mem.pprof")
			if perr != nil {
				return errors.Join(cpuCloseErr, perr)
			}

			runtime.GC()
			return errors.Join(cpuCloseErr, pprof.WriteHeapProfile(f), f.Close())
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&profile, "profile", "p", false, "record CPU and Mem pprof")

	rootCmd.AddCommand(ApiCmd(ctx))
	rootCmd.AddCommand(MigrateCmd(ctx))

	if err := rootCmd.Execute(); err != nil {
		return -1
	}

	return 0

}
