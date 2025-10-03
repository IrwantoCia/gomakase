package parser

import (
	"testing"
)

func TestRouter(t *testing.T) {
	filePath := "router.go"
	parser := NewASTParser(filePath)
	parser.AddImport("github.com/IrwantoCia/gomakase/internal/auth/application", "authApp")
	parser.WriteFile()
}

func TestAddDependencies(t *testing.T) {
	filePath := "router.go"
	parser := NewASTParser(filePath)
	code := `_ = "bar"`
	err := parser.AddDependencies([]string{code})
	if err != nil {
		t.Fatalf("Failed to parse statement: %v", err)
	}
	parser.WriteFile()
}

func TestAddRouter(t *testing.T) {
	filePath := "router.go"
	parser := NewASTParser(filePath)
	parser.AddRoute("router.GET(\"/login\", authHandler.LoginPage)")
	parser.WriteFile()
}
