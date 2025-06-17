package parser

import (
	"fmt"
	"path/filepath"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
)

// Parser represents a parser for GDScript
type Parser struct {
	lexer        *Lexer
	currentToken Token
	peekToken    Token
	errors       []error
	errorMode    ErrorMode
}

// Error represents a parser error
type Error struct {
	Line    int
	Column  int
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// ErrorMode defines how the parser handles errors
type ErrorMode int

const (
	// ErrorModeStrict stops parsing at the first error
	ErrorModeStrict ErrorMode = iota
	// ErrorModePanic uses panic mode recovery to continue parsing after errors
	ErrorModePanic
)

// NewParser creates a new parser
func NewParser(input string) *Parser {
	lexer := NewLexer(input)
	p := &Parser{
		lexer:  lexer,
		errors: []error{},
	}
	// Read two tokens to initialize currentToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// NewParserWithOptions creates a new parser with the given options
func NewParserWithOptions(input string, errorMode ErrorMode) *Parser {
	parser := NewParser(input)
	parser.errorMode = errorMode
	return parser
}

// Errors returns the parser errors
func (p *Parser) Errors() []error {
	return p.errors
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// expectPeek checks if the next token is of the expected type
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)

	// If in panic mode, try to recover
	if p.errorMode == ErrorModePanic {
		p.synchronize()
	}

	return false
}

// skipNewlines advances past any newline and indentation tokens
func (p *Parser) skipNewlines() {
	for p.currentToken.Type == NL || p.currentToken.Type == INDENT {
		p.nextToken()
	}
}

// synchronize attempts to recover from a parse error by advancing to a synchronization point
func (p *Parser) synchronize() {
	p.nextToken()

	for p.currentToken.Type != EOF {
		// Synchronization points are typically statement boundaries
		if p.currentToken.Type == SEMICOLON || p.currentToken.Type == NL {
			return
		}

		// Synchronize at the beginning of statements
		switch p.currentToken.Type {
		case CLASS, FUNC, VAR, CONST, IF, FOR, WHILE, MATCH, RETURN:
			return
		}

		p.nextToken()
	}
}

// peekError adds an error for an unexpected token
func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, Error{
		Line:    p.peekToken.Line,
		Column:  p.peekToken.Column,
		Message: msg,
	})
}

// Parse parses the input and returns an AST
func (p *Parser) Parse() *ast.AbstractSyntaxTree {
	tree := ast.NewAST()

	// Parse the global scope
	class := p.parseGlobalScope()
	if class != nil {
		tree.RootClass = class
		tree.Classes = append(tree.Classes, class)

		// Add any subclasses found in global scope to the tree
		for _, subClass := range class.SubClasses {
			tree.Classes = append(tree.Classes, subClass)
		}
	}

	return tree
}

// parseGlobalScope parses the global scope as a class
func (p *Parser) parseGlobalScope() *ast.Class {
	class := ast.NewClass("global scope", ast.Position{Line: 1, Column: 1})

	// Parse statements until EOF
	for p.currentToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			// Check if this is a function and add it to both statements and functions
			if function, ok := stmt.(*ast.Function); ok {
				class.AddFunction(function)
			} else if classStmt, ok := stmt.(*ast.Class); ok {
				// Add class as a subclass to the global scope
				class.AddSubClass(classStmt)
			} else {
				class.AddStatement(stmt)
			}
		}
		p.nextToken()
	}

	return class
}

// parseStatement parses a statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case PASS:
		return p.parsePassStatement()
	case RETURN:
		return p.parseReturnStatement()
	case VAR:
		return p.parseVarStatement()
	case CONST:
		return p.parseConstStatement()
	case FUNC:
		return p.parseFunctionDefinition()
	case CLASS:
		return p.parseClassDefinition()
	case IF:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForStatement()
	case MATCH:
		return p.parseMatchStatement()
	case BREAK:
		return p.parseBreakStatement()
	case CONTINUE:
		return p.parseContinueStatement()
	case IDENT, INT, FLOAT, STRING, RSTRING, TRUE, FALSE, NULL, SELF, LPAREN:
		return p.parseExpressionStatement()
	default:
		return nil
	}
}

// parsePassStatement parses a pass statement
func (p *Parser) parsePassStatement() ast.Statement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}
	return ast.NewPassStatement(pos)
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	stmt := ast.NewReturnStatement(pos, nil)
	p.nextToken()

	// If the next token is not a semicolon or newline, parse the return value
	if p.currentToken.Type != SEMICOLON && p.currentToken.Type != NL {
		stmt.Value = p.parseExpression(PREC_LOWEST)
	}

	return stmt
}

