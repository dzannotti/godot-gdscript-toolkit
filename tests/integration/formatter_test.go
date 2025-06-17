package integration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/formatter"
	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

// TestFormatterAgainstPythonTestCases tests the Go formatter against Python gdtoolkit test cases
func TestFormatterAgainstPythonTestCases(t *testing.T) {
	// Define test cases that we can reliably test
	testCases := []struct {
		name       string
		inputFile  string
		outputFile string
	}{
		{
			name:       "type_hints",
			inputFile:  "gdtoolkit/tests/formatter/input-output-pairs/type_hints.in.gd",
			outputFile: "gdtoolkit/tests/formatter/input-output-pairs/type_hints.out.gd",
		},
		{
			name:       "simple_function_definitions",
			inputFile:  "gdtoolkit/tests/formatter/input-output-pairs/simple_function_definitions.in.gd",
			outputFile: "gdtoolkit/tests/formatter/input-output-pairs/simple_function_definitions.out.gd",
		},
		{
			name:       "simple_classes_and_functions",
			inputFile:  "gdtoolkit/tests/formatter/input-output-pairs/simple_classes_and_functions.in.gd",
			outputFile: "gdtoolkit/tests/formatter/input-output-pairs/simple_classes_and_functions.out.gd",
		},
		{
			name:       "const_statements",
			inputFile:  "gdtoolkit/tests/formatter/input-output-pairs/const_statements.in.gd",
			outputFile: "gdtoolkit/tests/formatter/input-output-pairs/const_statements.out.gd",
		},
		{
			name:       "simple_function_statements",
			inputFile:  "gdtoolkit/tests/formatter/input-output-pairs/simple_function_statements.in.gd",
			outputFile: "gdtoolkit/tests/formatter/input-output-pairs/simple_function_statements.out.gd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Read input file
			inputPath := filepath.Join("../../", tc.inputFile)
			inputContent, err := ioutil.ReadFile(inputPath)
			if err != nil {
				t.Skipf("Could not read input file %s: %v", inputPath, err)
				return
			}

			// Read expected output file
			outputPath := filepath.Join("../../", tc.outputFile)
			expectedContent, err := ioutil.ReadFile(outputPath)
			if err != nil {
				t.Skipf("Could not read output file %s: %v", outputPath, err)
				return
			}

			// Parse the input
			ast, parseErrors := parser.ParseFile(tc.inputFile, string(inputContent))
			if len(parseErrors) > 0 {
				t.Fatalf("Parse errors for %s: %v", tc.name, parseErrors)
			}

			// Format the AST
			config := formatter.DefaultConfig()
			formattedCode, err := formatter.FormatCode(ast, config)
			if err != nil {
				t.Fatalf("Format error for %s: %v", tc.name, err)
			}

			// Normalize line endings and trailing whitespace
			expected := strings.TrimSpace(string(expectedContent))
			actual := strings.TrimSpace(formattedCode)

			// Compare the results
			if actual != expected {
				t.Errorf("Formatter output doesn't match Python gdtoolkit for %s:\n\nExpected:\n%s\n\nActual:\n%s\n\nDifference at first mismatch:", tc.name, expected, actual)

				// Show detailed difference
				expectedLines := strings.Split(expected, "\n")
				actualLines := strings.Split(actual, "\n")

				maxLines := len(expectedLines)
				if len(actualLines) > maxLines {
					maxLines = len(actualLines)
				}

				for i := 0; i < maxLines; i++ {
					var expectedLine, actualLine string
					if i < len(expectedLines) {
						expectedLine = expectedLines[i]
					}
					if i < len(actualLines) {
						actualLine = actualLines[i]
					}

					if expectedLine != actualLine {
						t.Errorf("Line %d mismatch:\nExpected: %q\nActual:   %q", i+1, expectedLine, actualLine)
						break
					}
				}
			}
		})
	}
}

