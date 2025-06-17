// Package validation provides tests that validate parsing against Python gdtoolkit fixtures
package validation

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/testutil"
)

// TestPythonGDToolkitFixtureCompatibility tests our parser against all Python gdtoolkit fixtures
// This is the most critical test for ensuring 1:1 compatibility
func TestPythonGDToolkitFixtureCompatibility(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	t.Run("ValidScriptFixtures", func(t *testing.T) {
		testFixtureDirectory(t, fixtures.ValidScripts, true)
	})

	t.Run("InvalidScriptFixtures", func(t *testing.T) {
		testFixtureDirectory(t, fixtures.InvalidScripts, false)
	})

	t.Run("FormatterInputOutputPairs", func(t *testing.T) {
		testFixtureDirectory(t, fixtures.FormatterPairs, true)
	})

	t.Run("PotentialGodotBugs", func(t *testing.T) {
		// These may be edge cases that might not parse correctly
		testFixtureDirectory(t, fixtures.PotentialBugs, true)
	})
}

// TestSpecificFixtureCategories tests specific categories of GDScript features
func TestSpecificFixtureCategories(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	// Test specific fixture files that represent important GDScript features
	importantFixtures := []struct {
		name        string
		file        string
		description string
	}{
		{
			name:        "Functions",
			file:        "functions.gd",
			description: "Core function definition and usage patterns",
		},
		{
			name:        "Expressions",
			file:        "expressions.gd",
			description: "Complex expression parsing and operator precedence",
		},
		{
			name:        "Match",
			file:        "match.gd",
			description: "Match statement with various pattern types",
		},
		{
			name:        "StaticTyping",
			file:        "static_typing.gd",
			description: "Type annotations and static typing features",
		},
		{
			name:        "Annotations",
			file:        "annotations.gd",
			description: "Annotation usage (@export, @tool, etc.)",
		},
		{
			name:        "Enums",
			file:        "enums.gd",
			description: "Enum definitions and usage",
		},
		{
			name:        "Constants",
			file:        "constants.gd",
			description: "Constant declarations",
		},
		{
			name:        "Signals",
			file:        "signals.gd",
			description: "Signal definitions",
		},
		{
			name:        "Properties",
			file:        "properties.gd",
			description: "Property definitions with getters/setters",
		},
		{
			name:        "ClassExtends",
			file:        "extends.gd",
			description: "Class inheritance patterns",
		},
	}

	for _, fixture := range importantFixtures {
		t.Run(fixture.name, func(t *testing.T) {
			filePath := filepath.Join(fixtures.ValidScripts, fixture.file)
			testSingleFixture(t, filePath, fixture.description, true)
		})
	}
}

// TestEdgeCaseFixtures tests known edge cases and potential problem areas
func TestEdgeCaseFixtures(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	edgeCases := []struct {
		name        string
		file        string
		description string
		shouldParse bool
	}{
		{
			name:        "MultilineExpressions",
			file:        "multiline_expressions.gd",
			description: "Expressions spanning multiple lines",
			shouldParse: true,
		},
		{
			name:        "MultilineLambdas",
			file:        "multiline_lambdas.gd",
			description: "Lambda expressions across multiple lines",
			shouldParse: true,
		},
		{
			name:        "ComplexDotChains",
			file:        "complex_dot_chains.gd",
			description: "Long property access chains",
			shouldParse: true,
		},
		{
			name:        "TypedArrays",
			file:        "typed_arrays.gd",
			description: "Typed array declarations and usage",
			shouldParse: true,
		},
		{
			name:        "TypedDicts",
			file:        "typed_dicts.gd",
			description: "Typed dictionary declarations",
			shouldParse: true,
		},
		{
			name:        "NodePaths",
			file:        "node_paths.gd",
			description: "NodePath syntax and usage",
			shouldParse: true,
		},
		{
			name:        "Lambdas",
			file:        "lambdas.gd",
			description: "Lambda function expressions",
			shouldParse: true,
		},
		{
			name:        "TrailingCommaParameters",
			file:        "trailing_comma_after_formal_parameter.gd",
			description: "Functions with trailing comma in parameters",
			shouldParse: true,
		},
		{
			name:        "TrailingCommaArguments",
			file:        "trailing_comma_after_actual_parameter.gd",
			description: "Function calls with trailing comma in arguments",
			shouldParse: true,
		},
	}

	for _, edgeCase := range edgeCases {
		t.Run(edgeCase.name, func(t *testing.T) {
			filePath := filepath.Join(fixtures.ValidScripts, edgeCase.file)
			testSingleFixture(t, filePath, edgeCase.description, edgeCase.shouldParse)
		})
	}
}

// TestFormatterCompatibility tests that formatter input/output pairs parse correctly
func TestFormatterCompatibility(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	// Get all formatter test files
	files, err := testutil.GetGDScriptFiles(fixtures.FormatterPairs)
	if err != nil {
		t.Fatalf("Failed to get formatter test files: %v", err)
	}

	inputFiles := []string{}
	outputFiles := []string{}

	for _, file := range files {
		fileName := filepath.Base(file)
		if strings.Contains(fileName, ".in.gd") {
			inputFiles = append(inputFiles, file)
		} else if strings.Contains(fileName, ".out.gd") {
			outputFiles = append(outputFiles, file)
		}
	}

	t.Logf("Testing %d formatter input files", len(inputFiles))
	t.Logf("Testing %d formatter output files", len(outputFiles))

	// Test that both input and output files parse correctly
	for _, file := range inputFiles {
		fileName := filepath.Base(file)
		t.Run("Input_"+fileName, func(t *testing.T) {
			testSingleFixture(t, file, "Formatter input file", true)
		})
	}

	for _, file := range outputFiles {
		fileName := filepath.Base(file)
		t.Run("Output_"+fileName, func(t *testing.T) {
			testSingleFixture(t, file, "Formatter output file", true)
		})
	}
}

