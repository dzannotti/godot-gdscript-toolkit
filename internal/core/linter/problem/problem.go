package problem

import (
	"fmt"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
)

// Severity represents the severity of a problem
type Severity string

const (
	// Error represents a critical problem that should be fixed
	Error Severity = "error"
	// Warning represents a non-critical problem that should be fixed
	Warning Severity = "warning"
	// Info represents an informational message
	Info Severity = "info"
)

// Problem represents a linting problem
type Problem struct {
	Position ast.Position
	Message  string
	RuleName string
	Severity Severity
}

// String returns a string representation of the problem
func (p Problem) String() string {
	return fmt.Sprintf("%s: %s (%s) at %s", p.Severity, p.Message, p.RuleName, p.Position)
}

// NewProblem creates a new problem
func NewProblem(pos ast.Position, message, ruleName string, severity Severity) Problem {
	return Problem{
		Position: pos,
		Message:  message,
		RuleName: ruleName,
		Severity: severity,
	}
}

// NewError creates a new error problem
func NewError(pos ast.Position, message, ruleName string) Problem {
	return NewProblem(pos, message, ruleName, Error)
}

// NewWarning creates a new warning problem
func NewWarning(pos ast.Position, message, ruleName string) Problem {
	return NewProblem(pos, message, ruleName, Warning)
}

// NewInfo creates a new info problem
func NewInfo(pos ast.Position, message, ruleName string) Problem {
	return NewProblem(pos, message, ruleName, Info)
}
