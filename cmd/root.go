package copy

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gocp",
		Short: "better cp",
		Long:  `better cp written in Go`,
		Args: func(cmd *cobra.Command, args []string) error {
			validator := cobra.MatchAll(
				cobra.OnlyValidArgs,
				cobra.MinimumNArgs(1),
			)
			return validator(cmd, args)
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hi")
		},
	}

	Recursive bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "copy directories recursively")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
