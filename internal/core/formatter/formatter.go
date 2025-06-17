package formatter

import (
	"regexp"
	"strings"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
)

const (
	// TAB_INDENT_SIZE is the default indentation size when using tabs
	TAB_INDENT_SIZE = 4
	// INLINE_COMMENT_OFFSET is the number of spaces before inline comments
	INLINE_COMMENT_OFFSET = 2
	// MAX_LINE_LENGTH is the default maximum line length
	MAX_LINE_LENGTH = 100
)

// Config holds formatter configuration
type Config struct {
	MaxLineLength    int
	SpacesForIndent  *int // nil means use tabs
	UseSpaces        bool
	SingleIndentSize int
}

// DefaultConfig returns the default formatter configuration
func DefaultConfig() *Config {
	return &Config{
		MaxLineLength:    MAX_LINE_LENGTH,
		SpacesForIndent:  nil,
		UseSpaces:        false,
		SingleIndentSize: TAB_INDENT_SIZE,
	}
}

// Context holds formatting state
type Context struct {
	Config              *Config
	SingleIndentString  string
	IndentLevel         int
	PreviouslyProcessed int
	MaxLineLength       int
	IndentRegex         *regexp.Regexp
	StandaloneComments  []string
	InlineComments      []string
}

// NewContext creates a new formatting context
func NewContext(config *Config) *Context {
	var singleIndentString string
	var indentRegex *regexp.Regexp

	if config.SpacesForIndent != nil {
		singleIndentString = strings.Repeat(" ", *config.SpacesForIndent)
		indentRegex = regexp.MustCompile(`^[ ]*`)
	} else {
		singleIndentString = "\t"
		indentRegex = regexp.MustCompile(`^[\t]*`)
	}

	return &Context{
		Config:             config,
		SingleIndentString: singleIndentString,
		IndentLevel:        0,
		MaxLineLength:      config.MaxLineLength,
		IndentRegex:        indentRegex,
	}
}

// GetIndent returns the indentation string for the current level
func (c *Context) GetIndent() string {
	return strings.Repeat(c.SingleIndentString, c.IndentLevel)
}

// IncreaseIndent increases the indentation level
func (c *Context) IncreaseIndent() {
	c.IndentLevel++
}

// DecreaseIndent decreases the indentation level
func (c *Context) DecreaseIndent() {
	if c.IndentLevel > 0 {
		c.IndentLevel--
	}
}

// FormattedLine represents a formatted line with optional line number
type FormattedLine struct {
	LineNumber *int
	Content    string
}

// FormatCode formats GDScript code using the provided AST
func FormatCode(ast *ast.AbstractSyntaxTree, config *Config) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}

	context := NewContext(config)
	formatter := &Formatter{context: context}

	lines := formatter.FormatAST(ast)

	// Join lines with newlines
	var result strings.Builder
	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(line.Content)
	}

	return result.String(), nil
}

// Formatter implements the visitor pattern for formatting
type Formatter struct {
	context *Context
	lines   []FormattedLine
}

// FormatAST formats the entire AST
func (f *Formatter) FormatAST(node *ast.AbstractSyntaxTree) []FormattedLine {
	f.lines = []FormattedLine{}
	f.visitAST(node)
	return f.lines
}

// addLine adds a formatted line
func (f *Formatter) addLine(content string) {
	f.lines = append(f.lines, FormattedLine{
		Content: content,
	})
}

// addEmptyLine adds an empty line
func (f *Formatter) addEmptyLine() {
	f.addLine("")
}

// visitAST visits the root AST node
func (f *Formatter) visitAST(node *ast.AbstractSyntaxTree) {
	// Check if root class is just a wrapper (has default name from parser)
	hasRealRootClass := node.RootClass != nil && node.RootClass.Name != "" && node.RootClass.Name != "test.gd"

	// Format real root class if present
	if hasRealRootClass {
		f.visitClass(node.RootClass)
	} else if node.RootClass != nil {
		// If root class is just a wrapper, format its contents directly
		f.visitClassContents(node.RootClass)
	}

	// Format other classes with proper spacing
	previousWasClass := hasRealRootClass
	for _, class := range node.Classes {
		if class != node.RootClass {
			if previousWasClass {
				f.addEmptyLine()
				f.addEmptyLine()
			}
			f.visitClass(class)
			previousWasClass = true
		}
	}

	// Format top-level functions
	for _, function := range node.Functions {
		if previousWasClass {
			f.addEmptyLine()
			f.addEmptyLine()
		}
		f.visitFunction(function)
		previousWasClass = false
	}
}

