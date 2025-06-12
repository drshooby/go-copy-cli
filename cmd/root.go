package copy

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Options struct {
	Recursive      bool
	Preserve       bool
	FollowSymlinks bool
}

var (
	rootCmd = &cobra.Command{
		Use:   "gocp",
		Short: "better cp",
		Long:  `better cp written in Go`,
		Args: func(cmd *cobra.Command, args []string) error {
			validator := cobra.MatchAll(
				cobra.OnlyValidArgs,
				cobra.MinimumNArgs(2),
			)
			return validator(cmd, args)
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
		},
	}

	opts = &Options{}
)

func Copy(paths []string, dest string) error {
	for _, src := range paths {
		info, err := os.Lstat(src)
		if err != nil {
			log.Printf("Warning: cannot access %s: %v", src, err)
			continue
		}
		if info.IsDir() {
			if !opts.Recursive {
				return fmt.Errorf("%s is a directory (use -r to copy directories)", src)
			}
			// copy directory
		} else {
			// copy file or symlink
		}
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&opts.Recursive, "recursive", "r", false, "copy directories recursively")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
