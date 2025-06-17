// Package validation provides AST structure validation tests
package validation

import (
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

// TestASTStructureValidation tests that the AST structure is correctly built for various GDScript constructs
func TestASTStructureValidation(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedAST    func(*ast.AbstractSyntaxTree) error
		shouldFail     bool
		failureMessage string
	}{
		{
			name:  "EmptyScript",
			input: "",
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil {
					return NewValidationError("AST should not be nil")
				}
				if tree.RootClass == nil {
					return NewValidationError("Root class should exist")
				}
				return nil
			},
		},
		{
			name: "SimpleFunctionDeclaration",
			input: `func test():
	pass`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil || tree.RootClass == nil {
					return NewValidationError("AST or root class is nil")
				}
				if len(tree.RootClass.Functions) != 1 {
					return NewValidationError("Expected 1 function, got %d", len(tree.RootClass.Functions))
				}
				fn := tree.RootClass.Functions[0]
				if fn.Name != "test" {
					return NewValidationError("Expected function name 'test', got '%s'", fn.Name)
				}
				if len(fn.Parameters) != 0 {
					return NewValidationError("Expected 0 parameters, got %d", len(fn.Parameters))
				}
				return nil
			},
		},
		{
			name: "FunctionWithParameters",
			input: `func test(a: int, b: String, c = 42):
	pass`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil || tree.RootClass == nil {
					return NewValidationError("AST or root class is nil")
				}
				if len(tree.RootClass.Functions) != 1 {
					return NewValidationError("Expected 1 function, got %d", len(tree.RootClass.Functions))
				}
				fn := tree.RootClass.Functions[0]
				if len(fn.Parameters) != 3 {
					return NewValidationError("Expected 3 parameters, got %d", len(fn.Parameters))
				}

				// Check parameter types
				if fn.Parameters[0].TypeHint != "int" {
					return NewValidationError("Expected first parameter type 'int', got '%s'", fn.Parameters[0].TypeHint)
				}
				if fn.Parameters[1].TypeHint != "String" {
					return NewValidationError("Expected second parameter type 'String', got '%s'", fn.Parameters[1].TypeHint)
				}
				if fn.Parameters[2].Default == nil {
					return NewValidationError("Expected third parameter to have default value")
				}
				return nil
			},
		},
		{
			name: "FunctionWithReturnType",
			input: `func test() -> int:
	return 42`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil || tree.RootClass == nil {
					return NewValidationError("AST or root class is nil")
				}
				if len(tree.RootClass.Functions) != 1 {
					return NewValidationError("Expected 1 function, got %d", len(tree.RootClass.Functions))
				}
				fn := tree.RootClass.Functions[0]
				if fn.ReturnType != "int" {
					return NewValidationError("Expected return type 'int', got '%s'", fn.ReturnType)
				}
				return nil
			},
		},
		{
			name: "ClassDefinition",
			input: `class MyClass:
	pass`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil {
					return NewValidationError("AST should not be nil")
				}
				if len(tree.Classes) < 2 { // root + MyClass
					return NewValidationError("Expected at least 2 classes, got %d", len(tree.Classes))
				}

				// Find MyClass (not the root class)
				var myClass *ast.Class
				for _, class := range tree.Classes {
					if class.Name == "MyClass" {
						myClass = class
						break
					}
				}
				if myClass == nil {
					return NewValidationError("MyClass not found in AST")
				}
				return nil
			},
		},
		{
			name: "ClassWithInheritance",
			input: `class MyClass extends Node:
	pass`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil {
					return NewValidationError("AST should not be nil")
				}

				var myClass *ast.Class
				for _, class := range tree.Classes {
					if class.Name == "MyClass" {
						myClass = class
						break
					}
				}
				if myClass == nil {
					return NewValidationError("MyClass not found in AST")
				}
				if myClass.Extends != "Node" {
					return NewValidationError("Expected extends 'Node', got '%s'", myClass.Extends)
				}
				return nil
			},
		},
		{
			name: "VariableDeclarations",
			input: `var x = 42
var y: int = 10
var z := "hello"`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil || tree.RootClass == nil {
					return NewValidationError("AST or root class is nil")
				}

				// Count variable statements
				varCount := 0
				for _, stmt := range tree.RootClass.Statements {
					if _, ok := stmt.(*ast.VarStatement); ok {
						varCount++
					}
				}

				if varCount != 3 {
					return NewValidationError("Expected 3 variable statements, got %d", varCount)
				}
				return nil
			},
		},
		{
			name: "StaticFunction",
			input: `static func utility():
	pass`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil || tree.RootClass == nil {
					return NewValidationError("AST or root class is nil")
				}
				if len(tree.RootClass.Functions) != 1 {
					return NewValidationError("Expected 1 function, got %d", len(tree.RootClass.Functions))
				}
				fn := tree.RootClass.Functions[0]
				if !fn.IsStatic {
					return NewValidationError("Expected function to be static")
				}
				return nil
			},
		},
		{
			name: "NestedClasses",
			input: `class Outer:
	class Inner:
		pass
	pass`,
			expectedAST: func(tree *ast.AbstractSyntaxTree) error {
				if tree == nil {
					return NewValidationError("AST should not be nil")
				}

				var outerClass *ast.Class
				for _, class := range tree.Classes {
					if class.Name == "Outer" {
						outerClass = class
						break
					}
				}
				if outerClass == nil {
					return NewValidationError("Outer class not found")
				}

				if len(outerClass.SubClasses) != 1 {
					return NewValidationError("Expected 1 subclass in Outer, got %d", len(outerClass.SubClasses))
				}

				if outerClass.SubClasses[0].Name != "Inner" {
					return NewValidationError("Expected subclass name 'Inner', got '%s'", outerClass.SubClasses[0].Name)
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the input
			p := parser.NewParser(tc.input)
			tree := p.Parse()
			errors := p.Errors()

			// Expect parsing to succeed unless shouldFail is true
			if tc.shouldFail {
				if len(errors) == 0 {
					t.Errorf("Expected parsing to fail, but it succeeded")
				}
				return
			}

			if len(errors) > 0 {
				t.Fatalf("Parsing failed with errors: %v", errors)
			}

			// Validate the AST structure
			if tc.expectedAST != nil {
				if err := tc.expectedAST(tree); err != nil {
					t.Errorf("AST validation failed: %v", err)
				}
			}
		})
	}
}

