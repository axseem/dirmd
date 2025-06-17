package processor

import (
	"testing"
)

func TestGetLanguage(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		expectedLang string
	}{
		{
			name:         "Go file by extension",
			path:         "project/main.go",
			expectedLang: "go",
		},
		{
			name:         "Dockerfile by filename",
			path:         "Dockerfile",
			expectedLang: "dockerfile",
		},
		{
			name:         "Makefile by filename",
			path:         "/usr/bin/Makefile",
			expectedLang: "makefile",
		},
		{
			name:         "go.mod by filename",
			path:         "go.mod",
			expectedLang: "go",
		},
		{
			name:         "Case-insensitive filename",
			path:         "GEMFILE",
			expectedLang: "ruby",
		},
		{
			name:         "Unknown extension",
			path:         "archive.zip",
			expectedLang: "zip",
		},
		{
			name:         "File with no extension",
			path:         "README",
			expectedLang: "",
		},
		{
			name:         "Dotfile with extension",
			path:         ".config.json",
			expectedLang: "json",
		},
		{
			name:         "Dotfile with no extension",
			path:         ".bashrc",
			expectedLang: "bashrc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lang := getLanguage(tc.path)
			if lang != tc.expectedLang {
				t.Errorf("getLanguage(%q) = %q; want %q", tc.path, lang, tc.expectedLang)
			}
		})
	}
}

func TestIsBinary(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "Simple text",
			data:     []byte("hello world"),
			expected: false,
		},
		{
			name:     "UTF-8 text",
			data:     []byte("こんにちは"),
			expected: false,
		},
		{
			name:     "Text with newline and tab",
			data:     []byte("line 1\n\tline 2"),
			expected: false,
		},
		{
			name:     "Data with null byte",
			data:     []byte("this is a \x00 binary"),
			expected: true,
		},
		{
			name:     "Image file header (PNG)",
			data:     []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			expected: false,
		},
		{
			name:     "ELF executable header",
			data:     []byte{0x7f, 'E', 'L', 'F', 0x02, 0x01, 0x01, 0x00},
			expected: true,
		},
		{
			name:     "Empty data",
			data:     []byte{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isBin := isBinary(tc.data)
			if isBin != tc.expected {
				t.Errorf("isBinary() = %v; want %v", isBin, tc.expected)
			}
		})
	}
}