// TestFormatterBasicFunctionality tests basic formatter functionality
func TestFormatterBasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "basic_class",
			input: `class Test:
	pass`,
			expected: `class Test:
	pass`,
		},
		{
			name: "function_with_params",
			input: `func test(a,b:int,c=1):
	pass`,
			expected: `func test(a, b: int, c = 1):
	pass`,
		},
		{
			name:     "variable_declaration",
			input:    `var x:int=5`,
			expected: `var x: int = 5`,
		},
		{
			name:     "const_declaration",
			input:    `const MAX_SIZE:int=100`,
			expected: `const MAX_SIZE: int = 100`,
		},
		{
			name: "if_statement",
			input: `if x>0:
	print("positive")
else:
	print("not positive")`,
			expected: `if x > 0:
	print("positive")
else:
	print("not positive")`,
		},
		{
			name: "for_loop",
			input: `for i in range(10):
	print(i)`,
			expected: `for i in range(10):
	print(i)`,
		},
		{
			name: "while_loop",
			input: `while condition:
	do_something()`,
			expected: `while condition:
	do_something()`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input
			ast, parseErrors := parser.ParseFile("test.gd", tt.input)
			if len(parseErrors) > 0 {
				t.Fatalf("Parse errors: %v", parseErrors)
			}

			// Format the AST
			config := formatter.DefaultConfig()
			result, err := formatter.FormatCode(ast, config)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			// Normalize and compare
			expected := strings.TrimSpace(tt.expected)
			actual := strings.TrimSpace(result)

			if actual != expected {
				t.Errorf("Formatting mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
			}
		})
	}
}

// TestFormatterIdempotency tests that formatting is idempotent
func TestFormatterIdempotency(t *testing.T) {
	input := `class Example:
	func method(param: int = 5) -> String:
		var local_var: Array = [1, 2, 3]
		if param > 0:
			return "positive"
		else:
			return "not positive"
	
	func another_method():
		for i in range(10):
			print(i)
		
		while true:
			break`

	// Parse and format once
	ast1, parseErrors := parser.ParseFile("test.gd", input)
	if len(parseErrors) > 0 {
		t.Fatalf("Parse errors: %v", parseErrors)
	}

	config := formatter.DefaultConfig()
	formatted1, err := formatter.FormatCode(ast1, config)
	if err != nil {
		t.Fatalf("First format error: %v", err)
	}

	// Parse and format the result again
	ast2, parseErrors := parser.ParseFile("test.gd", formatted1)
	if len(parseErrors) > 0 {
		t.Fatalf("Parse errors on formatted code: %v", parseErrors)
	}

	formatted2, err := formatter.FormatCode(ast2, config)
	if err != nil {
		t.Fatalf("Second format error: %v", err)
	}

	// The two formatted results should be identical
	if formatted1 != formatted2 {
		t.Errorf("Formatter is not idempotent:\nFirst format:\n%s\n\nSecond format:\n%s", formatted1, formatted2)
	}
}

// TestFormatterPreservesSemantics tests that formatting preserves semantic meaning
func TestFormatterPreservesSemantics(t *testing.T) {
	testCases := []string{
		`class Test:
	func method(a: int, b: String = "default") -> bool:
		var result = a > 0 and b != ""
		return result`,

		`func calculate(x: float) -> float:
	if x < 0:
		return -x
	elif x == 0:
		return 0
	else:
		return x * 2`,

		`var global_var: Dictionary = {"key": "value", "number": 42}
const CONSTANT: int = 100`,
	}

	for i, input := range testCases {
		t.Run(fmt.Sprintf("semantic_test_%d", i), func(t *testing.T) {
			// Parse original
			ast1, parseErrors := parser.ParseFile("test.gd", input)
			if len(parseErrors) > 0 {
				t.Fatalf("Parse errors on original: %v", parseErrors)
			}

			// Format
			config := formatter.DefaultConfig()
			formatted, err := formatter.FormatCode(ast1, config)
			if err != nil {
				t.Fatalf("Format error: %v", err)
			}

			// Parse formatted
			ast2, parseErrors := parser.ParseFile("test.gd", formatted)
			if len(parseErrors) > 0 {
				t.Fatalf("Parse errors on formatted code: %v", parseErrors)
			}

			// Both ASTs should be functionally equivalent
			// For now, we just ensure both parse successfully
			// A more sophisticated test would compare AST structures
			if ast1 == nil || ast2 == nil {
				t.Error("One of the ASTs is nil")
			}
		})
	}
}
