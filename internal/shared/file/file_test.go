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
	input := "{{ .Module }}/go.mod"
	expected := "test/go.mod"

	file := NewFile()
	filePath, err := file.ParseFilePath(input, map[string]string{
		"Module": "test",
	})
	if err != nil {
		t.Fatalf("Error parsing file path: %v", err)
	}

	assert.Equal(t, filePath, expected)
}

func TestFile_IsPathExists(t *testing.T) {
	input := "test.tmpl"
	expected := true

	file := NewFile()
	exists := file.IsPathExists(input)
	assert.Equal(t, exists, expected)
}
