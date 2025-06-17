package rules

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
)

// MaxLineLength checks for lines that are too long
type MaxLineLength struct{}

func (r *MaxLineLength) Name() string {
	return "max-line-length"
}

func (r *MaxLineLength) Description() string {
	return "Checks for lines that exceed the maximum length"
}

func (r *MaxLineLength) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// This rule needs the original source code, not just the AST
	// For now, we'll implement a basic check
	threshold := config.GetRuleSetting(r.Name(), "threshold", 100).(int)
	tabChars := config.GetRuleSetting("tab-characters", "value", 4).(int)

	// TODO: Get source code from tree context
	// For now, return empty problems as we need source access
	_ = threshold
	_ = tabChars

	return problems
}

// MaxFileLines checks for files that have too many lines
type MaxFileLines struct{}

func (r *MaxFileLines) Name() string {
	return "max-file-lines"
}

func (r *MaxFileLines) Description() string {
	return "Checks for files that exceed the maximum number of lines"
}

func (r *MaxFileLines) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	threshold := config.GetRuleSetting(r.Name(), "threshold", 1000).(int)

	// TODO: Get line count from tree context
	// For now, return empty problems as we need source access
	_ = threshold

	return problems
}

// TrailingWhitespace checks for trailing whitespace
type TrailingWhitespace struct{}

func (r *TrailingWhitespace) Name() string {
	return "trailing-whitespace"
}

func (r *TrailingWhitespace) Description() string {
	return "Checks for trailing whitespace in lines"
}

func (r *TrailingWhitespace) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// TODO: Get source code from tree context
	// For now, return empty problems as we need source access

	return problems
}

// MixedTabsAndSpaces checks for mixed tabs and spaces in indentation
type MixedTabsAndSpaces struct{}

func (r *MixedTabsAndSpaces) Name() string {
	return "mixed-tabs-and-spaces"
}

func (r *MixedTabsAndSpaces) Description() string {
	return "Checks for mixed tabs and spaces in indentation"
}

func (r *MixedTabsAndSpaces) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// TODO: Get source code from tree context
	// For now, return empty problems as we need source access

	return problems
}

// FormatChecker provides source-based format checking
type FormatChecker struct {
	source string
}

// NewFormatChecker creates a new format checker with source code
func NewFormatChecker(source string) *FormatChecker {
	return &FormatChecker{source: source}
}

// CheckMaxLineLength checks for lines that exceed maximum length
func (f *FormatChecker) CheckMaxLineLength(threshold, tabChars int) []problem.Problem {
	var problems []problem.Problem
	lines := strings.Split(f.source, "\n")

	for lineNum, line := range lines {
		// Replace tabs with spaces for length calculation
		expandedLine := strings.ReplaceAll(line, "\t", strings.Repeat(" ", tabChars))

		if len(expandedLine) > threshold {
			problems = append(problems, problem.NewWarning(
				ast.Position{Line: lineNum + 1, Column: 1},
				"Max allowed line length ("+strconv.Itoa(threshold)+") exceeded",
				"max-line-length",
			))
		}
	}

	return problems
}

// CheckMaxFileLines checks for files that exceed maximum number of lines
func (f *FormatChecker) CheckMaxFileLines(threshold int) []problem.Problem {
	var problems []problem.Problem
	lines := strings.Split(f.source, "\n")

	if len(lines) > threshold {
		problems = append(problems, problem.NewWarning(
			ast.Position{Line: len(lines), Column: 1},
			"Max allowed file lines num ("+strconv.Itoa(threshold)+") exceeded",
			"max-file-lines",
		))
	}

	return problems
}

// CheckTrailingWhitespace checks for trailing whitespace
func (f *FormatChecker) CheckTrailingWhitespace() []problem.Problem {
	var problems []problem.Problem
	lines := strings.Split(f.source, "\n")
	trailingWSRegex := regexp.MustCompile(`\s$`)

	for lineNum, line := range lines {
		if trailingWSRegex.MatchString(line) {
			problems = append(problems, problem.NewWarning(
				ast.Position{Line: lineNum + 1, Column: 1},
				"Trailing whitespace(s)",
				"trailing-whitespace",
			))
		}
	}

	return problems
}

// CheckMixedTabsAndSpaces checks for mixed tabs and spaces in indentation
func (f *FormatChecker) CheckMixedTabsAndSpaces() []problem.Problem {
	var problems []problem.Problem
	lines := strings.Split(f.source, "\n")
	mixedRegex := regexp.MustCompile(`^(\t+ +| +\t+)`)

	for lineNum, line := range lines {
		if mixedRegex.MatchString(line) {
			problems = append(problems, problem.NewWarning(
				ast.Position{Line: lineNum + 1, Column: 1},
				"Mixed tabs and spaces",
				"mixed-tabs-and-spaces",
			))
		}
	}

	return problems
}

// GetDefaultFormatRules returns the default format checking rules
func GetDefaultFormatRules() []linter.Rule {
	return []linter.Rule{
		&MaxLineLength{},
		&MaxFileLines{},
		&TrailingWhitespace{},
		&MixedTabsAndSpaces{},
	}
}
