package linter

import (
	"fmt"
	"sync"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
	"github.com/dzannotti/gdtoolkit/internal/core/parser"
)

// Rule represents a linting rule
type Rule interface {
	// Check applies the rule to an AST and returns any problems found
	Check(tree *ast.AbstractSyntaxTree, config Config) []problem.Problem
	// Name returns the name of the rule
	Name() string
	// Description returns a description of the rule
	Description() string
}

// Linter performs linting on GDScript code
type Linter struct {
	rules  []Rule
	config Config
}

// NewLinter creates a new linter with the given rules and configuration
func NewLinter(rules []Rule, config Config) *Linter {
	return &Linter{
		rules:  rules,
		config: config,
	}
}

// Lint lints the given code and returns any problems found
func (l *Linter) Lint(code string) ([]problem.Problem, error) {
	// Parse the code
	tree, errors := parser.ParseFile("", code)
	if len(errors) > 0 {
		return nil, fmt.Errorf("parsing errors: %v", errors)
	}

	return l.LintASTWithSource(tree, code), nil
}

// LintAST lints the given AST and returns any problems found
func (l *Linter) LintAST(tree *ast.AbstractSyntaxTree) []problem.Problem {
	return l.LintASTWithSource(tree, "")
}

// LintASTWithSource lints the given AST with source code for directive processing
func (l *Linter) LintASTWithSource(tree *ast.AbstractSyntaxTree, source string) []problem.Problem {
	var problems []problem.Problem

	// Parse directives from source if available
	var ruleContext *RuleContext
	if source != "" {
		directives := ParseDirectives(source)
		ruleContext = NewRuleContext()
		ruleContext.ProcessDirectives(directives)
	}

	// Apply each enabled rule
	for _, rule := range l.rules {
		if l.config.IsRuleEnabled(rule.Name()) {
			ruleProblems := rule.Check(tree, l.config)

			// Filter problems based on directives if we have source code
			if ruleContext != nil {
				var filteredProblems []problem.Problem
				for _, prob := range ruleProblems {
					if ruleContext.IsRuleEnabled(prob.RuleName, prob.Position) {
						filteredProblems = append(filteredProblems, prob)
					}
				}
				problems = append(problems, filteredProblems...)
			} else {
				problems = append(problems, ruleProblems...)
			}
		}
	}

	return problems
}

// LintFile lints the given file and returns any problems found
func (l *Linter) LintFile(filePath string) ([]problem.Problem, error) {
	// Parse the file
	tree, errors := parser.ParseFile(filePath, "")
	if len(errors) > 0 {
		return nil, fmt.Errorf("parsing errors: %v", errors)
	}

	return l.LintAST(tree), nil
}

// LintFiles lints the given files and returns any problems found
func (l *Linter) LintFiles(filePaths []string) (map[string][]problem.Problem, error) {
	results := make(map[string][]problem.Problem)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, filePath := range filePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			problems, err := l.LintFile(path)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results[path] = []problem.Problem{
					problem.NewError(
						ast.Position{Line: 1, Column: 1},
						fmt.Sprintf("Failed to lint file: %v", err),
						"internal",
					),
				}
				return
			}

			results[path] = problems
		}(filePath)
	}

	wg.Wait()
	return results, nil
}

// RegisterRule registers a rule with the linter
func (l *Linter) RegisterRule(rule Rule) {
	l.rules = append(l.rules, rule)
}