// parseBreakStatement parses a break statement
func (p *Parser) parseBreakStatement() *ast.BreakStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}
	return ast.NewBreakStatement(pos)
}

// parseContinueStatement parses a continue statement
func (p *Parser) parseContinueStatement() *ast.ContinueStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}
	return ast.NewContinueStatement(pos)
}

// parseExpressionStatement parses an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	stmt := ast.NewExpressionStatement(pos, p.parseExpression(PREC_LOWEST))

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// Precedence levels for operators
const (
	PREC_LOWEST     = iota
	PREC_ASSIGN     // =, +=, -=, etc.
	PREC_LOGICAL    // and, or
	PREC_COMPARISON // ==, !=, <, >, <=, >=
	PREC_BITWISE    // &, |, ^
	PREC_SUM        // +, -
	PREC_PRODUCT    // *, /, %
	PREC_PREFIX     // -X, !X
	PREC_CALL       // myFunction(X)
	PREC_INDEX      // array[index]
	PREC_DOT        // obj.property
)

// Operator precedence map
var precedences = map[TokenType]int{
	ASSIGN:     PREC_ASSIGN,
	PLUSEQ:     PREC_ASSIGN,
	MINUSEQ:    PREC_ASSIGN,
	ASTERISKEQ: PREC_ASSIGN,
	SLASHEQ:    PREC_ASSIGN,
	PERCENTEQ:  PREC_ASSIGN,
	AMPEQ:      PREC_ASSIGN,
	PIPEEQ:     PREC_ASSIGN,
	CARETEQ:    PREC_ASSIGN,
	LTLTEQ:     PREC_ASSIGN,
	GTGTEQ:     PREC_ASSIGN,
	POWEREQ:    PREC_ASSIGN,

	OR:       PREC_LOGICAL,
	AND:      PREC_LOGICAL,
	PIPEPIPE: PREC_LOGICAL,
	AMPAMP:   PREC_LOGICAL,

	EQ:     PREC_COMPARISON,
	NOT_EQ: PREC_COMPARISON,
	LT:     PREC_COMPARISON,
	GT:     PREC_COMPARISON,
	LTE:    PREC_COMPARISON,
	GTE:    PREC_COMPARISON,

	BITAND: PREC_BITWISE,
	BITOR:  PREC_BITWISE,
	BITXOR: PREC_BITWISE,
	LTLT:   PREC_BITWISE,
	GTGT:   PREC_BITWISE,

	PLUS:  PREC_SUM,
	MINUS: PREC_SUM,

	ASTERISK: PREC_PRODUCT,
	SLASH:    PREC_PRODUCT,
	PERCENT:  PREC_PRODUCT,
	POWER:    PREC_PRODUCT,

	LPAREN:   PREC_CALL,
	LBRACKET: PREC_INDEX,
	DOT:      PREC_DOT,
}

// parseExpression parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var leftExp ast.Expression

	// Parse prefix expressions
	switch p.currentToken.Type {
	case IDENT:
		leftExp = p.parseIdentifier()
	case INT:
		leftExp = p.parseIntegerLiteral()
	case FLOAT:
		leftExp = p.parseFloatLiteral()
	case STRING:
		leftExp = p.parseStringLiteral()
	case TRUE, FALSE:
		leftExp = p.parseBooleanLiteral()
	case NULL:
		leftExp = p.parseNullLiteral()
	case LPAREN:
		leftExp = p.parseGroupedExpression()
	case MINUS, BANG:
		leftExp = p.parsePrefixExpression()
	default:
		return nil
	}

	if leftExp == nil {
		return nil
	}

	// Parse infix expressions
	for p.peekToken.Type != SEMICOLON && p.peekToken.Type != NL && p.peekToken.Type != EOF &&
		p.peekToken.Type != COLON && p.peekToken.Type != RPAREN && p.peekToken.Type != COMMA &&
		p.peekToken.Type != RBRACE && p.peekToken.Type != RBRACKET {

		// Special case: conditional expression (ternary operator)
		if p.peekToken.Type == IF {
			leftExp = p.parseConditionalExpression(leftExp)
			continue
		}

		// Check if this is an infix operator and if precedence is higher
		peekPrec := p.peekPrecedence()
		if precedence >= peekPrec {
			break
		}

		// Only continue if we have a valid infix operator
		isInfixOp := false
		switch p.peekToken.Type {
		case PLUS, MINUS, ASTERISK, SLASH, PERCENT, POWER:
			isInfixOp = true
		case EQ, NOT_EQ, LT, GT, LTE, GTE:
			isInfixOp = true
		case AND, OR:
			isInfixOp = true
		case ASSIGN, PLUSEQ, MINUSEQ, ASTERISKEQ, SLASHEQ:
			isInfixOp = true
		case LPAREN:
			isInfixOp = true
		case LBRACKET:
			isInfixOp = true
		case DOT:
			isInfixOp = true
		}

		if !isInfixOp {
			break
		}

		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	}

	return leftExp
}