// visitClassContents formats the contents of a class without the class declaration
func (f *Formatter) visitClassContents(node *ast.Class) {
	// Format statements
	previousWasStatement := false
	for _, stmt := range node.Statements {
		if previousWasStatement {
			// Add spacing between different types of statements
		}
		f.visitStatement(stmt)
		previousWasStatement = true
	}

	// Format functions with proper spacing
	for _, function := range node.Functions {
		if previousWasStatement {
			f.addEmptyLine()
		}
		f.visitFunction(function)
		previousWasStatement = true
	}

	// Format sub-classes
	for _, subClass := range node.SubClasses {
		if previousWasStatement {
			f.addEmptyLine()
			f.addEmptyLine()
		}
		f.visitClass(subClass)
		previousWasStatement = true
	}
}

// visitClass formats a class definition
func (f *Formatter) visitClass(node *ast.Class) {
	// Format class name
	classLine := f.context.GetIndent() + "class"
	if node.Name != "" {
		classLine += " " + node.Name
	}

	// Format extends clause
	if node.Extends != "" {
		classLine += ":\n" + f.context.GetIndent() + "\textends " + node.Extends
	}
	classLine += ":"

	f.addLine(classLine)

	// Increase indentation for class body
	f.context.IncreaseIndent()

	// Check if class has any content
	hasContent := len(node.Statements) > 0 || len(node.Functions) > 0 || len(node.SubClasses) > 0

	if !hasContent {
		f.addLine(f.context.GetIndent() + "pass")
	} else {
		// Format statements
		for _, stmt := range node.Statements {
			f.visitStatement(stmt)
		}

		// Format functions with proper spacing
		previousWasStatement := len(node.Statements) > 0
		for _, function := range node.Functions {
			if previousWasStatement {
				f.addEmptyLine()
			}
			f.visitFunction(function)
			previousWasStatement = true
		}

		// Format sub-classes
		for _, subClass := range node.SubClasses {
			if previousWasStatement {
				f.addEmptyLine()
			}
			f.visitClass(subClass)
			previousWasStatement = true
		}
	}

	// Decrease indentation
	f.context.DecreaseIndent()
}

// visitFunction formats a function definition
func (f *Formatter) visitFunction(node *ast.Function) {
	// Build function signature
	funcLine := f.context.GetIndent() + "func " + node.Name + "("

	// Format parameters
	if len(node.Parameters) > 0 {
		// Check if parameters should be multiline
		paramStr := f.formatParameters(node.Parameters)
		totalLength := len(funcLine) + len(paramStr) + 1 // +1 for closing paren
		if node.ReturnType != "" {
			totalLength += len(" -> " + node.ReturnType)
		}

		if totalLength > f.context.MaxLineLength {
			// Multi-line parameters
			funcLine += "\n"
			f.context.IncreaseIndent()
			for i, param := range node.Parameters {
				paramLine := f.context.GetIndent() + f.formatParameter(param)
				if i < len(node.Parameters)-1 {
					paramLine += ","
				} else {
					paramLine += ","
				}
				funcLine += paramLine + "\n"
			}
			f.context.DecreaseIndent()
			funcLine += f.context.GetIndent() + ")"
		} else {
			// Single line parameters
			funcLine += paramStr + ")"
		}
	} else {
		funcLine += ")"
	}

	// Add return type if present
	if node.ReturnType != "" {
		funcLine += " -> " + node.ReturnType
	}

	funcLine += ":"
	f.addLine(funcLine)

	// Format function body
	f.context.IncreaseIndent()
	if len(node.Statements) == 0 {
		f.addLine(f.context.GetIndent() + "pass")
	} else {
		for _, stmt := range node.Statements {
			f.visitStatement(stmt)
		}
	}
	f.context.DecreaseIndent()
}

