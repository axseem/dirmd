package bundler

import (
	"testing"
)

func TestGenerateFileTree(t *testing.T) {
	rootDir := "/home/user/project"
	testCases := []struct {
		name           string
		paths          []string
		expectedOutput string
	}{
		{
			name: "simple structure",
			paths: []string{
				"/home/user/project/main.go",
				"/home/user/project/README.md",
				"/home/user/project/internal/server.go",
			},
			expectedOutput: `# Structure of ` + "`project`" + `

- ` + "`project/`" + `
  - ` + "`README.md`" + `
  - ` + "`internal/`" + `
    - ` + "`server.go`" + `
  - ` + "`main.go`" + `
`,
		},
		{
			name:  "empty paths",
			paths: []string{},
			expectedOutput: `# Structure of ` + "`project`" + `

- ` + "`project/`" + `
`,
		},
		{
			name: "deeply nested structure with unsorted input",
			paths: []string{
				"/home/user/project/go.mod",
				"/home/user/project/api/v1/handler.go",
				"/home/user/project/main.go",
				"/home/user/project/api/v2/handler.go",
				"/home/user/project/api/v1/model.go",
			},
			expectedOutput: `# Structure of ` + "`project`" + `

- ` + "`project/`" + `
  - ` + "`api/`" + `
    - ` + "`v1/`" + `
      - ` + "`handler.go`" + `
      - ` + "`model.go`" + `
    - ` + "`v2/`" + `
      - ` + "`handler.go`" + `
  - ` + "`go.mod`" + `
  - ` + "`main.go`" + `
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := generateFileTree(rootDir, tc.paths)
			if got != tc.expectedOutput {
				t.Errorf("generateFileTree() mismatch:\n--- EXPECTED ---\n%s\n\n--- GOT ---\n%s", tc.expectedOutput, got)
			}
		})
	}
}
