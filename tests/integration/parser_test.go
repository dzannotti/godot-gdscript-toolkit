package integration

import (
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/testutil"
)

func TestParserOnAllValidFiles(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	// Test valid scripts
	t.Run("ValidScripts", func(t *testing.T) {
		testutil.TestParserOnValidFiles(t, fixtures.ValidScripts)
	})

	// Test formatter input-output pairs (these should also be valid)
	t.Run("FormatterPairs", func(t *testing.T) {
		testutil.TestParserOnValidFiles(t, fixtures.FormatterPairs)
	})
}

func TestParserOnAllInvalidFiles(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	// Test invalid scripts
	t.Run("InvalidScripts", func(t *testing.T) {
		testutil.TestParserOnInvalidFiles(t, fixtures.InvalidScripts)
	})
}

// Test specific GDScript features that are critical for 1:1 compatibility
func TestParserSpecificFeatures(t *testing.T) {
	testCases := []struct {
		name string
		code string
	}{
		{
			name: "SimpleVariableDeclaration",
			code: "var x = 42",
		},
		{
			name: "FunctionDefinition",
			code: "func foo():\n\tpass",
		},
		{
			name: "ClassDefinition",
			code: "class MyClass:\n\tpass",
		},
		{
			name: "IfStatement",
			code: "if true:\n\tpass",
		},
		{
			name: "ForLoop",
			code: "for i in range(10):\n\tpass",
		},
		{
			name: "WhileLoop",
			code: "while true:\n\tpass",
		},
		{
			name: "MatchStatement",
			code: "match x:\n\t1:\n\t\tpass",
		},
		{
			name: "TypedVariable",
			code: "var x: int = 42",
		},
		{
			name: "FunctionWithParameters",
			code: "func foo(a: int, b: String):\n\tpass",
		},
		{
			name: "FunctionWithReturnType",
			code: "func foo() -> int:\n\treturn 42",
		},
		{
			name: "Expressions",
			code: "1 + 2 * 3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := testutil.ParseCode(t, tc.code)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tc.name, err)
			}
		})
	}
}
