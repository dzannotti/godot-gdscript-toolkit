package rules

import (
	"regexp"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
)

// NameCheckRule represents a generic name checking rule
type NameCheckRule struct {
	name        string
	description string
	pattern     string
	compiled    *regexp.Regexp
}

// NewNameCheckRule creates a new name checking rule
func NewNameCheckRule(name, description, pattern string) *NameCheckRule {
	compiled, _ := regexp.Compile("^" + pattern + "$")
	return &NameCheckRule{
		name:        name,
		description: description,
		pattern:     pattern,
		compiled:    compiled,
	}
}

// Name returns the name of the rule
func (r *NameCheckRule) Name() string {
	return r.name
}

// Description returns a description of the rule
func (r *NameCheckRule) Description() string {
	return r.description
}

// Check applies the rule to an AST and returns any problems found
func (r *NameCheckRule) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Get the pattern from config if available, otherwise use default
	pattern := config.GetRuleSetting(r.name, "pattern", r.pattern).(string)
	compiled, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		compiled = r.compiled // fallback to default
	}

	visitor := &nameCheckVisitor{
		problems: &problems,
		ruleName: r.name,
		pattern:  compiled,
	}

	ast.Walk(visitor, tree)
	return problems
}

type nameCheckVisitor struct {
	problems *[]problem.Problem
	ruleName string
	pattern  *regexp.Regexp
}

func (v *nameCheckVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Function:
		if v.ruleName == "function-name" {
			v.checkName(n.Name, n.Position(), "Function")
		}
		// Check function arguments
		if v.ruleName == "function-argument-name" {
			for _, param := range n.Parameters {
				v.checkName(param.Name, param.Pos, "Function argument")
			}
		}
	case *ast.Class:
		if v.ruleName == "sub-class-name" && n.Name != "" {
			v.checkName(n.Name, n.Position(), "Class")
		}
	case *ast.VarStatement:
		if v.ruleName == "function-variable-name" {
			// Check if this is in a function scope
			if v.isInFunctionScope(n) && !v.hasLoadCall(n) {
				v.checkName(n.Name, n.Position(), "Function variable")
			}
		}
		if v.ruleName == "function-preload-variable-name" {
			if v.isInFunctionScope(n) && v.hasPreloadCall(n) {
				v.checkName(n.Name, n.Position(), "Function preload variable")
			}
		}
		if v.ruleName == "class-variable-name" {
			if !v.isInFunctionScope(n) && !v.hasLoadCall(n) {
				v.checkName(n.Name, n.Position(), "Class variable")
			}
		}
		if v.ruleName == "class-load-variable-name" {
			if !v.isInFunctionScope(n) && v.hasLoadCall(n) {
				v.checkName(n.Name, n.Position(), "Class load variable")
			}
		}
		// Check constants (VarStatement with IsConst = true)
		if v.ruleName == "constant-name" && n.IsConst && !v.hasLoadCall(n) {
			v.checkName(n.Name, n.Position(), "Constant")
		}
		if v.ruleName == "load-constant-name" && n.IsConst && v.hasLoadCall(n) {
			v.checkName(n.Name, n.Position(), "Load constant")
		}
	}
	return v
}

func (v *nameCheckVisitor) checkName(name string, pos ast.Position, context string) {
	if name != "" && !v.pattern.MatchString(name) {
		*v.problems = append(*v.problems, problem.NewWarning(
			pos,
			context+" name \""+name+"\" is not valid",
			v.ruleName,
		))
	}
}

func (v *nameCheckVisitor) isInFunctionScope(node ast.Node) bool {
	// TODO: This is a simplified check. In a complete implementation,
	// we would need to track the scope properly during AST traversal
	return false // For now, assume global scope
}

func (v *nameCheckVisitor) hasLoadCall(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.VarStatement:
		return v.isLoadOrPreloadCall(n.Value)
	}
	return false
}

func (v *nameCheckVisitor) hasPreloadCall(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.VarStatement:
		return v.isPreloadCall(n.Value)
	}
	return false
}

func (v *nameCheckVisitor) isLoadOrPreloadCall(expr ast.Expression) bool {
	if callExpr, ok := expr.(*ast.CallExpression); ok {
		if ident, ok := callExpr.Function.(*ast.Identifier); ok {
			return ident.Value == "load" || ident.Value == "preload"
		}
	}
	return false
}

func (v *nameCheckVisitor) isPreloadCall(expr ast.Expression) bool {
	if callExpr, ok := expr.(*ast.CallExpression); ok {
		if ident, ok := callExpr.Function.(*ast.Identifier); ok {
			return ident.Value == "preload"
		}
	}
	return false
}

// GetDefaultNameRules returns the default name checking rules
func GetDefaultNameRules() []linter.Rule {
	return []linter.Rule{
		NewNameCheckRule("function-name", "Function name should follow naming conventions", `^(_|[a-z])[a-z0-9]*(_[a-z0-9]+)*$`),
		NewNameCheckRule("sub-class-name", "Sub-class name should follow naming conventions", `^(_?[A-Z][a-zA-Z0-9]*)$`),
		NewNameCheckRule("class-name", "Class name should follow naming conventions", `^[A-Z][a-zA-Z0-9]*$`),
		NewNameCheckRule("signal-name", "Signal name should follow naming conventions", `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`),
		NewNameCheckRule("enum-name", "Enum name should follow naming conventions", `^[A-Z][a-zA-Z0-9]*$`),
		NewNameCheckRule("enum-element-name", "Enum element name should follow naming conventions", `^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$`),
		NewNameCheckRule("loop-variable-name", "Loop variable name should follow naming conventions", `^(_|[a-z])[a-z0-9]*(_[a-z0-9]+)*$`),
		NewNameCheckRule("function-argument-name", "Function argument name should follow naming conventions", `^(_|[a-z])[a-z0-9]*(_[a-z0-9]+)*$`),
		NewNameCheckRule("function-variable-name", "Function variable name should follow naming conventions", `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`),
		NewNameCheckRule("function-preload-variable-name", "Function preload variable name should follow naming conventions", `^[A-Z][a-zA-Z0-9]*$`),
		NewNameCheckRule("constant-name", "Constant name should follow naming conventions", `^(_?[A-Z][A-Z0-9]*(_[A-Z0-9]+)*)$`),
		NewNameCheckRule("load-constant-name", "Load constant name should follow naming conventions", `^[A-Z][a-zA-Z0-9]*$`),
		NewNameCheckRule("class-variable-name", "Class variable name should follow naming conventions", `^(_?[a-z][a-z0-9]*(_[a-z0-9]+)*)$`),
		NewNameCheckRule("class-load-variable-name", "Class load variable name should follow naming conventions", `^(_?[a-z][a-z0-9]*(_[a-z0-9]+)*|[A-Z][a-zA-Z0-9]*)$`),
	}
}