// TestRegressionPrevention tests specific cases that have caused issues in the past
func TestRegressionPrevention(t *testing.T) {
	regressionTests := []struct {
		name        string
		code        string
		description string
		shouldParse bool
	}{
		{
			name:        "EmptyFile",
			code:        "",
			description: "Empty GDScript files should parse without error",
			shouldParse: true,
		},
		{
			name: "CommentsOnly",
			code: `# This is a comment
# Another comment`,
			description: "Files with only comments should parse",
			shouldParse: true,
		},
		{
			name:        "SinglePassStatement",
			code:        "pass",
			description: "Single pass statement should parse",
			shouldParse: true,
		},
		{
			name:        "ClassNameStatement",
			code:        "class_name MyScript",
			description: "class_name declaration should parse",
			shouldParse: true,
		},
		{
			name:        "ExtendsStatement",
			code:        `extends Node`,
			description: "extends declaration should parse",
			shouldParse: true,
		},
		{
			name: "ToolScript",
			code: `@tool
extends EditorPlugin`,
			description: "@tool annotation with extends should parse",
			shouldParse: true,
		},
		{
			name: "FunctionWithDocstring",
			code: `func test():
	"""This is a docstring"""
	pass`,
			description: "Functions with docstrings should parse",
			shouldParse: true,
		},
		{
			name: "NestedFunctionCalls",
			code: `func test():
	result = func1(func2(func3(arg)))`,
			description: "Deeply nested function calls should parse",
			shouldParse: true,
		},
		{
			name: "ComplexTypeAnnotations",
			code: `func test(callback: Callable[[int, String], bool]):
	pass`,
			description: "Complex type annotations should parse",
			shouldParse: true,
		},
		{
			name: "MultilineArrayWithTrailingComma",
			code: `var array = [
	1,
	2,
	3,
]`,
			description: "Multiline arrays with trailing commas should parse",
			shouldParse: true,
		},
	}

	for _, test := range regressionTests {
		t.Run(test.name, func(t *testing.T) {
			result := validateGDScriptCode(t, test.code, test.name)

			if test.shouldParse && !result.ParsedOK {
				t.Errorf("Expected code to parse successfully, but it failed: %v\nDescription: %s",
					result.ParseErrors, test.description)
			} else if !test.shouldParse && result.ParsedOK {
				t.Errorf("Expected code to fail parsing, but it succeeded\nDescription: %s", test.description)
			}

			if result.ParsedOK {
				t.Logf("✓ %s: %s", test.name, test.description)
			}
		})
	}
}

// Helper functions

func testFixtureDirectory(t *testing.T, dir string, shouldParseSuccessfully bool) {
	files, err := testutil.GetGDScriptFiles(dir)
	if err != nil {
		t.Fatalf("Failed to get GDScript files from %s: %v", dir, err)
	}

	if len(files) == 0 {
		t.Skipf("No GDScript files found in %s", dir)
		return
	}

	successCount := 0
	failureCount := 0

	for _, file := range files {
		fileName := filepath.Base(file)

		// Skip known problematic files
		if shouldSkipFixture(fileName) {
			t.Logf("⚠ Skipping known problematic fixture: %s", fileName)
			continue
		}

		success := testSingleFixture(t, file, "", shouldParseSuccessfully)
		if success {
			successCount++
		} else {
			failureCount++
		}
	}

	total := successCount + failureCount
	if total > 0 {
		successRate := float64(successCount) / float64(total) * 100
		t.Logf("Fixture parsing results for %s: %d/%d (%.1f%%) successful",
			filepath.Base(dir), successCount, total, successRate)
	}
}

func testSingleFixture(t *testing.T, filePath, description string, shouldParseSuccessfully bool) bool {
	fileName := filepath.Base(filePath)

	content, err := testutil.LoadTestFile(filePath)
	if err != nil {
		t.Errorf("Failed to load fixture %s: %v", fileName, err)
		return false
	}

	result := validateGDScriptCode(t, content, fileName)

	if shouldParseSuccessfully {
		if !result.ParsedOK {
			t.Errorf("✗ Fixture %s failed to parse: %v", fileName, result.ParseErrors)
			if description != "" {
				t.Errorf("  Description: %s", description)
			}
			return false
		} else {
			t.Logf("✓ Fixture %s parsed successfully", fileName)
			if description != "" {
				t.Logf("  Description: %s", description)
			}
			return true
		}
	} else {
		if result.ParsedOK {
			t.Errorf("✗ Fixture %s was expected to fail parsing but succeeded", fileName)
			return false
		} else {
			t.Logf("✓ Fixture %s correctly failed to parse (as expected)", fileName)
			return true
		}
	}
}

func shouldSkipFixture(fileName string) bool {
	// List of fixtures that are known to be problematic or not yet supported
	skipList := []string{
		"bug_326_multistatement_lambda_corner_case",
		// Add other problematic fixtures as discovered
	}

	for _, skip := range skipList {
		if strings.Contains(fileName, skip) {
			return true
		}
	}
	return false
}
