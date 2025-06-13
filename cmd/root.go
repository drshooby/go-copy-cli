package copy

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Options struct {
	Recursive      bool
	Preserve       bool
	FollowSymlinks bool
	Verbose        bool
}

var (
	rootCmd = &cobra.Command{
		Use:   "gocp",
		Short: "better cp",
		Long:  `better cp written in Go`,
		Args: func(cmd *cobra.Command, args []string) error {
			validator := cobra.MatchAll(
				cobra.MinimumNArgs(2),
			)
			return validator(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Copy(args[:len(args)-1], args[len(args)-1])
		},
	}

	opts = &Options{}
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func CopyDirectory(src, dst string) error {

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if opts.Verbose {
			log.Printf("Warning: directory at path \"%s\" doesn't exist, creating path...", dst)
		}
		if err := os.MkdirAll(dst, 0777); err != nil {
			return fmt.Errorf("error: failed to create directory %v", err)
		}
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("there was en error walking %v", err)
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(dstPath, 0777); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dstPath, err)

			}
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0777); err != nil {
			return err
		}

		if opts.Verbose {
			log.Printf("Copying file: %s into %s", path, dstPath)
		}
		return CopyFile(path, dstPath)
	})
}

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
			if err = CopyDirectory(src, dest); err != nil {
				fmt.Println(err)
			}
		} else {
			// copy file or symlink
		}
	}
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&opts.Recursive, "recursive", "r", false, "copy directories recursively")
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "display verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