// TestExpressionParsing tests that various expressions are parsed correctly
func TestExpressionParsing(t *testing.T) {
	expressionTests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "ArithmeticExpressions",
			input:       "1 + 2 * 3",
			description: "Basic arithmetic with operator precedence",
		},
		{
			name:        "ComparisonExpressions",
			input:       "x > 5 and y < 10",
			description: "Comparison with logical operators",
		},
		{
			name:        "FunctionCall",
			input:       "func_name(arg1, arg2)",
			description: "Function call with arguments",
		},
		{
			name:        "AttributeAccess",
			input:       "player.health.current",
			description: "Chained attribute access",
		},
		{
			name:        "ArraySubscript",
			input:       "array[index][0]",
			description: "Chained array subscript",
		},
		{
			name:        "TernaryExpression",
			input:       "value if condition else default",
			description: "Ternary conditional expression",
		},
		{
			name:        "TypeCasting",
			input:       "value as int",
			description: "Type casting expression",
		},
		{
			name:        "ArrayLiteral",
			input:       "[1, 2, 3, [4, 5]]",
			description: "Nested array literal",
		},
		{
			name:        "DictionaryLiteral",
			input:       `{"key": "value", "nested": {"inner": true}}`,
			description: "Nested dictionary literal",
		},
	}

	for _, tc := range expressionTests {
		t.Run(tc.name, func(t *testing.T) {
			// Wrap expression in a function to make it parseable
			input := "func test():\n\t" + tc.input

			p := parser.NewParser(input)
			tree := p.Parse()
			errors := p.Errors()

			if len(errors) > 0 {
				t.Errorf("Failed to parse %s (%s): %v", tc.name, tc.description, errors)
				return
			}

			if tree == nil || tree.RootClass == nil || len(tree.RootClass.Functions) == 0 {
				t.Errorf("Expression parsing failed: no function found in AST for %s", tc.name)
				return
			}

			t.Logf("✓ Successfully parsed %s: %s", tc.name, tc.description)
		})
	}
}

// TestControlFlowParsing tests parsing of control flow statements
func TestControlFlowParsing(t *testing.T) {
	controlFlowTests := []struct {
		name  string
		input string
	}{
		{
			name: "IfStatement",
			input: `if condition:
	pass`,
		},
		{
			name: "IfElseStatement",
			input: `if condition:
	pass
else:
	pass`,
		},
		{
			name: "IfElifElseStatement",
			input: `if condition1:
	pass
elif condition2:
	pass
else:
	pass`,
		},
		{
			name: "ForLoop",
			input: `for item in collection:
	pass`,
		},
		{
			name: "WhileLoop",
			input: `while condition:
	pass`,
		},
		{
			name: "MatchStatement",
			input: `match value:
	1:
		pass
	"string":
		pass
	_:
		pass`,
		},
		{
			name: "NestedControlFlow",
			input: `if outer_condition:
	for item in items:
		if inner_condition:
			match item:
				"special":
					break
				_:
					continue`,
		},
	}

	for _, tc := range controlFlowTests {
		t.Run(tc.name, func(t *testing.T) {
			// Wrap in a function
			input := "func test():\n\t" + tc.input

			p := parser.NewParser(input)
			tree := p.Parse()
			errors := p.Errors()

			if len(errors) > 0 {
				t.Errorf("Failed to parse control flow %s: %v", tc.name, errors)
				return
			}

			if tree == nil {
				t.Errorf("Control flow parsing failed: AST is nil for %s", tc.name)
				return
			}

			t.Logf("✓ Successfully parsed control flow: %s", tc.name)
		})
	}
}

// ValidationError represents an AST validation error
type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

// NewValidationError creates a new validation error
func NewValidationError(format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		message: sprintf(format, args...),
	}
}

// sprintf is a simple sprintf implementation for our validation errors
func sprintf(format string, args ...interface{}) string {
	// This is a simplified implementation
	// In a real implementation, you'd use fmt.Sprintf
	result := format
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			result = replaceFirst(result, "%s", v)
		case int:
			result = replaceFirst(result, "%d", intToString(v))
		}
	}
	return result
}

func replaceFirst(s, old, new string) string {
	// Simple replacement - in real code you'd use strings.Replace
	for i := 0; i < len(s)-len(old)+1; i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}

func intToString(i int) string {
	if i == 0 {
		return "0"
	}

	negative := i < 0
	if negative {
		i = -i
	}

	digits := []byte{}
	for i > 0 {
		digits = append([]byte{byte(i%10 + '0')}, digits...)
		i /= 10
	}

	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}
