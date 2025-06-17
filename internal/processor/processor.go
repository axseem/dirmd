package processor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// langExtMap maps file extensions and specific filenames to Markdown language identifiers.
var langExtMap = map[string]string{
	"dockerfile":  "dockerfile",
	"gemfile":     "ruby",
	"go.mod":      "go",
	"go.sum":      "text",
	"makefile":    "makefile",
	"vagrantfile": "ruby",

	".ada":        "ada",
	".as":         "actionscript",
	".bat":        "batch",
	".bash":       "bash",
	".c":          "c",
	".h":          "c",
	".clj":        "clojure",
	".coffee":     "coffeescript",
	".conf":       "ini",
	".cpp":        "cpp",
	".hpp":        "cpp",
	".cxx":        "cpp",
	".cs":         "csharp",
	".css":        "css",
	".csv":        "csv",
	".d":          "d",
	".dart":       "dart",
	".diff":       "diff",
	".elm":        "elm",
	".erl":        "erlang",
	".ex":         "elixir",
	".exs":        "elixir",
	".f90":        "fortran",
	".fs":         "fsharp",
	".go":         "go",
	".groovy":     "groovy",
	".hcl":        "hcl",
	".hs":         "haskell",
	".html":       "html",
	".ini":        "ini",
	".java":       "java",
	".jl":         "julia",
	".js":         "javascript",
	".mjs":        "javascript",
	".json":       "json",
	".jsx":        "jsx",
	".kt":         "kotlin",
	".kts":        "kotlin",
	".less":       "less",
	".lisp":       "lisp",
	".lua":        "lua",
	".m":          "objectivec",
	".md":         "markdown",
	".mk":         "makefile",
	".ml":         "ocaml",
	".patch":      "diff",
	".perl":       "perl",
	".pl":         "perl",
	".php":        "php",
	".phtml":      "php",
	".properties": "properties",
	".proto":      "protobuf",
	".ps1":        "powershell",
	".py":         "python",
	".r":          "r",
	".rb":         "ruby",
	".rs":         "rust",
	".sass":       "sass",
	".scala":      "scala",
	".scm":        "scheme",
	".scss":       "scss",
	".sh":         "shell",
	".sql":        "sql",
	".svelte":     "svelte",
	".swift":      "swift",
	".tcl":        "tcl",
	".tex":        "latex",
	".tf":         "terraform",
	".toml":       "toml",
	".ts":         "typescript",
	".tsx":        "tsx",
	".vb":         "vbnet",
	".vbs":        "vbscript",
	".vue":        "vue",
	".xml":        "xml",
	".xsd":        "xml",
	".xsl":        "xml",
	".yaml":       "yaml",
	".yml":        "yaml",
	".zig":        "zig",
	".zsh":        "zsh",
}

// Result represents the processed content of a single file.
type Result struct {
	Path      string
	Content   []byte
	Language  string
	IsBinary  bool
	ReadError error
}

// ProcessFile reads a file and returns its raw content and metadata.
func ProcessFile(path string) Result {
	content, err := os.ReadFile(path)
	if err != nil {
		return Result{Path: path, ReadError: fmt.Errorf("reading file: %w", err)}
	}

	if isBinary(content) {
		return Result{Path: path, IsBinary: true}
	}

	lang := getLanguage(path)

	return Result{
		Path:     path,
		Content:  content,
		Language: lang,
	}
}

// getLanguage determines the language for syntax highlighting.
// It first checks the full filename, then the file extension.
func getLanguage(path string) string {
	filename := strings.ToLower(filepath.Base(path))
	if lang, ok := langExtMap[filename]; ok {
		return lang
	}

	ext := filepath.Ext(filename)
	if lang, ok := langExtMap[ext]; ok {
		return lang
	}

	return strings.TrimPrefix(ext, ".")
}

func isBinary(data []byte) bool {
	return bytes.Contains(data, []byte{0})
}