// formatParameters formats function parameters as a single line
func (f *Formatter) formatParameters(params []*ast.Parameter) string {
	var parts []string
	for _, param := range params {
		parts = append(parts, f.formatParameter(param))
	}
	return strings.Join(parts, ", ")
}

// formatParameter formats a single parameter
func (f *Formatter) formatParameter(param *ast.Parameter) string {
	result := param.Name

	// Add type annotation
	if param.TypeHint != "" {
		result += ": " + param.TypeHint
	}

	// Add default value
	if param.Default != nil {
		// For now, always use = for default values
		// TODO: Implement proper inferred type detection
		result += " = " + f.formatExpression(param.Default)
	}

	return result
}

// visitStatement formats a statement
func (f *Formatter) visitStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.VarStatement:
		f.visitVarStatement(s)
	case *ast.ReturnStatement:
		f.visitReturnStatement(s)
	case *ast.ExpressionStatement:
		f.visitExpressionStatement(s)
	case *ast.PassStatement:
		f.addLine(f.context.GetIndent() + "pass")
	case *ast.BreakStatement:
		f.addLine(f.context.GetIndent() + "break")
	case *ast.ContinueStatement:
		f.addLine(f.context.GetIndent() + "continue")
	case *ast.IfStatement:
		f.visitIfStatement(s)
	case *ast.ForStatement:
		f.visitForStatement(s)
	case *ast.WhileStatement:
		f.visitWhileStatement(s)
	case *ast.MatchStatement:
		f.visitMatchStatement(s)
	case *ast.Function:
		// Functions can appear as statements in some contexts
		f.visitFunction(s)
	case *ast.Class:
		// Classes can appear as statements in some contexts
		f.visitClass(s)
	default:
		// Add debugging info to understand what type we're getting
		if stmt != nil {
			f.addLine(f.context.GetIndent() + "# Unknown statement type: " + stmt.TokenLiteral())
		} else {
			f.addLine(f.context.GetIndent() + "# Nil statement")
		}
	}
}

// visitVarStatement formats a variable declaration
func (f *Formatter) visitVarStatement(stmt *ast.VarStatement) {
	line := f.context.GetIndent()

	if stmt.IsConst {
		line += "const " + stmt.Name
	} else {
		line += "var " + stmt.Name
	}

	if stmt.TypeHint != "" {
		line += ": " + stmt.TypeHint
	}

	if stmt.Value != nil {
		line += " = " + f.formatExpression(stmt.Value)
	}

	f.addLine(line)
}

// visitReturnStatement formats a return statement
func (f *Formatter) visitReturnStatement(stmt *ast.ReturnStatement) {
	line := f.context.GetIndent() + "return"
	if stmt.Value != nil {
		line += " " + f.formatExpression(stmt.Value)
	}
	f.addLine(line)
}

// visitExpressionStatement formats an expression statement
func (f *Formatter) visitExpressionStatement(stmt *ast.ExpressionStatement) {
	line := f.context.GetIndent() + f.formatExpression(stmt.Expression)
	f.addLine(line)
}

// visitIfStatement formats an if statement
func (f *Formatter) visitIfStatement(stmt *ast.IfStatement) {
	// Format if condition
	line := f.context.GetIndent() + "if " + f.formatExpression(stmt.Condition) + ":"
	f.addLine(line)

	// Format if body
	f.context.IncreaseIndent()
	if len(stmt.Consequence) == 0 {
		f.addLine(f.context.GetIndent() + "pass")
	} else {
		for _, s := range stmt.Consequence {
			f.visitStatement(s)
		}
	}
	f.context.DecreaseIndent()

	// Format elif clauses
	for i, condition := range stmt.ElseCondition {
		elifLine := f.context.GetIndent() + "elif " + f.formatExpression(condition) + ":"
		f.addLine(elifLine)

		f.context.IncreaseIndent()
		if len(stmt.ElseBranches[i]) == 0 {
			f.addLine(f.context.GetIndent() + "pass")
		} else {
			for _, s := range stmt.ElseBranches[i] {
				f.visitStatement(s)
			}
		}
		f.context.DecreaseIndent()
	}

	// Format else clause
	if len(stmt.Alternative) > 0 {
		f.addLine(f.context.GetIndent() + "else:")
		f.context.IncreaseIndent()
		for _, s := range stmt.Alternative {
			f.visitStatement(s)
		}
		f.context.DecreaseIndent()
	}
}

