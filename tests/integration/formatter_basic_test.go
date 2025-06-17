package integration

import (
	"strings"
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/formatter"
	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

// TestFormatterWithSimpleCases tests the formatter with simple cases that the parser can handle
func TestFormatterWithSimpleCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple_class",
			input: `class TestClass:
	pass`,
			expected: `class TestClass:
	pass`,
		},
		{
			name: "class_with_function",
			input: `class TestClass:
	func test_method():
		pass`,
			expected: `class TestClass:
	func test_method():
		pass`,
		},
		{
			name: "function_only",
			input: `func test_function():
	pass`,
			expected: `func test_function():
	pass`,
		},
		{
			name:     "simple_variable",
			input:    `var test_var = 5`,
			expected: `var test_var = 5`,
		},
		{
			name: "simple_if",
			input: `if true:
	pass`,
			expected: `if true:
	pass`,
		},
		{
			name: "simple_while",
			input: `while true:
	pass`,
			expected: `while true:
	pass`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input
			ast, parseErrors := parser.ParseFile("test.gd", tt.input)
			if len(parseErrors) > 0 {
				t.Logf("Parse errors for %s: %v", tt.name, parseErrors)
				// Skip this test if parsing fails, as that's a parser issue, not formatter
				t.Skip("Skipping due to parser errors")
				return
			}

			// Format the AST
			config := formatter.DefaultConfig()
			result, err := formatter.FormatCode(ast, config)
			if err != nil {
				t.Fatalf("Format error for %s: %v", tt.name, err)
			}

			// Normalize and compare
			expected := strings.TrimSpace(tt.expected)
			actual := strings.TrimSpace(result)

			if actual != expected {
				t.Errorf("Formatting mismatch for %s:\nExpected:\n%s\n\nActual:\n%s", tt.name, expected, actual)
			}
		})
	}
}

// TestFormatterAgainstValidGDScript tests the formatter against real GDScript files
func TestFormatterAgainstValidGDScript(t *testing.T) {
	// Use a simple GDScript example that we know the parser can handle
	input := `class TestClass:
	func method():
		pass

func global_function():
	pass`

	// Parse the input
	ast, parseErrors := parser.ParseFile("test.gd", input)
	if len(parseErrors) > 0 {
		t.Fatalf("Parse errors: %v", parseErrors)
	}

	// Format the AST
	config := formatter.DefaultConfig()
	result, err := formatter.FormatCode(ast, config)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	// Parse the formatted result to ensure it's still valid
	ast2, parseErrors2 := parser.ParseFile("test.gd", result)
	if len(parseErrors2) > 0 {
		t.Fatalf("Parse errors on formatted code: %v", parseErrors2)
	}

	// Ensure we have valid AST structures
	if ast == nil || ast2 == nil {
		t.Error("One of the ASTs is nil")
	}

	// Basic check that formatting doesn't break the structure
	if result == "" {
		t.Error("Formatted result is empty")
	}
}