// parseCallExpression parses function call expressions
func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		Function: fn,
	}

	exp.Arguments = p.parseCallArguments()
	return exp
}

// parseCallArguments parses function call arguments
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(PREC_LOWEST))

	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(PREC_LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return args
}

// parseIdentifier parses an identifier expression
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		Value: p.currentToken.Literal,
	}
}

// parseIntegerLiteral parses an integer literal
func (p *Parser) parseIntegerLiteral() ast.Expression {
	return &ast.NumberLiteral{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		IsInt:    true,
		Original: p.currentToken.Literal,
	}
}

// parseFloatLiteral parses a float literal
func (p *Parser) parseFloatLiteral() ast.Expression {
	return &ast.NumberLiteral{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		IsInt:    false,
		Original: p.currentToken.Literal,
	}
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		Value: p.currentToken.Literal,
	}
}

// parseBooleanLiteral parses a boolean literal
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		Value: p.currentToken.Type == TRUE,
	}
}

// parseNullLiteral parses a null literal
func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
	}
}

// parseGroupedExpression parses a grouped expression (parentheses)
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // consume '('

	exp := p.parseExpression(PREC_LOWEST)
	if exp == nil {
		return nil
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp
}

// parsePrefixExpression parses prefix expressions like -x, !x
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREC_PREFIX)

	return expression
}

// parseInfixExpression parses infix expressions like x + y and function calls
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// Handle function calls
	if p.currentToken.Type == LPAREN {
		return p.parseCallExpression(left)
	}

	// Handle regular infix expressions
	expression := &ast.InfixExpression{
		BaseExpression: ast.BaseExpression{
			Pos: ast.Position{
				Line:   p.currentToken.Line,
				Column: p.currentToken.Column,
				Offset: p.currentToken.Offset,
			},
		},
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return PREC_LOWEST
}

// curPrecedence returns the precedence of the current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return PREC_LOWEST
}

// parseVarStatement parses a variable declaration statement
func (p *Parser) parseVarStatement() *ast.VarStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'var' token
	p.nextToken()

	// Skip any newlines between 'var' and variable name
	p.skipNewlines()

	// Parse variable name
	if p.currentToken.Type != IDENT {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected identifier, got %s", p.currentToken.Type),
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	varName := p.currentToken.Literal
	isTyped := false
	stmt := ast.NewVarStatement(pos, varName, isTyped)

	// Move past the identifier
	p.nextToken()

	// Check for type hint or type inference
	if p.currentToken.Type == COLON {
		isTyped = true
		p.nextToken() // Skip colon

		if p.currentToken.Type != IDENT {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected type identifier, got %s", p.currentToken.Type),
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			return nil
		}

		typeHint := p.currentToken.Literal
		stmt.SetTypeHint(typeHint)
		p.nextToken() // Skip type identifier
	} else if p.currentToken.Type == COLONASSIGN {
		// Type inference (:=)
		isTyped = true
		stmt.SetInferred() // We need to add this method
		p.nextToken()      // Skip :=

		// Parse the value expression for type inference
		value := p.parseExpression(PREC_LOWEST)
		stmt.SetValue(value)
		return stmt
	}

	// Check for assignment
	if p.currentToken.Type == ASSIGN {
		p.nextToken() // Skip equals sign
		value := p.parseExpression(PREC_LOWEST)
		stmt.SetValue(value)

		// After parsing expression, advance past it
		p.nextToken()
	}

	return stmt
}