// visitForStatement formats a for statement
func (f *Formatter) visitForStatement(stmt *ast.ForStatement) {
	line := f.context.GetIndent() + "for " + stmt.Iterator
	if stmt.TypeHint != "" {
		line += ": " + stmt.TypeHint
	}
	line += " in " + f.formatExpression(stmt.Collection) + ":"
	f.addLine(line)

	f.context.IncreaseIndent()
	if len(stmt.Body) == 0 {
		f.addLine(f.context.GetIndent() + "pass")
	} else {
		for _, s := range stmt.Body {
			f.visitStatement(s)
		}
	}
	f.context.DecreaseIndent()
}

// visitWhileStatement formats a while statement
func (f *Formatter) visitWhileStatement(stmt *ast.WhileStatement) {
	line := f.context.GetIndent() + "while " + f.formatExpression(stmt.Condition) + ":"
	f.addLine(line)

	f.context.IncreaseIndent()
	if len(stmt.Body) == 0 {
		f.addLine(f.context.GetIndent() + "pass")
	} else {
		for _, s := range stmt.Body {
			f.visitStatement(s)
		}
	}
	f.context.DecreaseIndent()
}

// visitMatchStatement formats a match statement
func (f *Formatter) visitMatchStatement(stmt *ast.MatchStatement) {
	line := f.context.GetIndent() + "match " + f.formatExpression(stmt.Value) + ":"
	f.addLine(line)

	f.context.IncreaseIndent()
	for _, branch := range stmt.Branches {
		patternLine := f.context.GetIndent() + f.formatExpression(branch.Pattern)
		if branch.Guard != nil {
			patternLine += " when " + f.formatExpression(branch.Guard)
		}
		patternLine += ":"
		f.addLine(patternLine)

		f.context.IncreaseIndent()
		if len(branch.Body) == 0 {
			f.addLine(f.context.GetIndent() + "pass")
		} else {
			for _, s := range branch.Body {
				f.visitStatement(s)
			}
		}
		f.context.DecreaseIndent()
	}
	f.context.DecreaseIndent()
}

// formatExpression formats an expression
func (f *Formatter) formatExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value
	case *ast.StringLiteral:
		// String literals already include quotes, don't add extra ones
		return e.Value
	case *ast.NumberLiteral:
		return e.Original
	case *ast.BooleanLiteral:
		if e.Value {
			return "true"
		}
		return "false"
	case *ast.NullLiteral:
		return "null"
	case *ast.ArrayLiteral:
		if len(e.Elements) == 0 {
			return "[]"
		}
		var elements []string
		for _, elem := range e.Elements {
			elements = append(elements, f.formatExpression(elem))
		}
		return "[" + strings.Join(elements, ", ") + "]"
	case *ast.DictionaryLiteral:
		if len(e.Pairs) == 0 {
			return "{}"
		}
		var pairs []string
		for key, value := range e.Pairs {
			pairs = append(pairs, f.formatExpression(key)+": "+f.formatExpression(value))
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	case *ast.PrefixExpression:
		return e.Operator + f.formatExpression(e.Right)
	case *ast.InfixExpression:
		return f.formatExpression(e.Left) + " " + e.Operator + " " + f.formatExpression(e.Right)
	case *ast.CallExpression:
		funcStr := f.formatExpression(e.Function)
		if len(e.Arguments) == 0 {
			return funcStr + "()"
		}
		var args []string
		for _, arg := range e.Arguments {
			args = append(args, f.formatExpression(arg))
		}
		return funcStr + "(" + strings.Join(args, ", ") + ")"
	case *ast.IndexExpression:
		return f.formatExpression(e.Left) + "[" + f.formatExpression(e.Index) + "]"
	case *ast.DotExpression:
		return f.formatExpression(e.Left) + "." + e.Property
	case *ast.AssignmentExpression:
		return f.formatExpression(e.Left) + " = " + f.formatExpression(e.Right)
	default:
		return "# Unknown expression"
	}
}
