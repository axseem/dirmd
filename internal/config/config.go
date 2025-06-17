package config

import "runtime"

// Config holds all the configuration for the dirmd tool.
type Config struct {
	// RootDir is the source directory to bundle.
	RootDir string
	// OutputPath is the path to the output markdown file.
	OutputPath string
	// IgnoreFilePath is the path to a custom .gitignore-style file.
	IgnoreFilePath string
	// Workers is the number of concurrent workers to use for file processing.
	Workers int
	// IncludeHidden specifies whether to include hidden files and directories.
	IncludeHidden bool
}

// NewDefaultConfig creates a new configuration with default values.
func NewDefaultConfig() *Config {
	return &Config{
		OutputPath:    "bundle.md",
		Workers:       runtime.NumCPU(),
		IncludeHidden: false,
	}
}
