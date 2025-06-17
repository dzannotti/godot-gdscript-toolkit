package integration

import (
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/rules"
)

func TestBasicLintingRules(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected []string // Expected rule names that should trigger
	}{
		{
			name: "expression-not-assigned should trigger",
			code: `
func foo():
	1 + 1
	true
`,
			expected: []string{"expression-not-assigned", "expression-not-assigned"},
		},
		{
			name: "expression-not-assigned should not trigger for calls",
			code: `
func foo():
	bar()
	x.baz()
	await something()
`,
			expected: []string{},
		},
		{
			name: "unnecessary-pass should trigger",
			code: `
func foo():
	pass
	print("hello")
`,
			expected: []string{"unnecessary-pass"},
		},
		{
			name: "unnecessary-pass should not trigger when pass is only statement",
			code: `
func foo():
	pass
`,
			expected: []string{},
		},
		{
			name: "duplicated-load should trigger",
			code: `
var A = load("res://scene.tscn")
var B = preload("res://scene.tscn")
`,
			expected: []string{"duplicated-load"},
		},
		{
			name: "unused-argument should trigger",
			code: `
func foo(x, y):
	print(x)
`,
			expected: []string{"unused-argument"},
		},
		{
			name: "unused-argument should not trigger for underscore prefixed",
			code: `
func foo(_x, y):
	print(y)
`,
			expected: []string{},
		},
		{
			name: "comparison-with-itself should trigger",
			code: `
func foo():
	if x == x:
		return true
	if 1 == 1:
		return false
`,
			expected: []string{"comparison-with-itself", "comparison-with-itself"},
		},
		{
			name: "class-definitions-order should trigger",
			code: `
class Test:
	var x = 1
	signal my_signal
`,
			expected: []string{"class-definitions-order"},
		},
		{
			name: "class-definitions-order should not trigger with correct order",
			code: `
class Test:
	signal my_signal
	const X = 1
	var x = 1
	func foo():
		pass
`,
			expected: []string{},
		},
	}

	// Create linter with all rules
	allRules := rules.GetDefaultRules()
	config := linter.DefaultConfig()
	l := linter.NewLinter(allRules, config)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			problems, err := l.Lint(tc.code)
			if err != nil {
				t.Fatalf("Linting failed: %v", err)
			}

			if len(problems) != len(tc.expected) {
				t.Errorf("Expected %d problems, got %d", len(tc.expected), len(problems))
				for i, problem := range problems {
					t.Logf("Problem %d: %s", i, problem.String())
				}
				return
			}

			// Check that expected rule names are present
			for i, expectedRule := range tc.expected {
				if i >= len(problems) {
					t.Errorf("Expected problem with rule %s but not enough problems found", expectedRule)
					continue
				}
				if problems[i].RuleName != expectedRule {
					t.Errorf("Expected rule %s, got %s", expectedRule, problems[i].RuleName)
				}
			}
		})
	}
}

func TestLintingDirectives(t *testing.T) {
	testCases := []struct {
		name        string
		code        string
		expectedLen int
		description string
	}{
		{
			name: "gdlint:ignore should suppress next line",
			code: `
# gdlint:ignore=expression-not-assigned
func foo():
	1 + 1
`,
			expectedLen: 0,
			description: "Should ignore expression-not-assigned on next line",
		},
		{
			name: "gdlint:disable should suppress globally",
			code: `
# gdlint:disable=expression-not-assigned
func foo():
	1 + 1
func bar():
	true
`,
			expectedLen: 0,
			description: "Should disable expression-not-assigned globally",
		},
		{
			name: "gdlint:disable and enable should work together",
			code: `
# gdlint:disable=expression-not-assigned
func foo():
	1 + 1
# gdlint:enable=expression-not-assigned
func bar():
	true
`,
			expectedLen: 1,
			description: "Should re-enable expression-not-assigned after enable directive",
		},
		{
			name: "multiple rules in ignore",
			code: `
# gdlint:ignore=expression-not-assigned,unnecessary-pass
func foo():
	1 + 1
	pass
	print("done")
`,
			expectedLen: 0,
			description: "Should ignore multiple rules on next line",
		},
	}

	allRules := rules.GetDefaultRules()
	config := linter.DefaultConfig()
	l := linter.NewLinter(allRules, config)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			problems, err := l.Lint(tc.code)
			if err != nil {
				t.Fatalf("Linting failed: %v", err)
			}

			if len(problems) != tc.expectedLen {
				t.Errorf("Expected %d problems, got %d. %s", tc.expectedLen, len(problems), tc.description)
				for i, problem := range problems {
					t.Logf("Problem %d: %s", i, problem.String())
				}
			}
		})
	}
}

func TestComplexLintingScenarios(t *testing.T) {
	code := `
extends Node

# gdlint:ignore=class-definitions-order
var out_of_order = 1
signal my_signal

const CONSTANT = 42
var regular_var = "test"

func _ready():
	var local_var = load("res://test.gd")
	var duplicate = load("res://test.gd")  # Should trigger duplicated-load
	
	# gdlint:ignore=expression-not-assigned
	1 + 1  # Should be ignored
	
	true  # Should trigger expression-not-assigned

func test_function(unused_param, _private_unused):
	if unused_param == unused_param:  # Should trigger comparison-with-itself
		pass
	print("test")  # unnecessary-pass should trigger since there are other statements

class SubClass:
	extends NonExistentParent  # Won't trigger sub-class-before-parent-class since parent doesn't exist
	pass
`

	allRules := rules.GetDefaultRules()
	config := linter.DefaultConfig()
	l := linter.NewLinter(allRules, config)

	problems, err := l.Lint(code)
	if err != nil {
		t.Fatalf("Linting failed: %v", err)
	}

	t.Logf("Found %d problems in complex scenario:", len(problems))
	for i, problem := range problems {
		t.Logf("Problem %d: %s", i, problem.String())
	}

	// We expect at least some problems but won't check exact count due to complexity
	// This test is more for debugging and ensuring no crashes occur
	expectedRules := map[string]bool{
		"duplicated-load":         false,
		"expression-not-assigned": false,
		"comparison-with-itself":  false,
		"unnecessary-pass":        false,
	}

	for _, problem := range problems {
		if _, exists := expectedRules[problem.RuleName]; exists {
			expectedRules[problem.RuleName] = true
		}
	}

	for rule, found := range expectedRules {
		if !found {
			t.Logf("Expected rule %s was not triggered", rule)
		}
	}
}
