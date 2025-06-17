// Package ast provides the Abstract Syntax Tree structures for GDScript
package ast

import (
	"fmt"
)

// Position represents a position in the source code
type Position struct {
	Line   int
	Column int
	Offset int
}

// String returns a string representation of the position
func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.Line, p.Column)
}

// Node is the base interface for all AST nodes
type Node interface {
	// Position returns the position of the node in the source code
	Position() Position
	// TokenLiteral returns the literal value of the token
	TokenLiteral() string
}

// Statement is the interface for all statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface for all expression nodes
type Expression interface {
	Node
	expressionNode()
}

// AbstractSyntaxTree represents the root of the AST
type AbstractSyntaxTree struct {
	Pos       Position
	RootClass *Class
	Classes   []*Class
	Functions []*Function
}

// Position returns the position of the AST in the source code
func (ast *AbstractSyntaxTree) Position() Position {
	return ast.Pos
}

// TokenLiteral returns the literal value of the token
func (ast *AbstractSyntaxTree) TokenLiteral() string {
	return "program"
}

// NewAST creates a new AST
func NewAST() *AbstractSyntaxTree {
	return &AbstractSyntaxTree{
		Pos:       Position{Line: 1, Column: 1},
		Classes:   make([]*Class, 0),
		Functions: make([]*Function, 0),
	}
}

// AddClass adds a class to the AST
func (ast *AbstractSyntaxTree) AddClass(class *Class) {
	if ast.RootClass == nil {
		ast.RootClass = class
	}
	ast.Classes = append(ast.Classes, class)
}

// AddFunction adds a function to the AST
func (ast *AbstractSyntaxTree) AddFunction(function *Function) {
	ast.Functions = append(ast.Functions, function)
}
