package rules

import (
	"github.com/dzannotti/gdtoolkit/internal/core/ast"
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
	"github.com/dzannotti/gdtoolkit/internal/core/linter/problem"
)

// ClassDefinitionsOrder checks for proper class member ordering
type ClassDefinitionsOrder struct{}

// Name returns the name of the rule
func (r *ClassDefinitionsOrder) Name() string {
	return "class-definitions-order"
}

// Description returns a description of the rule
func (r *ClassDefinitionsOrder) Description() string {
	return "Checks for proper class member ordering"
}

// Check applies the rule to an AST and returns any problems found
func (r *ClassDefinitionsOrder) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Create a visitor to check class definitions order
	visitor := &classDefinitionsOrderVisitor{
		problems: &problems,
	}

	// Walk the AST
	ast.Walk(visitor, tree)

	return problems
}

// classDefinitionsOrderVisitor is a visitor that checks class member ordering
type classDefinitionsOrderVisitor struct {
	problems *[]problem.Problem
}

// Visit is called for each node in the AST
func (v *classDefinitionsOrderVisitor) Visit(node ast.Node) ast.Visitor {
	if class, ok := node.(*ast.Class); ok {
		v.checkClassOrder(class)
	}
	return v
}

// Expected order for class members according to gdtoolkit
type memberType int

const (
	memberPass memberType = iota
	memberClassName
	memberExtends
	memberDocstring
	memberSignal
	memberEnum
	memberConst
	memberStaticVar
	memberExportGroup
	memberExportVar
	memberVar
	memberPrivateVar
	memberOnreadyVar
	memberPrivateOnreadyVar
	memberInnerClass
	memberFunc
	memberStaticFunc
)

// getMemberType determines the type of a class member
func getMemberType(stmt ast.Statement) memberType {
	switch s := stmt.(type) {
	case *ast.PassStatement:
		return memberPass
	case *ast.VarStatement:
		// Check if it's a const
		if s.IsConst {
			return memberConst
		}

		// Check annotations for @export, @onready, etc.
		hasExport := false
		hasOnready := false
		for _, annotation := range s.Annotations {
			if annotation.Name == "export" {
				hasExport = true
			} else if annotation.Name == "export_group" {
				return memberExportGroup
			} else if annotation.Name == "onready" {
				hasOnready = true
			}
		}

		// Check if variable name starts with underscore (private)
		isPrivate := len(s.Name) > 0 && s.Name[0] == '_'

		if hasExport {
			return memberExportVar
		} else if hasOnready {
			if isPrivate {
				return memberPrivateOnreadyVar
			}
			return memberOnreadyVar
		} else if isPrivate {
			return memberPrivateVar
		}
		return memberVar
	case *ast.ExpressionStatement:
		// Check for string literals (docstrings)
		if _, ok := s.Expression.(*ast.StringLiteral); ok {
			return memberDocstring
		}
	case *ast.Class:
		return memberInnerClass
	case *ast.Function:
		// Check if it's a static function
		for _, annotation := range s.Annotations {
			if annotation.Name == "static" {
				return memberStaticFunc
			}
		}
		return memberFunc
	}

	// For other statement types, try to infer from kind
	if stmt, ok := stmt.(*ast.BaseStatement); ok {
		switch stmt.Kind {
		case "signal_stmt":
			return memberSignal
		case "enum_stmt":
			return memberEnum
		}
	}

	return memberVar // Default fallback
}

// checkClassOrder validates the ordering of class members
func (v *classDefinitionsOrderVisitor) checkClassOrder(class *ast.Class) {
	if len(class.Statements) <= 1 {
		return // No ordering issues with 0 or 1 statements
	}

	var lastType memberType = -1

	for _, stmt := range class.Statements {
		currentType := getMemberType(stmt)

		// Check if current member type should come after the last type
		if int(currentType) < int(lastType) {
			*v.problems = append(*v.problems, problem.NewWarning(
				stmt.Position(),
				"Class member is not in the correct order",
				"class-definitions-order",
			))
		}

		lastType = currentType
	}
}

// SubClassBeforeParentClass checks for subclasses defined before their parent class
type SubClassBeforeParentClass struct{}

// Name returns the name of the rule
func (r *SubClassBeforeParentClass) Name() string {
	return "sub-class-before-parent-class"
}

// Description returns a description of the rule
func (r *SubClassBeforeParentClass) Description() string {
	return "Checks for subclasses defined before their parent class"
}

// Check applies the rule to an AST and returns any problems found
func (r *SubClassBeforeParentClass) Check(tree *ast.AbstractSyntaxTree, config linter.Config) []problem.Problem {
	var problems []problem.Problem

	// Collect all class names and their definitions
	classDefinitions := make(map[string]*ast.Class)
	var classOrder []string

	// First pass: collect all classes
	visitor := &classCollector{
		classes: classDefinitions,
		order:   &classOrder,
	}
	ast.Walk(visitor, tree)

	// Second pass: check inheritance order
	for _, className := range classOrder {
		class := classDefinitions[className]
		if class.Extends != "" {
			// Check if parent class is defined after this class
			parentIndex := -1
			currentIndex := -1

			for i, name := range classOrder {
				if name == class.Extends {
					parentIndex = i
				}
				if name == className {
					currentIndex = i
				}
			}

			// If parent is defined after current class, it's an error
			if parentIndex > currentIndex && parentIndex != -1 {
				problems = append(problems, problem.NewError(
					class.Position(),
					"Subclass '"+className+"' is defined before its parent class '"+class.Extends+"'",
					"sub-class-before-parent-class",
				))
			}
		}
	}

	return problems
}

// classCollector collects all class definitions in order
type classCollector struct {
	classes map[string]*ast.Class
	order   *[]string
}

// Visit is called for each node in the AST
func (v *classCollector) Visit(node ast.Node) ast.Visitor {
	if class, ok := node.(*ast.Class); ok {
		if class.Name != "" {
			v.classes[class.Name] = class
			*v.order = append(*v.order, class.Name)
		}
	}
	return v
}
