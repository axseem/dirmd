package bundler

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/axseem/dirmd/internal/config"
	"github.com/axseem/dirmd/internal/ignorer"
	"github.com/axseem/dirmd/internal/processor"
)

// Bundler orchestrates the file bundling process.
type Bundler struct {
	cfg     *config.Config
	ignorer *ignorer.Ignorer
}

// New creates a new Bundler instance.
func New(cfg *config.Config) (*Bundler, error) {
	ign, err := ignorer.New(cfg.RootDir, cfg.IgnoreFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ignorer: %w", err)
	}
	return &Bundler{
		cfg:     cfg,
		ignorer: ign,
	}, nil
}

// Bundle finds, processes, and bundles all relevant files into a single markdown file.
func (b *Bundler) Bundle() error {
	fmt.Println("- Collecting files...")
	filePaths, err := b.collectFiles()
	if err != nil {
		return fmt.Errorf("error collecting files: %w", err)
	}

	if len(filePaths) == 0 {
		fmt.Println("No files to bundle.")
		return nil
	}
	fmt.Printf("- Found %d files to bundle.\n", len(filePaths))

	fmt.Printf("- Processing files with %d workers...\n", b.cfg.Workers)
	results := b.processFiles(filePaths)

	fmt.Println("- Assembling markdown file...")
	markdownContent, err := b.assembleMarkdown(filePaths, results)
	if err != nil {
		return fmt.Errorf("error assembling markdown: %w", err)
	}

	err = os.WriteFile(b.cfg.OutputPath, []byte(markdownContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to output file %s: %w", b.cfg.OutputPath, err)
	}

	fmt.Printf("- Successfully bundled project to %s\n", b.cfg.OutputPath)
	return nil
}

func (b *Bundler) collectFiles() ([]string, error) {
	var files []string
	err := filepath.WalkDir(b.cfg.RootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(b.cfg.RootDir, path)
		if err != nil {
			return err
		}
		relativePath = filepath.ToSlash(relativePath)
		if relativePath == "." {
			return nil
		}

		if b.ignorer.IsIgnored(relativePath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !b.cfg.IncludeHidden && strings.HasPrefix(d.Name(), ".") && d.Name() != ".gitignore" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func (b *Bundler) processFiles(paths []string) map[string]processor.Result {
	jobs := make(chan string, len(paths))
	resultsChan := make(chan processor.Result, len(paths))
	var wg sync.WaitGroup

	for i := 0; i < b.cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				resultsChan <- processor.ProcessFile(path)
			}
		}()
	}

	for _, path := range paths {
		jobs <- path
	}
	close(jobs)

	wg.Wait()
	close(resultsChan)

	resultsMap := make(map[string]processor.Result)
	for res := range resultsChan {
		resultsMap[res.Path] = res
	}
	return resultsMap
}

type treeNode struct {
	children map[string]*treeNode
	isDir    bool
}

func generateFileTree(rootDir string, paths []string) string {
	root := &treeNode{children: make(map[string]*treeNode), isDir: true}
	for _, path := range paths {
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			continue
		}
		components := strings.Split(filepath.ToSlash(relPath), "/")
		currentNode := root
		for i, component := range components {
			if _, ok := currentNode.children[component]; !ok {
				currentNode.children[component] = &treeNode{children: make(map[string]*treeNode)}
			}
			currentNode = currentNode.children[component]
			if i < len(components)-1 {
				currentNode.isDir = true
			}
		}
	}

	var builder strings.Builder
	dirName := filepath.Base(rootDir)
	builder.WriteString(fmt.Sprintf("# Structure of `%s`\n\n", dirName))

	builder.WriteString(fmt.Sprintf("- `%s/`\n", dirName))
	buildTreeStringRecursive(&builder, root, "  ")
	return builder.String()
}

func buildTreeStringRecursive(builder *strings.Builder, node *treeNode, prefix string) {
	keys := make([]string, 0, len(node.children))
	for k := range node.children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		childNode := node.children[key]
		name := key
		if childNode.isDir {
			name += "/"
		}
		builder.WriteString(fmt.Sprintf("%s- `%s`\n", prefix, name))
		if len(childNode.children) > 0 {
			buildTreeStringRecursive(builder, childNode, prefix+"  ")
		}
	}
}

func (b *Bundler) assembleMarkdown(sortedPaths []string, results map[string]processor.Result) (string, error) {
	var markdownParts []string

	tree := generateFileTree(b.cfg.RootDir, sortedPaths)
	markdownParts = append(markdownParts, tree)

	for _, path := range sortedPaths {
		result, ok := results[path]
		if !ok {
			return "", fmt.Errorf("internal error: result not found for path %s", path)
		}

		if result.ReadError != nil {
			fmt.Fprintf(os.Stderr, "- Could not process file %s: %v\n", path, result.ReadError)
			continue
		}
		if result.IsBinary {
			fmt.Printf("- Skipping binary file: %s\n", path)
			continue
		}

		relPath, err := filepath.Rel(b.cfg.RootDir, path)
		if err != nil {
			relPath = path
		}
		relPath = filepath.ToSlash(relPath)

		var fileContentBuilder strings.Builder
		fileContentBuilder.WriteString("`" + relPath + "`\n")
		fileContentBuilder.WriteString(fmt.Sprintf("```%s\n", result.Language))
		fileContentBuilder.Write(bytes.TrimSpace(result.Content))
		fileContentBuilder.WriteString("\n```")

		markdownParts = append(markdownParts, fileContentBuilder.String())
	}

	return strings.Join(markdownParts, "\n\n"), nil
}
