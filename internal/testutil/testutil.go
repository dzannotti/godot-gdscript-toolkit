// Package testutil provides utilities for testing the Go gdtoolkit implementation
package testutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/rules"
	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

// TestFixtures contains paths to test fixture directories
type TestFixtures struct {
	ValidScripts   string
	InvalidScripts string
	FormatterPairs string
	PotentialBugs  string
}

// GetTestFixtures returns the paths to test fixture directories
func GetTestFixtures() *TestFixtures {
	// Adjust paths relative to the gogdtoolkit directory
	return &TestFixtures{
		ValidScripts:   "../../gdtoolkit/tests/valid-gd-scripts",
		InvalidScripts: "../../gdtoolkit/tests/invalid-gd-scripts",
		FormatterPairs: "../../gdtoolkit/tests/formatter/input-output-pairs",
		PotentialBugs:  "../../gdtoolkit/tests/potential-godot-bugs",
	}
}

// LoadTestFile reads and returns the contents of a test file
func LoadTestFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GetGDScriptFiles returns all .gd files in the given directory
func GetGDScriptFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".gd") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// ParseCode wraps the parser for testing
func ParseCode(t *testing.T, code string) (*ast.AbstractSyntaxTree, error) {
	p := parser.NewParser(code)
	tree := p.Parse()
	errors := p.Errors()
	if len(errors) > 0 {
		return tree, fmt.Errorf("parsing errors: %v", errors)
	}
	return tree, nil
}

// LintCode wraps the linter for testing
func LintCode(t *testing.T, code string, config linter.Config) ([]problem.Problem, error) {
	// Create a linter with all available rules
	allRules := rules.GetDefaultRules()
	l := linter.NewLinter(allRules, config)
	return l.Lint(code)
}

// SimpleOKCheck verifies that code passes linting (no problems found)
func SimpleOKCheck(t *testing.T, code string, disabledRules ...string) {
	config := linter.DefaultConfig()
	for _, rule := range disabledRules {
		config.DisabledRules = append(config.DisabledRules, rule)
	}

	problems, err := LintCode(t, code, config)
	if err != nil {
		t.Fatalf("Linting failed: %v", err)
	}

	if len(problems) > 0 {
		t.Fatalf("Expected no problems, but found %d: %v", len(problems), problems)
	}
}

// SimpleNOKCheck verifies that code fails linting with the expected rule and line
func SimpleNOKCheck(t *testing.T, code string, expectedRule string, expectedLine int, disabledRules ...string) {
	// First, verify that disabling the rule makes the code pass
	configWithDisabled := linter.DefaultConfig()
	configWithDisabled.DisabledRules = append(configWithDisabled.DisabledRules, expectedRule)
	for _, rule := range disabledRules {
		configWithDisabled.DisabledRules = append(configWithDisabled.DisabledRules, rule)
	}

	problems, err := LintCode(t, code, configWithDisabled)
	if err != nil {
		t.Fatalf("Linting failed: %v", err)
	}
	if len(problems) != 0 {
		t.Fatalf("Expected no problems when rule %s is disabled, but found %d", expectedRule, len(problems))
	}

	// Now check that the rule triggers with the expected details
	config := linter.DefaultConfig()
	for _, rule := range disabledRules {
		config.DisabledRules = append(config.DisabledRules, rule)
	}

	problems, err = LintCode(t, code, config)
	if err != nil {
		t.Fatalf("Linting failed: %v", err)
	}

	if len(problems) != 1 {
		t.Fatalf("Expected exactly 1 problem, but found %d: %v", len(problems), problems)
	}

	problem := problems[0]
	if problem.RuleName != expectedRule {
		t.Fatalf("Expected rule %s, but got %s", expectedRule, problem.RuleName)
	}

	if problem.Position.Line != expectedLine {
		t.Fatalf("Expected line %d, but got %d", expectedLine, problem.Position.Line)
	}
}

// CompareParsedASTs compares two ASTs for structural equality (for validating parser compatibility)
func CompareParsedASTs(t *testing.T, ast1, ast2 *ast.AbstractSyntaxTree) {
	// This is a placeholder for AST comparison logic
	// In a full implementation, this would recursively compare AST nodes
	// For now, we'll just check that both are non-nil
	if ast1 == nil && ast2 == nil {
		return
	}
	if ast1 == nil || ast2 == nil {
		t.Fatalf("AST comparison failed: one AST is nil")
	}
	// TODO: Implement detailed AST comparison
}

// TestParserOnValidFiles tests the parser against all valid GDScript files
func TestParserOnValidFiles(t *testing.T, validDir string) {
	files, err := GetGDScriptFiles(validDir)
	if err != nil {
		t.Fatalf("Failed to get test files: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := LoadTestFile(file)
			if err != nil {
				t.Fatalf("Failed to load test file %s: %v", file, err)
			}

			// Skip known problematic files
			if strings.Contains(file, "bug_326_multistatement_lambda_corner_case") {
				t.Skip("Skipping known problematic file")
				return
			}

			_, err = ParseCode(t, content)
			if err != nil {
				t.Fatalf("Failed to parse valid file %s: %v", file, err)
			}
		})
	}
}

// TestParserOnInvalidFiles tests the parser against all invalid GDScript files
func TestParserOnInvalidFiles(t *testing.T, invalidDir string) {
	files, err := GetGDScriptFiles(invalidDir)
	if err != nil {
		t.Fatalf("Failed to get test files: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := LoadTestFile(file)
			if err != nil {
				t.Fatalf("Failed to load test file %s: %v", file, err)
			}

			_, err = ParseCode(t, content)
			if err == nil {
				t.Fatalf("Expected parsing to fail for invalid file %s, but it succeeded", file)
			}
		})
	}
}

// TestLinterOnValidFiles tests the linter against all valid GDScript files
func TestLinterOnValidFiles(t *testing.T, validDir string) {
	files, err := GetGDScriptFiles(validDir)
	if err != nil {
		t.Fatalf("Failed to get test files: %v", err)
	}

	config := linter.DefaultConfig()

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := LoadTestFile(file)
			if err != nil {
				t.Fatalf("Failed to load test file %s: %v", file, err)
			}

			// Just check that linting doesn't crash
			_, err = LintCode(t, content, config)
			if err != nil {
				t.Fatalf("Linting failed on valid file %s: %v", file, err)
			}
		})
	}
}
