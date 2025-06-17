package rules

import (
	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
)

// NoElifReturn checks for unnecessary elif after return
type NoElifReturn struct{}

func (r *NoElifReturn) Name() string {
	return "no-elif-return"
}

func (r *NoElifReturn) Description() string {
	return "Checks for unnecessary elif after return"
}

func (r *NoElifReturn) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	visitor := &noElifReturnVisitor{
		problems: &problems,
	}

	ast.Walk(visitor, tree)
	return problems
}

type noElifReturnVisitor struct {
	problems *[]problem.Problem
}

func (v *noElifReturnVisitor) Visit(node ast.Node) ast.Visitor {
	if ifStmt, ok := node.(*ast.IfStatement); ok {
		// Check if the if branch always returns
		if v.branchAlwaysReturns(ifStmt.Consequence) {
			// Check each elif branch
			for i := range ifStmt.ElseBranches {
				// Create a position for the elif branch
				pos := ast.Position{Line: ifStmt.Position().Line + i + 1, Column: 1}
				*v.problems = append(*v.problems, problem.NewWarning(
					pos,
					"Unnecessary \"elif\" after \"return\"",
					"no-elif-return",
				))
			}
		}
	}
	return v
}

// NoElseReturn checks for unnecessary else after return
type NoElseReturn struct{}

func (r *NoElseReturn) Name() string {
	return "no-else-return"
}

func (r *NoElseReturn) Description() string {
	return "Checks for unnecessary else after return"
}

func (r *NoElseReturn) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	visitor := &noElseReturnVisitor{
		problems: &problems,
	}

	ast.Walk(visitor, tree)
	return problems
}

type noElseReturnVisitor struct {
	problems *[]problem.Problem
}

func (v *noElseReturnVisitor) Visit(node ast.Node) ast.Visitor {
	if ifStmt, ok := node.(*ast.IfStatement); ok {
		// Check if all non-else branches always return
		allNonElseBranchesReturn := true

		// Check if branch
		if !v.branchAlwaysReturns(ifStmt.Consequence) {
			allNonElseBranchesReturn = false
		}

		// Check elif branches
		for _, elifBranch := range ifStmt.ElseBranches {
			if !v.branchAlwaysReturns(elifBranch) {
				allNonElseBranchesReturn = false
				break
			}
		}

		// If all non-else branches return and there's an else branch
		if allNonElseBranchesReturn && len(ifStmt.Alternative) > 0 {
			// Create a position for the else branch
			pos := ast.Position{Line: ifStmt.Position().Line + len(ifStmt.ElseBranches) + 1, Column: 1}
			*v.problems = append(*v.problems, problem.NewWarning(
				pos,
				"Unnecessary \"else\" after \"return\"",
				"no-else-return",
			))
		}
	}
	return v
}

// branchAlwaysReturns checks if a branch always returns
func (v *noElifReturnVisitor) branchAlwaysReturns(statements []ast.Statement) bool {
	return v.hasReturnStatement(statements) || v.hasIfThatAlwaysReturns(statements)
}

// branchAlwaysReturns checks if a branch always returns (for else return visitor)
func (v *noElseReturnVisitor) branchAlwaysReturns(statements []ast.Statement) bool {
	return v.hasReturnStatement(statements) || v.hasIfThatAlwaysReturns(statements)
}

// hasReturnStatement checks if there's a return statement among the statements
func (v *noElifReturnVisitor) hasReturnStatement(statements []ast.Statement) bool {
	for _, stmt := range statements {
		if _, ok := stmt.(*ast.ReturnStatement); ok {
			return true
		}
	}
	return false
}

// hasReturnStatement checks if there's a return statement among the statements (for else return visitor)
func (v *noElseReturnVisitor) hasReturnStatement(statements []ast.Statement) bool {
	for _, stmt := range statements {
		if _, ok := stmt.(*ast.ReturnStatement); ok {
			return true
		}
	}
	return false
}

// hasIfThatAlwaysReturns checks if there's an if statement that always returns
func (v *noElifReturnVisitor) hasIfThatAlwaysReturns(statements []ast.Statement) bool {
	for _, stmt := range statements {
		if ifStmt, ok := stmt.(*ast.IfStatement); ok {
			if v.ifAlwaysReturns(ifStmt) {
				return true
			}
		}
	}
	return false
}

// hasIfThatAlwaysReturns checks if there's an if statement that always returns (for else return visitor)
func (v *noElseReturnVisitor) hasIfThatAlwaysReturns(statements []ast.Statement) bool {
	for _, stmt := range statements {
		if ifStmt, ok := stmt.(*ast.IfStatement); ok {
			if v.ifAlwaysReturns(ifStmt) {
				return true
			}
		}
	}
	return false
}

// ifAlwaysReturns checks if an if statement always returns
func (v *noElifReturnVisitor) ifAlwaysReturns(ifStmt *ast.IfStatement) bool {
	// Must have an else branch to always return
	if len(ifStmt.Alternative) == 0 {
		return false
	}

	// All branches must return
	if !v.branchAlwaysReturns(ifStmt.Consequence) {
		return false
	}

	for _, elifBranch := range ifStmt.ElseBranches {
		if !v.branchAlwaysReturns(elifBranch) {
			return false
		}
	}

	if !v.branchAlwaysReturns(ifStmt.Alternative) {
		return false
	}

	return true
}

// ifAlwaysReturns checks if an if statement always returns (for else return visitor)
func (v *noElseReturnVisitor) ifAlwaysReturns(ifStmt *ast.IfStatement) bool {
	// Must have an else branch to always return
	if len(ifStmt.Alternative) == 0 {
		return false
	}

	// All branches must return
	if !v.branchAlwaysReturns(ifStmt.Consequence) {
		return false
	}

	for _, elifBranch := range ifStmt.ElseBranches {
		if !v.branchAlwaysReturns(elifBranch) {
			return false
		}
	}

	if !v.branchAlwaysReturns(ifStmt.Alternative) {
		return false
	}

	return true
}

// GetDefaultIfReturnRules returns the default if-return checking rules
func GetDefaultIfReturnRules() []linter.Rule {
	return []linter.Rule{
		&NoElifReturn{},
		&NoElseReturn{},
	}
}
