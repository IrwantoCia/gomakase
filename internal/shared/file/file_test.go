package file

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestFile_ParseTemplate(t *testing.T) {
	input := "test module"
	expected := "output " + input

	content, err := os.ReadFile("test.tmpl")
	if err != nil {
		t.Fatalf("Error reading template: %v", err)
	}

	file := NewFile()
	content, err = file.ParseTemplate(content, map[string]string{
		"Module": input,
	})

	if err != nil {
		t.Fatalf("Error parsing template: %v", err)
	}

	fmt.Printf("Content: %s\n", string(content))
	assert.Equal(t, string(content), expected)
}

func TestFile_ParseFilePath(t *testing.T) {
	tests := []struct {
		input    string
		data     map[string]string
		expected string
	}{
		{
			input:    "{{ .Module }}/go.mod",
			data:     map[string]string{"Module": "test"},
			expected: "test/go.mod",
		},
		{
			input:    "internal/{{ .ContextName | lower }}/domain/{{ .ContextName | lower }}.entity.go",
			data:     map[string]string{"ContextName": "myapp"},
			expected: "internal/myapp/domain/myapp.entity.go",
		},
		{
			input:    "internal/{{ .ContextName | title }}/delivery/{{ .ContextName | title }}.handler.go",
			data:     map[string]string{"ContextName": "myapp"},
			expected: "internal/Myapp/delivery/Myapp.handler.go",
		},
	}

	file := NewFile()
	for _, tt := range tests {
		filePath, err := file.ParseFilePath(tt.input, tt.data)
		if err != nil {
			t.Fatalf("Error parsing file path: %v", err)
		}
		assert.Equal(t, filePath, tt.expected)
	}
}

func TestFile_IsPathExists(t *testing.T) {
	input := "test.tmpl"
	expected := true

	file := NewFile()
	exists := file.IsPathExists(input)
	assert.Equal(t, exists, expected)
}
