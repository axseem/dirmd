package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/axseem/dirmd/internal/bundler"
	"github.com/axseem/dirmd/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cfg := config.NewDefaultConfig()

	cmd := &cobra.Command{
		Use:   "dirmd <directory>",
		Short: "Bundles all files from a directory into a single markdown file.",
		Long: `dirmd is a CLI tool that traverses a specified directory,
reads all non-ignored files, and bundles them into a single,
well-formatted markdown file.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootDir, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid directory path: %w", err)
			}
			info, err := os.Stat(rootDir)
			if err != nil {
				return fmt.Errorf("cannot access directory: %w", err)
			}
			if !info.IsDir() {
				return fmt.Errorf("path is not a directory: %s", rootDir)
			}
			cfg.RootDir = rootDir

			b, err := bundler.New(cfg)
			if err != nil {
				return err
			}

			return b.Bundle()
		},
	}

	cmd.Flags().StringVarP(&cfg.OutputPath, "output", "o", cfg.OutputPath, "Path for the output markdown file")
	cmd.Flags().StringVarP(&cfg.IgnoreFilePath, "ignore-file", "i", "", "Path to a custom .gitignore-style file to use for ignoring files")
	cmd.Flags().IntVarP(&cfg.Workers, "workers", "w", cfg.Workers, "Number of concurrent workers for processing files")
	cmd.Flags().BoolVar(&cfg.IncludeHidden, "include-hidden", cfg.IncludeHidden, "Include hidden files and directories (those starting with a dot)")

	return cmd
}