// parseConstStatement parses a constant declaration statement
func (p *Parser) parseConstStatement() *ast.VarStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'const' token
	p.nextToken()

	// Skip any newlines between 'const' and constant name
	p.skipNewlines()

	// Parse constant name
	if p.currentToken.Type != IDENT {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected identifier, got %s", p.currentToken.Type),
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	constName := p.currentToken.Literal
	isTyped := false
	stmt := ast.NewVarStatement(pos, constName, isTyped)
	stmt.SetConst()

	// Move past the identifier
	p.nextToken()

	// Check for type hint
	if p.currentToken.Type == COLON {
		isTyped = true
		p.nextToken() // Skip colon

		if p.currentToken.Type != IDENT {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected type identifier, got %s", p.currentToken.Type),
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			return nil
		}

		typeHint := p.currentToken.Literal
		stmt.SetTypeHint(typeHint)
		p.nextToken() // Skip type identifier
	}

	// Constants must have an assignment
	if p.currentToken.Type != ASSIGN {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "constants must be assigned a value",
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	p.nextToken() // Skip equals sign
	value := p.parseExpression(PREC_LOWEST)
	stmt.SetValue(value)

	// After parsing expression, advance past it
	p.nextToken()

	return stmt
}

// parseFunctionDefinition parses a function definition
func (p *Parser) parseFunctionDefinition() *ast.Function {
	startPos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Check if function is static
	isStatic := false
	if p.currentToken.Type == STATIC {
		isStatic = true
		p.nextToken() // Skip 'static'
		if p.currentToken.Type != FUNC {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected 'func' after 'static', got %s", p.currentToken.Type),
			})
			return nil
		}
	}

	// Skip 'func' token
	p.nextToken()

	// Parse function name - must be immediately after 'func'
	if p.currentToken.Type != IDENT {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected function name, got %s", p.currentToken.Type),
		})
		return nil
	}

	funcName := p.currentToken.Literal
	function := ast.NewFunction(funcName, startPos)
	function.IsStatic = isStatic

	// Move past the function name
	p.nextToken()

	// Parse parameters - must have parentheses
	if p.currentToken.Type != LPAREN {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected '(' after function name, got %s", p.currentToken.Type),
		})
		return nil
	}

	p.parseParameterList(function)

	// Check for return type - FIXED: Check current token, not peek
	if p.currentToken.Type == ARROW {
		p.nextToken() // Skip arrow to get to return type

		if p.currentToken.Type != IDENT {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected return type, got %s", p.currentToken.Type),
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			return nil
		}

		function.ReturnType = p.currentToken.Literal
		p.nextToken() // Skip return type
	}

	// Expect colon
	if p.currentToken.Type != COLON {
		if !p.expectPeek(COLON) {
			return nil
		}
	}

	// Parse function body - handle flexible indentation
	// Skip any newlines
	for p.peekToken.Type == NL || p.peekToken.Type == INDENT {
		p.nextToken()
	}

	// Try to expect indentation, but if we find statement tokens, proceed anyway
	if p.peekToken.Type == INDENT {
		p.nextToken() // consume INDENT
	} else if p.peekToken.Type == VAR || p.peekToken.Type == RETURN || p.peekToken.Type == PASS {
		// If we find statement tokens directly, the lexer handled indentation implicitly
		// This is acceptable - proceed with statement parsing
	} else if p.peekToken.Type == EOF {
		// If we reach EOF, there's no body to parse
		return function
	} else {
		p.errors = append(p.errors, Error{
			Line:    p.peekToken.Line,
			Column:  p.peekToken.Column,
			Message: fmt.Sprintf("expected indentation or statement, got %s", p.peekToken.Type),
		})
		return nil
	}

	// Parse statements until dedent
	for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
		// Skip any newlines or indentation tokens before parsing statements
		if p.currentToken.Type == NL {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			function.AddStatement(stmt)
		}
		p.nextToken()
	}

	return function
}

// parseParameterList parses the parameter list of a function
func (p *Parser) parseParameterList(function *ast.Function) {
	p.nextToken() // Skip '('

	if p.currentToken.Type == RPAREN {
		p.nextToken() // Skip ')'
		return        // Empty parameter list
	}

	// Parse first parameter
	param := p.parseParameter()
	if param == nil {
		return // Error already reported by parseParameter
	}
	function.AddParameter(param)

	// Parse additional parameters
	for p.currentToken.Type == COMMA {
		p.nextToken() // Skip ','

		// Check for trailing comma
		if p.currentToken.Type == RPAREN {
			break
		}

		param = p.parseParameter()
		if param == nil {
			return // Error already reported by parseParameter
		}
		function.AddParameter(param)
	}

	// Must end with right parenthesis
	if p.currentToken.Type != RPAREN {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected ')' at end of parameter list",
		})
		return
	}
	p.nextToken() // Skip ')'
}

