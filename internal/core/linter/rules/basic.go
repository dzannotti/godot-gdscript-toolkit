package rules

import (
	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
)

// ExpressionNotAssigned checks for expressions that are evaluated but not assigned or used
type ExpressionNotAssigned struct{}

// Name returns the name of the rule
func (r *ExpressionNotAssigned) Name() string {
	return "expression-not-assigned"
}

// Description returns a description of the rule
func (r *ExpressionNotAssigned) Description() string {
	return "Checks for expressions that are evaluated but not assigned or used"
}

// Check applies the rule to an AST and returns any problems found
func (r *ExpressionNotAssigned) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Create a visitor to find expression statements
	visitor := &expressionNotAssignedVisitor{
		problems: &problems,
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	return problems
}

// expressionNotAssignedVisitor is a visitor that finds expression statements
type expressionNotAssignedVisitor struct {
	problems *[]problem.Problem
}

// Visit is called for each node in the AST
func (v *expressionNotAssignedVisitor) Visit(node ast.Node) ast.Visitor {
	// Check if the node is an expression statement
	if exprStmt, ok := node.(*ast.ExpressionStatement); ok {
		// Check if the expression is a call expression (allowed)
		if _, ok := exprStmt.Expression.(*ast.CallExpression); ok {
			return v
		}

		// Check if the expression is an assignment expression (allowed)
		if _, ok := exprStmt.Expression.(*ast.AssignmentExpression); ok {
			return v
		}

		// Check if the expression is a string literal (docstring, allowed)
		if _, ok := exprStmt.Expression.(*ast.StringLiteral); ok {
			return v
		}

		// Only flag simple expressions that are actually not used
		switch exprStmt.Expression.(type) {
		case *ast.InfixExpression, *ast.NumberLiteral, *ast.BooleanLiteral, *ast.Identifier:
			*v.problems = append(*v.problems, problem.NewWarning(
				exprStmt.Position(),
				"Expression is not assigned to a variable or used",
				"expression-not-assigned",
			))
		}
	}

	return v
}

// UnnecessaryPass checks for unnecessary pass statements
type UnnecessaryPass struct{}

// Name returns the name of the rule
func (r *UnnecessaryPass) Name() string {
	return "unnecessary-pass"
}

// Description returns a description of the rule
func (r *UnnecessaryPass) Description() string {
	return "Checks for unnecessary pass statements"
}

// Check applies the rule to an AST and returns any problems found
func (r *UnnecessaryPass) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Create a visitor to find unnecessary pass statements
	visitor := &unnecessaryPassVisitor{
		problems: &problems,
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	return problems
}

// unnecessaryPassVisitor is a visitor that finds unnecessary pass statements
type unnecessaryPassVisitor struct {
	problems *[]problem.Problem
	inBlock  bool
	hasStmt  bool
}

// Visit is called for each node in the AST
func (v *unnecessaryPassVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Function:
		// Check if the function has statements other than pass
		hasStmt := false
		hasPass := false
		var passStmt *ast.PassStatement

		for _, stmt := range n.Statements {
			if _, ok := stmt.(*ast.PassStatement); ok {
				hasPass = true
				passStmt = stmt.(*ast.PassStatement)
			} else {
				hasStmt = true
			}
		}

		// If the function has both pass and other statements, the pass is unnecessary
		if hasPass && hasStmt && passStmt != nil {
			*v.problems = append(*v.problems, problem.NewWarning(
				passStmt.Position(),
				"Unnecessary pass statement",
				"unnecessary-pass",
			))
		}

	case *ast.Class:
		// Check if the class has statements other than pass
		hasStmt := false
		hasPass := false
		var passStmt *ast.PassStatement

		for _, stmt := range n.Statements {
			if _, ok := stmt.(*ast.PassStatement); ok {
				hasPass = true
				passStmt = stmt.(*ast.PassStatement)
			} else {
				hasStmt = true
			}
		}

		// If the class has both pass and other statements, the pass is unnecessary
		if hasPass && hasStmt && passStmt != nil {
			*v.problems = append(*v.problems, problem.NewWarning(
				passStmt.Position(),
				"Unnecessary pass statement",
				"unnecessary-pass",
			))
		}
	}

	return v
}

// DuplicatedLoad checks for duplicated load/preload statements
type DuplicatedLoad struct{}

// Name returns the name of the rule
func (r *DuplicatedLoad) Name() string {
	return "duplicated-load"
}

// Description returns a description of the rule
func (r *DuplicatedLoad) Description() string {
	return "Checks for duplicated load/preload statements"
}

// Check applies the rule to an AST and returns any problems found
func (r *DuplicatedLoad) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Create a visitor to find duplicated load statements
	visitor := &duplicatedLoadVisitor{
		problems: &problems,
		loads:    make(map[string]ast.Position),
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	return problems
}

// duplicatedLoadVisitor is a visitor that finds duplicated load statements
type duplicatedLoadVisitor struct {
	problems *[]problem.Problem
	loads    map[string]ast.Position
}

