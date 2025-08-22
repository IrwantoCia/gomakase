package main

import (
	"fmt"
	"testing"

	"github.com/IrwantoCia/gomakase/cmd/gomakase/libs"
	"github.com/stretchr/testify/assert"
)

const TEMPLATE_ROOT = "templates"

func TestListFiles(t *testing.T) {
	service := libs.NewService("test", TEMPLATE_ROOT+"/init", templatesFS)

	_, err := service.ListFiles(TEMPLATE_ROOT + "/init")
	assert.NoError(t, err)
}

func TestGetOutpathTableDriven(t *testing.T) {
	service := libs.NewService("test", TEMPLATE_ROOT+"/init", templatesFS)

	tests := []struct {
		input    string
		expected string
	}{
		{input: TEMPLATE_ROOT + "/init/package.json.tmpl", expected: "test/package.json"},
		{input: TEMPLATE_ROOT + "/init/cmd/server/main.go.tmpl", expected: "test/cmd/server/main.go"},
		{input: TEMPLATE_ROOT + "/init/Makefile", expected: "test/Makefile"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("input: %s, expected: %s", test.input, test.expected)
		t.Run(testname, func(t *testing.T) {
			outpath, err := service.GetOutpath(test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, outpath)
		})
	}
}

func TestGetInitTasks(t *testing.T) {
	service := libs.NewService("foo", TEMPLATE_ROOT+"/init", templatesFS)

	_, err := service.GetInitTasks()
	assert.NoError(t, err)
}

func TestGetContextTasks(t *testing.T) {
	service := libs.NewService(".", TEMPLATE_ROOT+"/context", templatesFS)

	tasks, err := service.GetContextTasks("foo")
	for _, task := range tasks {
		fmt.Println(task.Outpath)
	}
	assert.NoError(t, err)
}

func TestConvertToLowerCamelCaseTableDriven(t *testing.T) {
	service := libs.NewService(".", TEMPLATE_ROOT+"/context", templatesFS)

	tests := []struct {
		input    string
		expected string
	}{
		{input: "foo", expected: "foo"},
		{input: "bar", expected: "bar"},
		{input: "baz", expected: "baz"},
		{input: "foo_bar", expected: "fooBar"},
		{input: "FooBar", expected: "fooBar"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("input: %s, expected: %s", test.input, test.expected)
		t.Run(testname, func(t *testing.T) {
			assert.Equal(t, test.expected, service.ConvertToLowerCamelCase(test.input))
		})
	}
}

func TestConvertToUpperCamelCaseTableDriven(t *testing.T) {
	service := libs.NewService(".", TEMPLATE_ROOT+"/context", templatesFS)

	tests := []struct {
		input    string
		expected string
	}{
		{input: "foo", expected: "Foo"},
		{input: "bar", expected: "Bar"},
		{input: "baz", expected: "Baz"},
		{input: "foo_bar", expected: "FooBar"},
		{input: "FooBar", expected: "FooBar"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("input: %s, expected: %s", test.input, test.expected)
		t.Run(testname, func(t *testing.T) {
			assert.Equal(t, test.expected, service.ConvertToUpperCamelCase(test.input))
		})
	}
}