// parseParameter parses a single parameter
func (p *Parser) parseParameter() *ast.Parameter {
	// Parameter must start with an identifier
	if p.currentToken.Type != IDENT {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected parameter name, got %s", p.currentToken.Type),
		})
		return nil
	}

	param := &ast.Parameter{
		Pos: ast.Position{
			Line:   p.currentToken.Line,
			Column: p.currentToken.Column,
			Offset: p.currentToken.Offset,
		},
		Name: p.currentToken.Literal,
	}

	p.nextToken() // Skip parameter name

	// Handle type hint if present
	if p.currentToken.Type == COLON {
		p.nextToken() // Skip ':'
		if p.currentToken.Type != IDENT {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected type name after ':', got %s", p.currentToken.Type),
			})
			return nil
		}
		param.TypeHint = p.currentToken.Literal
		p.nextToken() // Skip type name
	}

	// Handle default value if present
	if p.currentToken.Type == ASSIGN {
		p.nextToken() // Skip '='
		param.Default = p.parseExpression(PREC_LOWEST)
		if param.Default == nil {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: "invalid default value expression",
			})
			return nil
		}
		// After parsing default value expression, advance past it
		p.nextToken()
	}

	return param
}

// parseClassDefinition parses a class definition
func (p *Parser) parseClassDefinition() *ast.Class {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'class' token
	p.nextToken()

	// Skip any newlines between 'class' and class name
	p.skipNewlines()

	// Parse class name
	if p.currentToken.Type != IDENT {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected class name, got %s", p.currentToken.Type),
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	className := p.currentToken.Literal
	class := ast.NewClass(className, pos)

	// Move past the class name
	p.nextToken()

	// Check for extends
	if p.currentToken.Type == EXTENDS {
		p.nextToken() // Skip extends keyword

		if p.currentToken.Type != IDENT {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected parent class name, got %s", p.currentToken.Type),
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			return nil
		}

		class.Extends = p.currentToken.Literal
		p.nextToken() // Skip parent class name
	}

	// Expect colon
	if p.currentToken.Type != COLON {
		if !p.expectPeek(COLON) {
			return nil
		}
	}

	// Parse class body - handle flexible indentation
	// Skip any newlines
	for p.peekToken.Type == NL {
		p.nextToken()
	}

	// Try to expect indentation, but if we find statement tokens, proceed anyway
	if p.peekToken.Type == INDENT {
		p.nextToken() // consume INDENT
	} else if p.peekToken.Type == VAR || p.peekToken.Type == FUNC || p.peekToken.Type == PASS || p.peekToken.Type == CLASS {
		// If we find statement tokens directly, the lexer handled indentation implicitly
		// This is acceptable - proceed with statement parsing
	} else {
		p.errors = append(p.errors, Error{
			Line:    p.peekToken.Line,
			Column:  p.peekToken.Column,
			Message: fmt.Sprintf("expected indentation or statement, got %s", p.peekToken.Type),
		})
		return nil
	}

	// Parse class body statements until dedent
	for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
		// Skip any newlines or indentation tokens before parsing statements
		if p.currentToken.Type == NL || p.currentToken.Type == INDENT {
			p.nextToken()
			continue
		}

		switch p.currentToken.Type {
		case FUNC:
			function := p.parseFunctionDefinition()
			if function != nil {
				class.AddFunction(function)
			}
		case CLASS:
			subClass := p.parseClassDefinition()
			if subClass != nil {
				class.AddSubClass(subClass)
			}
		default:
			stmt := p.parseStatement()
			if stmt != nil {
				class.AddStatement(stmt)
			}
		}
		p.nextToken()
	}

	return class
}

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() *ast.IfStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'if' token
	p.nextToken()

	// Parse condition
	condition := p.parseExpression(PREC_LOWEST)
	if condition == nil {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected condition expression after 'if'",
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	stmt := ast.NewIfStatement(pos, condition)

	// Expect colon
	if !p.expectPeek(COLON) {
		return nil
	}

	// Parse if body - handle flexible indentation
	// Skip any newlines
	for p.peekToken.Type == NL {
		p.nextToken()
	}

	// Try to expect indentation, but if we find statement tokens, proceed anyway
	if p.peekToken.Type == INDENT {
		p.nextToken() // consume INDENT
	} else if p.peekToken.Type == VAR || p.peekToken.Type == RETURN || p.peekToken.Type == PASS ||
		p.peekToken.Type == IF || p.peekToken.Type == FOR || p.peekToken.Type == WHILE {
		// If we find statement tokens directly, the lexer handled indentation implicitly
		// This is acceptable - proceed with statement parsing
	} else {
		p.errors = append(p.errors, Error{
			Line:    p.peekToken.Line,
			Column:  p.peekToken.Column,
			Message: fmt.Sprintf("expected indentation or statement, got %s", p.peekToken.Type),
		})
		return nil
	}

	// Parse if body
	for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
		// Skip any newlines or indentation tokens before parsing statements
		if p.currentToken.Type == NL || p.currentToken.Type == INDENT {
			p.nextToken()
			continue
		}

		bodyStmt := p.parseStatement()
		if bodyStmt != nil {
			stmt.AddConsequenceStatement(bodyStmt)
		}
		p.nextToken()
	}

	// After the if body, advance past DEDENT to check for elif/else at same indentation level
	if p.currentToken.Type == DEDENT {
		p.nextToken() // Skip DEDENT to get to same level as 'if'
	}

	// Check for elif branches
	for p.currentToken.Type == ELIF {
		p.nextToken() // Skip 'elif'

		// Parse elif condition
		elifCondition := p.parseExpression(PREC_LOWEST)
		if elifCondition == nil {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: "expected condition expression after 'elif'",
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			return nil
		}

		// Expect colon
		if !p.expectPeek(COLON) {
			return nil
		}

		// Parse elif body - handle flexible indentation
		// Skip any newlines
		for p.peekToken.Type == NL {
			p.nextToken()
		}

		// Try to expect indentation, but if we find statement tokens, proceed anyway
		if p.peekToken.Type == INDENT {
			p.nextToken() // consume INDENT
		}

		// Parse elif body
		elifBody := make([]ast.Statement, 0)
		for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
			// Skip any newlines or indentation tokens before parsing statements
			if p.currentToken.Type == NL || p.currentToken.Type == INDENT {
				p.nextToken()
				continue
			}

			elifStmt := p.parseStatement()
			if elifStmt != nil {
				elifBody = append(elifBody, elifStmt)
			}
			p.nextToken()
		}

		stmt.AddElseIfBranch(elifCondition, elifBody)

		// After elif body, advance past DEDENT for next elif/else
		if p.currentToken.Type == DEDENT {
			p.nextToken()
		}
	}

	// Check for else branch
	if p.currentToken.Type == ELSE {
		p.nextToken() // Skip 'else'

		// Expect colon - check current token first
		if p.currentToken.Type == COLON {
			// Already have colon, no need to advance
		} else if p.peekToken.Type == COLON {
			if !p.expectPeek(COLON) {
				return nil
			}
		} else {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected ':' after 'else', got %s", p.currentToken.Type),
			})
			return nil
		}

		// Parse else body - handle flexible indentation like if body
		// Skip any newlines
		for p.peekToken.Type == NL {
			p.nextToken()
		}

		// Try to expect indentation, but if we find statement tokens, proceed anyway
		if p.peekToken.Type == INDENT {
			p.nextToken() // consume INDENT
		} else if p.peekToken.Type == VAR || p.peekToken.Type == RETURN || p.peekToken.Type == PASS ||
			p.peekToken.Type == IF || p.peekToken.Type == FOR || p.peekToken.Type == WHILE || p.peekToken.Type == IDENT {
			// If we find statement tokens directly, the lexer handled indentation implicitly
			// This is acceptable - proceed with statement parsing
		} else {
			p.errors = append(p.errors, Error{
				Line:    p.peekToken.Line,
				Column:  p.peekToken.Column,
				Message: fmt.Sprintf("expected indentation or statement after else, got %s", p.peekToken.Type),
			})
			return nil
		}

		// Parse else body
		elseBody := make([]ast.Statement, 0)
		for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
			// Skip any newlines or indentation tokens before parsing statements
			if p.currentToken.Type == NL || p.currentToken.Type == INDENT {
				p.nextToken()
				continue
			}

			elseStmt := p.parseStatement()
			if elseStmt != nil {
				elseBody = append(elseBody, elseStmt)
			}
			p.nextToken()
		}

		stmt.SetAlternative(elseBody)
	}

	return stmt
}

