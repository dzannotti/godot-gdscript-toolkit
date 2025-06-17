package parser

import (
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
)

func TestParser_Parse_SimpleFunction(t *testing.T) {
	input := `
func test():
	var x = 5
	return x
`

	p := NewParser(input)
	tree := p.Parse()

	// Check for parser errors
	errors := p.Errors()
	if len(errors) > 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %s", err)
		}
		t.FailNow()
	}

	// Verify the AST structure
	if tree == nil {
		t.Fatalf("Parse() returned nil")
	}

	// The root class should exist
	if tree.RootClass == nil {
		t.Fatalf("No root class in AST")
	}

	// The root class should be named "global scope"
	if tree.RootClass.Name != "global scope" {
		t.Errorf("root class name wrong. expected=%q, got=%q",
			"global scope", tree.RootClass.Name)
	}

	// For now, we're just checking that parsing completes without errors
	// In a more complete implementation, we would verify the structure of the AST
}

func TestParser_Parse_ClassDefinition(t *testing.T) {
	input := `
class MyClass:
	var x = 5
	
	func test():
		return x
`

	p := NewParser(input)
	tree := p.Parse()

	// Check for parser errors
	errors := p.Errors()
	if len(errors) > 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %s", err)
		}
		t.FailNow()
	}

	// Verify the AST structure
	if tree == nil {
		t.Fatalf("Parse() returned nil")
	}

	// The root class should exist
	if tree.RootClass == nil {
		t.Fatalf("No root class in AST")
	}

	// For now, we're just checking that parsing completes without errors
	// In a more complete implementation, we would verify the structure of the AST
}

func TestParser_ParseFile(t *testing.T) {
	content := `
func test():
	var x = 5
	return x
`

	tree, errors := ParseFile("test.gd", content)

	// Check for parser errors
	if len(errors) > 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, err := range errors {
			t.Errorf("parser error: %s", err)
		}
		t.FailNow()
	}

	// Verify the AST structure
	if tree == nil {
		t.Fatalf("ParseFile() returned nil")
	}

	// The root class should exist
	if tree.RootClass == nil {
		t.Fatalf("No root class in AST")
	}

	// The root class should be named after the file
	if tree.RootClass.Name != "test.gd" {
		t.Errorf("root class name wrong. expected=%q, got=%q",
			"test.gd", tree.RootClass.Name)
	}
}

func TestVisitor(t *testing.T) {
	// Create a simple AST
	tree := ast.NewAST()
	class := ast.NewClass("TestClass", ast.Position{Line: 1, Column: 1})
	tree.RootClass = class
	tree.Classes = append(tree.Classes, class)

	// Create a visitor that counts nodes
	count := 0
	visitor := &testVisitor{
		countFunc: func(node ast.Node) {
			count++
		},
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	// We should have visited at least 2 nodes (the AST and the class)
	if count < 2 {
		t.Errorf("visitor didn't visit enough nodes, count=%d", count)
	}
}

// testVisitor is a simple visitor that counts nodes
type testVisitor struct {
	countFunc func(ast.Node)
}

func (v *testVisitor) Visit(node ast.Node) ast.Visitor {
	v.countFunc(node)
	return v
}
