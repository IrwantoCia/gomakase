package file

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type File interface {
	CreateFile(path string, content []byte) error
	IsPathExists(path string) bool
	ParseFilePath(path string, data map[string]string) (string, error)
	ParseTemplate(content []byte, data map[string]string) ([]byte, error)
}

type file struct {
}

func NewFile() File {
	return &file{}
}

func (f *file) CreateFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	return os.WriteFile(path, content, 0644)
}

func (f *file) IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (f *file) ParseFilePath(path string, data map[string]string) (string, error) {
	tmpl, err := template.New("output").Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"title": cases.Title(language.English).String,
	}).Parse(path)
	if err != nil {
		return path, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return path, err
	}
	return buf.String(), nil
}

func (f *file) ParseTemplate(content []byte, data map[string]string) ([]byte, error) {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"title": cases.Title(language.English).String,
	}).Parse(string(content))
	if err != nil {
		return []byte{}, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
