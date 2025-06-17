package ignorer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIgnorer(t *testing.T) {
	rootDir := t.TempDir()

	gitignoreContent := `
# Comments should be ignored
*.log
/node_modules/
build/
secret.txt
`
	err := os.WriteFile(filepath.Join(rootDir, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write .gitignore: %v", err)
	}

	customIgnorePath := filepath.Join(rootDir, "custom.ignore")
	customIgnoreContent := `
*.tmp
/dist/
`
	err = os.WriteFile(customIgnorePath, []byte(customIgnoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write custom ignore file: %v", err)
	}

	t.Run("with .gitignore only", func(t *testing.T) {
		ign, err := New(rootDir, "")
		if err != nil {
			t.Fatalf("New() with .gitignore failed: %v", err)
		}

		testPaths := map[string]bool{
			"app.log":             true,
			"node_modules/react":  true,
			"build/output.bin":    true,
			"src/build/component": true,
			"secret.txt":          true,
			"src/main.go":         false,
			"README.md":           false,
		}

		for path, shouldBeIgnored := range testPaths {
			if ign.IsIgnored(path) != shouldBeIgnored {
				t.Errorf("path %q: expected ignored=%v, got %v", path, shouldBeIgnored, !shouldBeIgnored)
			}
		}
	})

	t.Run("with custom ignore file only", func(t *testing.T) {
		emptyDir := t.TempDir()
		ign, err := New(emptyDir, customIgnorePath)
		if err != nil {
			t.Fatalf("New() with custom ignore file failed: %v", err)
		}

		testPaths := map[string]bool{
			"data.tmp":           true,
			"dist/bundle.js":     true,
			"src/main.go":        false,
			"node_modules/react": false,
		}

		for path, shouldBeIgnored := range testPaths {
			if ign.IsIgnored(path) != shouldBeIgnored {
				t.Errorf("path %q: expected ignored=%v, got %v", path, shouldBeIgnored, !shouldBeIgnored)
			}
		}
	})

	t.Run("with both .gitignore and custom ignore file", func(t *testing.T) {
		ign, err := New(rootDir, customIgnorePath)
		if err != nil {
			t.Fatalf("New() with both files failed: %v", err)
		}

		testPaths := map[string]bool{
			// From .gitignore
			"app.log":            true,
			"node_modules/react": true,
			// From custom ignore
			"data.tmp":       true,
			"dist/bundle.js": true,
			// Not ignored
			"src/main.go": false,
		}

		for path, shouldBeIgnored := range testPaths {
			if ign.IsIgnored(path) != shouldBeIgnored {
				t.Errorf("path %q: expected ignored=%v, got %v", path, shouldBeIgnored, !shouldBeIgnored)
			}
		}
	})

	t.Run("no ignore files present", func(t *testing.T) {
		emptyDir := t.TempDir()
		ign, err := New(emptyDir, "non-existent-file.ignore")
		if err != nil {
			t.Fatalf("New() with no files failed: %v", err)
		}
		if ign.IsIgnored("any/path/at/all.txt") {
			t.Error("Expected nothing to be ignored when no ignore files exist, but a path was ignored")
		}
	})
}
