package ignorer

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Ignorer determines whether a file or directory should be ignored.
type Ignorer struct {
	gitIgnorer *ignore.GitIgnore
}

// New creates a new Ignorer, loading rules from the
// root directory's .gitignore and a custom ignore file.
func New(rootDir, customIgnoreFile string) (*Ignorer, error) {
	var allPatterns []string

	if gitIgnorePath := filepath.Join(rootDir, ".gitignore"); isFileExist(gitIgnorePath) {
		lines, err := readLines(gitIgnorePath)
		if err != nil {
			return nil, err
		}
		allPatterns = append(allPatterns, lines...)
	}

	if customIgnoreFile != "" && isFileExist(customIgnoreFile) {
		lines, err := readLines(customIgnoreFile)
		if err != nil {
			return nil, err
		}
		allPatterns = append(allPatterns, lines...)
	}

	ignorer := ignore.CompileIgnoreLines(allPatterns...)

	return &Ignorer{gitIgnorer: ignorer}, nil
}

func (i *Ignorer) IsIgnored(path string) bool {
	return i.gitIgnorer.MatchesPath(path)
}

func isFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
