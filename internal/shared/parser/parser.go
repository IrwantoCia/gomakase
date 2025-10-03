package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

type ASTParser interface {
	AddDependencies(codes []string) error
	AddImport(importPath string, alias string)
	AddRoute(route string)
	WriteFile()
}

type astParser struct {
	filePath string
	file     *ast.File
	fset     *token.FileSet
}

func NewASTParser(
	filePath string,
) ASTParser {
	src, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, src, 0)
	if err != nil {
		log.Fatalf("Error parsing file: %v", err)
	}

	return &astParser{
		file:     file,
		fset:     fset,
		filePath: filePath,
	}
}

func (r *astParser) AddImport(importPath string, alias string) {
	// Check if this exact import (path and alias) already exists to avoid duplicates.
	for _, i := range r.file.Imports {
		// Check path match
		if i.Path.Value == `"`+importPath+`"` {
			// If we want an alias, check if it has one and if it matches.
			if alias != "" {
				if i.Name != nil && i.Name.Name == alias {
					return // Exact match found, do nothing.
				}
			} else {
				// If we don't want an alias, check that it doesn't have one.
				if i.Name == nil {
					return // Exact match found, do nothing.
				}
			}
		}
	}

	// Create a new import spec
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + importPath + `"`,
		},
	}

	if alias != "" {
		newImport.Name = ast.NewIdent(alias)
	}

	// Find the import declaration block
	var importDecl *ast.GenDecl
	for _, decl := range r.file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.IMPORT {
			importDecl = genDecl
			break
		}
	}

	// Add the new import to the existing block or create a new one
	if importDecl != nil {
		importDecl.Specs = append(importDecl.Specs, newImport)
	} else {
		// No import block found, create a new one
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{newImport},
		}
		// Add the new import block to the top of the file's declarations
		r.file.Decls = append([]ast.Decl{importDecl}, r.file.Decls...)
	}

}

func (r *astParser) AddRoute(route string) {
	ast.Inspect(r.file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "Routes" {
			return true
		}

		insertionIndex := -1
		// ... logic to find insertionIndex ...
		for i, stmt := range funcDecl.Body.List {
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						// here we will insert after the last router
						if selExpr.X.(*ast.Ident).Name == "router" {
							insertionIndex = i
						}
					}
				}
			}
		}

		if insertionIndex == -1 {
			insertionIndex = len(funcDecl.Body.List)
		}

		insertionIndex++

		originalStmts := funcDecl.Body.List
		var newStmts []ast.Stmt
		stmt, err := r.parseStmt(route)
		if err != nil {
			log.Fatalf("Failed to parse statement: %v %v", route, err)
			return false
		}
		newStmts = append(newStmts, stmt)

		// check if the code is already in the list
		var cleanNewStmts []ast.Stmt
		for _, newStmt := range newStmts {
			var isExists = false
			for _, stmt := range originalStmts {
				isExists = r.isStmtExists(stmt, newStmt)
				if isExists {
					break
				}
			}
			if !isExists {
				cleanNewStmts = append(cleanNewStmts, newStmt)
			}
		}

		funcDecl.Body.List = append(
			originalStmts[:insertionIndex],
			append(cleanNewStmts, originalStmts[insertionIndex:]...)...)

		return false // Stop searching
	})
}

func (r *astParser) WriteFile() {
	var buf bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	err := cfg.Fprint(&buf, r.fset, r.file) // Use Fprint from the configured printer
	if err != nil {
		log.Fatalf("Failed to format code: %v", err)
	}
	os.WriteFile(r.filePath, buf.Bytes(), 0644)
}

// AddDependencies adds dependencies to the filepath.
// This function parses the AST to identify the appropriate location for adding
// dependencies initialization.
//
// The function will:
// - Look for an existing dependencies initialization
// - Insert the dependencies initialization before the router
//
// Returns an error if the dependencies cannot be added due to parsing issues
// or if the file structure is incompatible.
func (r *astParser) AddDependencies(codes []string) error {
	var parseErr error
	ast.Inspect(r.file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "Routes" {
			return true
		}

		insertionIndex := -1
		// ... logic to find insertionIndex ...
		for i, stmt := range funcDecl.Body.List {
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if selExpr.X.(*ast.Ident).Name == "router" {
							insertionIndex = i
							break
						}
					}
				}
			}
		}

		if insertionIndex == -1 {
			insertionIndex = len(funcDecl.Body.List)
		}

		originalStmts := funcDecl.Body.List
		var newStmts []ast.Stmt
		for _, code := range codes {
			stmt, err := r.parseStmt(code)
			if err != nil {
				log.Fatalf("Failed to parse statement: %v %v", code, err)
				parseErr = err
				return false
			}
			newStmts = append(newStmts, stmt)
		}

		// check if the code is already in the list
		var cleanNewStmts []ast.Stmt
		for _, newStmt := range newStmts {
			var isExists = false
			for _, stmt := range originalStmts {
				isExists = r.isStmtExists(stmt, newStmt)
				if isExists {
					break
				}
			}
			if !isExists {
				cleanNewStmts = append(cleanNewStmts, newStmt)
			}
		}

		funcDecl.Body.List = append(
			originalStmts[:insertionIndex],
			append(cleanNewStmts, originalStmts[insertionIndex:]...)...)

		return false // Stop searching
	})
	return parseErr
}

// ParseStmt takes a string containing a single Go statement and returns the corresponding AST node.
// It now also strips all position information from the new node, which is CRITICAL
// for preventing the go/printer from mis-formatting it upon insertion.
func (r *astParser) parseStmt(code string) (ast.Stmt, error) {
	src := fmt.Sprintf("package p\n\nfunc f() { %s }", code)
	file, err := parser.ParseFile(token.NewFileSet(), "src.go", src, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse code string: %w", err)
	}

	fn, ok := file.Decls[0].(*ast.FuncDecl)
	if !ok || len(fn.Body.List) == 0 {
		return nil, fmt.Errorf("failed to find statement in parsed code")
	}

	// Get the statement we care about
	stmt := fn.Body.List[0]

	// *** THE CRITICAL FIX IS HERE ***
	// Walk the new AST snippet and remove all position info.
	// This tells the printer to format it based on its structure, not its
	// (now meaningless) original position.
	ast.Inspect(stmt, func(n ast.Node) bool {
		if n != nil {
			// A slightly more robust way to do this, which works
			// for different node types.
			switch node := n.(type) {
			case *ast.AssignStmt:
				node.TokPos = token.NoPos
			case *ast.CallExpr:
				node.Lparen = token.NoPos
				node.Rparen = token.NoPos
			case *ast.SelectorExpr:
				// No Pos field to reset directly here, it's composed of X and Sel
			case *ast.Ident:
				node.NamePos = token.NoPos
			}
		}
		return true
	})

	return stmt, nil
}

func (r *astParser) nodeToString(n ast.Node) (string, error) {
	var buf bytes.Buffer
	// Use a new, empty fileset to erase all position info.
	// This is the key to getting a canonical, comparable string.
	fset := token.NewFileSet()
	err := printer.Fprint(&buf, fset, n)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func (r *astParser) isStmtExists(stmt ast.Stmt, code ast.Stmt) bool {
	stmtString, err := r.nodeToString(stmt)
	if err != nil {
		return false
	}
	codeString, err := r.nodeToString(code)
	if err != nil {
		return false
	}
	return stmtString == codeString
}
