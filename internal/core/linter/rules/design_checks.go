package rules

import (
	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
)

// MaxPublicMethods checks for too many public methods in a class
type MaxPublicMethods struct{}

func (r *MaxPublicMethods) Name() string {
	return "max-public-methods"
}

func (r *MaxPublicMethods) Description() string {
	return "Checks for too many public methods in a class"
}

func (r *MaxPublicMethods) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem
	threshold := config.GetRuleSetting(r.Name(), "threshold", 20).(int)

	visitor := &maxPublicMethodsVisitor{
		problems:  &problems,
		threshold: threshold,
	}

	ast.Walk(visitor, tree)
	return problems
}

type maxPublicMethodsVisitor struct {
	problems  *[]problem.Problem
	threshold int
}

func (v *maxPublicMethodsVisitor) Visit(node ast.Node) ast.Visitor {
	if class, ok := node.(*ast.Class); ok {
		publicMethods := 0
		for _, stmt := range class.Statements {
			if function, ok := stmt.(*ast.Function); ok {
				if isPublicFunction(function.Name) {
					publicMethods++
				}
			}
		}

		if publicMethods > v.threshold {
			className := "Global scope class"
			if class.Name != "" {
				className = "Class " + class.Name
			}

			*v.problems = append(*v.problems, problem.NewWarning(
				class.Position(),
				className+" has more than "+string(rune(v.threshold+48))+" public methods (functions)",
				"max-public-methods",
			))
		}
	}
	return v
}

// MaxReturns checks for too many return statements in a function
type MaxReturns struct{}

func (r *MaxReturns) Name() string {
	return "max-returns"
}

func (r *MaxReturns) Description() string {
	return "Checks for too many return statements in a function"
}

func (r *MaxReturns) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem
	threshold := config.GetRuleSetting(r.Name(), "threshold", 6).(int)

	visitor := &maxReturnsVisitor{
		problems:  &problems,
		threshold: threshold,
	}

	ast.Walk(visitor, tree)
	return problems
}

type maxReturnsVisitor struct {
	problems  *[]problem.Problem
	threshold int
}

func (v *maxReturnsVisitor) Visit(node ast.Node) ast.Visitor {
	if function, ok := node.(*ast.Function); ok {
		returnCount := 0
		var lastReturn *ast.ReturnStatement

		// Count return statements in the function
		for _, stmt := range function.Statements {
			if retStmt, ok := stmt.(*ast.ReturnStatement); ok {
				returnCount++
				lastReturn = retStmt
			}
		}

		if returnCount > v.threshold && lastReturn != nil {
			*v.problems = append(*v.problems, problem.NewWarning(
				lastReturn.Position(),
				"Function \""+function.Name+"\" has more than "+string(rune(v.threshold+48))+" return statements",
				"max-returns",
			))
		}
	}
	return v
}

// FunctionArgumentsNumber checks for too many function arguments
type FunctionArgumentsNumber struct{}

func (r *FunctionArgumentsNumber) Name() string {
	return "function-arguments-number"
}

func (r *FunctionArgumentsNumber) Description() string {
	return "Checks for too many function arguments"
}

func (r *FunctionArgumentsNumber) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem
	threshold := config.GetRuleSetting(r.Name(), "threshold", 10).(int)

	visitor := &functionArgumentsNumberVisitor{
		problems:  &problems,
		threshold: threshold,
	}

	ast.Walk(visitor, tree)
	return problems
}

type functionArgumentsNumberVisitor struct {
	problems  *[]problem.Problem
	threshold int
}

func (v *functionArgumentsNumberVisitor) Visit(node ast.Node) ast.Visitor {
	if function, ok := node.(*ast.Function); ok {
		if len(function.Parameters) > v.threshold {
			*v.problems = append(*v.problems, problem.NewWarning(
				function.Position(),
				"Function \""+function.Name+"\" has more than "+string(rune(v.threshold+48))+" arguments",
				"function-arguments-number",
			))
		}
	}
	return v
}

// isPublicFunction checks if a function name indicates it's public
func isPublicFunction(name string) bool {
	return len(name) > 0 && name[0] != '_'
}

// GetDefaultDesignRules returns the default design checking rules
func GetDefaultDesignRules() []linter.Rule {
	return []linter.Rule{
		&MaxPublicMethods{},
		&MaxReturns{},
		&FunctionArgumentsNumber{},
	}
}