// parseForStatement parses a for loop statement
func (p *Parser) parseForStatement() *ast.ForStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'for' token
	p.nextToken()

	// Parse iterator variable
	if p.currentToken.Type != IDENT {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected iterator variable name, got %s", p.currentToken.Type),
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	iteratorName := p.currentToken.Literal
	iteratorPos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Move past iterator name
	p.nextToken()

	isTyped := false
	typeHint := ""

	// Check for type hint
	if p.currentToken.Type == COLON {
		isTyped = true
		p.nextToken() // Skip colon

		if p.currentToken.Type != IDENT {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: fmt.Sprintf("expected type identifier, got %s", p.currentToken.Type),
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			return nil
		}

		typeHint = p.currentToken.Literal
		p.nextToken() // Skip type identifier
	}

	// Expect 'in' keyword
	if p.currentToken.Type != IN {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: fmt.Sprintf("expected 'in' keyword, got %s", p.currentToken.Type),
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	// Skip 'in' keyword
	p.nextToken()

	// Parse collection expression
	collection := p.parseExpression(PREC_LOWEST)
	if collection == nil {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected collection expression after 'in'",
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	// Create for statement
	stmt := ast.NewForStatement(pos, iteratorName, iteratorPos, collection, isTyped)

	// Set type hint if provided
	if isTyped {
		stmt.SetTypeHint(typeHint)
	}

	// Expect colon
	if !p.expectPeek(COLON) {
		return nil
	}

	// Parse for body - handle flexible indentation
	// Skip any newlines
	for p.peekToken.Type == NL {
		p.nextToken()
	}

	// Try to expect indentation, but if we find statement tokens, proceed anyway
	if p.peekToken.Type == INDENT {
		p.nextToken() // consume INDENT
	} else if p.peekToken.Type == VAR || p.peekToken.Type == RETURN || p.peekToken.Type == PASS ||
		p.peekToken.Type == IF || p.peekToken.Type == FOR || p.peekToken.Type == WHILE {
		// If we find statement tokens directly, the lexer handled indentation implicitly
		// This is acceptable - proceed with statement parsing
	} else {
		p.errors = append(p.errors, Error{
			Line:    p.peekToken.Line,
			Column:  p.peekToken.Column,
			Message: fmt.Sprintf("expected indentation or statement, got %s", p.peekToken.Type),
		})
		return nil
	}

	// Parse for loop body
	for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
		// Skip any newlines or indentation tokens before parsing statements
		if p.currentToken.Type == NL || p.currentToken.Type == INDENT {
			p.nextToken()
			continue
		}

		bodyStmt := p.parseStatement()
		if bodyStmt != nil {
			stmt.AddBodyStatement(bodyStmt)
		}
		p.nextToken()
	}

	return stmt
}

// parseWhileStatement parses a while loop statement
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'while' token
	p.nextToken()

	// Parse condition
	condition := p.parseExpression(PREC_LOWEST)
	if condition == nil {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected condition expression after 'while'",
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	stmt := ast.NewWhileStatement(pos, condition)

	// Expect colon
	if !p.expectPeek(COLON) {
		return nil
	}

	// Parse while body - handle flexible indentation
	// Skip any newlines
	for p.peekToken.Type == NL {
		p.nextToken()
	}

	// Try to expect indentation, but if we find statement tokens, proceed anyway
	if p.peekToken.Type == INDENT {
		p.nextToken() // consume INDENT
	} else if p.peekToken.Type == VAR || p.peekToken.Type == RETURN || p.peekToken.Type == PASS ||
		p.peekToken.Type == IF || p.peekToken.Type == FOR || p.peekToken.Type == WHILE {
		// If we find statement tokens directly, the lexer handled indentation implicitly
		// This is acceptable - proceed with statement parsing
	} else {
		p.errors = append(p.errors, Error{
			Line:    p.peekToken.Line,
			Column:  p.peekToken.Column,
			Message: fmt.Sprintf("expected indentation or statement, got %s", p.peekToken.Type),
		})
		return nil
	}

	// Parse while loop body
	for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
		// Skip any newlines or indentation tokens before parsing statements
		if p.currentToken.Type == NL || p.currentToken.Type == INDENT {
			p.nextToken()
			continue
		}

		bodyStmt := p.parseStatement()
		if bodyStmt != nil {
			stmt.AddBodyStatement(bodyStmt)
		}
		p.nextToken()
	}

	return stmt
}