// Visit is called for each node in the AST
func (v *duplicatedLoadVisitor) Visit(node ast.Node) ast.Visitor {
	// Check for variable statements with load/preload calls
	if varStmt, ok := node.(*ast.VarStatement); ok {
		if callExpr, ok := varStmt.Value.(*ast.CallExpression); ok {
			if ident, ok := callExpr.Function.(*ast.Identifier); ok {
				if ident.Value == "load" || ident.Value == "preload" {
					// Check if we have arguments
					if len(callExpr.Arguments) > 0 {
						if stringLit, ok := callExpr.Arguments[0].(*ast.StringLiteral); ok {
							path := stringLit.Value
							if prevPos, exists := v.loads[path]; exists {
								*v.problems = append(*v.problems, problem.NewWarning(
									varStmt.Position(),
									"Duplicated load statement for '"+path+"'",
									"duplicated-load",
								))
								_ = prevPos // suppress unused variable warning
							} else {
								v.loads[path] = varStmt.Position()
							}
						}
					}
				}
			}
		}
	}

	return v
}

// UnusedArgument checks for unused function arguments
type UnusedArgument struct{}

// Name returns the name of the rule
func (r *UnusedArgument) Name() string {
	return "unused-argument"
}

// Description returns a description of the rule
func (r *UnusedArgument) Description() string {
	return "Checks for unused function arguments"
}

// Check applies the rule to an AST and returns any problems found
func (r *UnusedArgument) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Create a visitor to find unused arguments
	visitor := &unusedArgumentVisitor{
		problems: &problems,
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	return problems
}

// unusedArgumentVisitor is a visitor that finds unused arguments
type unusedArgumentVisitor struct {
	problems *[]problem.Problem
}

// Visit is called for each node in the AST
func (v *unusedArgumentVisitor) Visit(node ast.Node) ast.Visitor {
	if function, ok := node.(*ast.Function); ok {
		// Track used identifiers in the function body
		usedIdentifiers := make(map[string]bool)

		// Simple check: look for identifiers in the function body
		bodyVisitor := &identifierCollector{used: usedIdentifiers}
		for _, stmt := range function.Statements {
			ast.Walk(bodyVisitor, stmt)
		}

		// Check each parameter
		for _, param := range function.Parameters {
			// Skip parameters that start with underscore (conventional way to mark unused)
			if len(param.Name) > 0 && param.Name[0] == '_' {
				continue
			}

			if !usedIdentifiers[param.Name] {
				*v.problems = append(*v.problems, problem.NewWarning(
					param.Pos,
					"Unused argument '"+param.Name+"'",
					"unused-argument",
				))
			}
		}
	}

	return v
}

// identifierCollector collects used identifiers
type identifierCollector struct {
	used map[string]bool
}

// Visit is called for each node in the AST
func (v *identifierCollector) Visit(node ast.Node) ast.Visitor {
	if ident, ok := node.(*ast.Identifier); ok {
		v.used[ident.Value] = true
	}
	return v
}

// ComparisonWithItself checks for comparisons of identical expressions
type ComparisonWithItself struct{}

// Name returns the name of the rule
func (r *ComparisonWithItself) Name() string {
	return "comparison-with-itself"
}

// Description returns a description of the rule
func (r *ComparisonWithItself) Description() string {
	return "Checks for comparisons of identical expressions"
}

// Check applies the rule to an AST and returns any problems found
func (r *ComparisonWithItself) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Create a visitor to find self-comparisons
	visitor := &comparisonWithItselfVisitor{
		problems: &problems,
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	return problems
}

// comparisonWithItselfVisitor is a visitor that finds self-comparisons
type comparisonWithItselfVisitor struct {
	problems *[]problem.Problem
}

// Visit is called for each node in the AST
func (v *comparisonWithItselfVisitor) Visit(node ast.Node) ast.Visitor {
	if binExpr, ok := node.(*ast.InfixExpression); ok {
		// Check if it's a comparison operator
		if binExpr.Operator == "==" || binExpr.Operator == "!=" ||
			binExpr.Operator == "<" || binExpr.Operator == ">" ||
			binExpr.Operator == "<=" || binExpr.Operator == ">=" {

			// Check if left and right expressions are the same
			if expressionsEqual(binExpr.Left, binExpr.Right) {
				*v.problems = append(*v.problems, problem.NewWarning(
					binExpr.Position(),
					"Comparison of identical expressions",
					"comparison-with-itself",
				))
			}
		}
	}

	return v
}

// expressionsEqual checks if two expressions are structurally equal
func expressionsEqual(left, right ast.Expression) bool {
	// Simple implementation - compare string representation
	// In a more complete implementation, this would do structural comparison
	switch l := left.(type) {
	case *ast.Identifier:
		if r, ok := right.(*ast.Identifier); ok {
			return l.Value == r.Value
		}
	case *ast.NumberLiteral:
		if r, ok := right.(*ast.NumberLiteral); ok {
			return l.Original == r.Original
		}
	case *ast.StringLiteral:
		if r, ok := right.(*ast.StringLiteral); ok {
			return l.Value == r.Value
		}
	case *ast.InfixExpression:
		if r, ok := right.(*ast.InfixExpression); ok {
			return l.Operator == r.Operator &&
				expressionsEqual(l.Left, r.Left) &&
				expressionsEqual(l.Right, r.Right)
		}
	}
	return false
}