// parseMatchStatement parses a match statement
func (p *Parser) parseMatchStatement() *ast.MatchStatement {
	pos := ast.Position{
		Line:   p.currentToken.Line,
		Column: p.currentToken.Column,
		Offset: p.currentToken.Offset,
	}

	// Skip 'match' token
	p.nextToken()

	// Parse value to match
	value := p.parseExpression(PREC_LOWEST)
	if value == nil {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected expression after 'match'",
		})

		if p.errorMode == ErrorModePanic {
			p.synchronize()
		}

		return nil
	}

	stmt := ast.NewMatchStatement(pos, value)

	// Expect colon
	if !p.expectPeek(COLON) {
		return nil
	}

	// Parse match body - handle flexible indentation
	// Skip any newlines
	for p.peekToken.Type == NL {
		p.nextToken()
	}

	// Try to expect indentation, but if we find statement tokens, proceed anyway
	if p.peekToken.Type == INDENT {
		p.nextToken() // consume INDENT
	} else {
		// For match statements, we can proceed directly to parsing patterns
		// Move to the next token to start parsing the first pattern
		p.nextToken()
	}

	// Parse match branches
	for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
		// Parse pattern
		pattern := p.parseExpression(PREC_LOWEST)
		if pattern == nil {
			p.errors = append(p.errors, Error{
				Line:    p.currentToken.Line,
				Column:  p.currentToken.Column,
				Message: "expected pattern expression in match branch",
			})

			if p.errorMode == ErrorModePanic {
				p.synchronize()
			}

			break
		}

		branch := ast.NewMatchBranch(ast.Position{
			Line:   p.currentToken.Line,
			Column: p.currentToken.Column,
			Offset: p.currentToken.Offset,
		}, pattern)

		// Check for guard condition (when)
		if p.peekToken.Type == WHEN {
			p.nextToken() // Move to 'when'
			p.nextToken() // Skip 'when'

			guard := p.parseExpression(PREC_LOWEST)
			if guard == nil {
				p.errors = append(p.errors, Error{
					Line:    p.currentToken.Line,
					Column:  p.currentToken.Column,
					Message: "expected guard expression after 'when'",
				})

				if p.errorMode == ErrorModePanic {
					p.synchronize()
				}

				break
			}

			branch.SetGuard(guard)
		}

		// Expect colon
		if !p.expectPeek(COLON) {
			break
		}

		// Handle branch body
		if p.peekToken.Type == NL {
			// Multi-line branch body
			p.nextToken() // Skip newline

			// Expect indentation
			if !p.expectPeek(INDENT) {
				break
			}

			// Parse branch body
			for p.currentToken.Type != DEDENT && p.currentToken.Type != EOF {
				bodyStmt := p.parseStatement()
				if bodyStmt != nil {
					branch.AddBodyStatement(bodyStmt)
				}
				p.nextToken()
			}
		} else {
			// Single-line branch body
			p.nextToken() // Move past colon

			bodyStmt := p.parseStatement()
			if bodyStmt != nil {
				branch.AddBodyStatement(bodyStmt)
			}

			// Skip to next line
			for p.currentToken.Type != NL && p.currentToken.Type != EOF {
				p.nextToken()
			}

			if p.currentToken.Type == NL {
				p.nextToken() // Skip newline
			}
		}

		stmt.AddBranch(branch)
	}

	return stmt
}

// ParseFile parses a GDScript file
func ParseFile(filePath string, content string) (*ast.AbstractSyntaxTree, []error) {
	// Use panic mode recovery by default for file parsing
	parser := NewParserWithOptions(content, ErrorModePanic)
	tree := parser.Parse()

	// Set the file name in the AST
	if tree.RootClass != nil {
		tree.RootClass.Name = filepath.Base(filePath)
	}

	return tree, parser.Errors()
}

// parseConditionalExpression parses a conditional expression (ternary operator)
// Format: value_if_true if condition else value_if_false
func (p *Parser) parseConditionalExpression(valueIfTrue ast.Expression) ast.Expression {
	p.nextToken() // consume 'if'

	// Parse the condition
	condition := p.parseExpression(PREC_LOWEST)
	if condition == nil {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected condition expression in conditional expression",
		})
		return nil
	}

	// Expect 'else'
	if !p.expectPeek(ELSE) {
		return nil
	}

	p.nextToken() // consume 'else' and move to value_if_false

	// Parse the value if false
	valueIfFalse := p.parseExpression(PREC_LOWEST)
	if valueIfFalse == nil {
		p.errors = append(p.errors, Error{
			Line:    p.currentToken.Line,
			Column:  p.currentToken.Column,
			Message: "expected expression after 'else' in conditional expression",
		})
		return nil
	}

	// Create a conditional expression node
	return &ast.ConditionalExpression{
		BaseExpression: ast.BaseExpression{
			Pos: valueIfTrue.Position(),
		},
		Condition:    condition,
		ValueIfTrue:  valueIfTrue,
		ValueIfFalse: valueIfFalse,
	}
}
