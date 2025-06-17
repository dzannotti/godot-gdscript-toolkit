package parser

import (
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes GDScript source code
type Lexer struct {
	input        string
	position     int     // current position in input (points to current char)
	readPosition int     // current reading position in input (after current char)
	ch           rune    // current char under examination
	line         int     // current line number
	column       int     // current column number
	indentStack  []int   // stack of indentation levels
	indentLevel  int     // current indentation level
	tokens       []Token // tokens to be returned before continuing lexing
}

// NewLexer creates a new Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:       input,
		line:        1,
		column:      1,
		indentStack: []int{0}, // start with 0 indentation
		indentLevel: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position in the input string
func (l *Lexer) readChar() {
	l.position = l.readPosition
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.readPosition += size
	}

	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0 // EOF
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	// If we have tokens queued up (like INDENT/DEDENT), return them first
	if len(l.tokens) > 0 {
		tok := l.tokens[0]
		l.tokens = l.tokens[1:]
		return tok
	}

	l.skipWhitespace()

	var tok Token

	switch l.ch {
	case '+':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(PLUSEQ, string(ch)+"=")
		} else {
			tok = l.newToken(PLUS, string(l.ch))
		}
	case '-':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(MINUSEQ, string(ch)+"=")
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(ARROW, string(ch)+">")
		} else {
			tok = l.newToken(MINUS, string(l.ch))
		}
	case '*':
		if l.peekChar() == '*' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.newToken(POWEREQ, string(ch)+"*=")
			} else {
				tok = l.newToken(POWER, string(ch)+"*")
			}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(ASTERISKEQ, string(ch)+"=")
		} else {
			tok = l.newToken(ASTERISK, string(l.ch))
		}
	case '/':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(SLASHEQ, string(ch)+"=")
		} else {
			tok = l.newToken(SLASH, string(l.ch))
		}
	case '%':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(PERCENTEQ, string(ch)+"=")
		} else {
			tok = l.newToken(PERCENT, string(l.ch))
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(AMPAMP, string(ch)+"&")
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(AMPEQ, string(ch)+"=")
		} else {
			tok = l.newToken(BITAND, string(l.ch))
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(PIPEPIPE, string(ch)+"|")
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(PIPEEQ, string(ch)+"=")
		} else {
			tok = l.newToken(BITOR, string(l.ch))
		}
	case '^':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(CARETEQ, string(ch)+"=")
		} else {
			tok = l.newToken(CARET, string(l.ch))
		}
	case '~':
		tok = l.newToken(BITNOT, string(l.ch))
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(LTE, string(ch)+"=")
		} else if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.newToken(LTLTEQ, string(ch)+"<=")
			} else {
				tok = l.newToken(LTLT, string(ch)+"<")
			}
		} else {
			tok = l.newToken(LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(GTE, string(ch)+"=")
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = l.newToken(GTGTEQ, string(ch)+">=")
			} else {
				tok = l.newToken(GTGT, string(ch)+">")
			}
		} else {
			tok = l.newToken(GT, string(l.ch))
		}
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(EQ, string(ch)+"=")
		} else {
			tok = l.newToken(ASSIGN, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(NOT_EQ, string(ch)+"=")
		} else {
			tok = l.newToken(BANG, string(l.ch))
		}
	case '.':
		tok = l.newToken(DOT, string(l.ch))
	case ',':
		tok = l.newToken(COMMA, string(l.ch))
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(COLONASSIGN, string(ch)+"=")
		} else {
			tok = l.newToken(COLON, string(l.ch))
		}
	case ';':
		tok = l.newToken(SEMICOLON, string(l.ch))
	case '(':
		tok = l.newToken(LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(RPAREN, string(l.ch))
	case '{':
		tok = l.newToken(LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(RBRACE, string(l.ch))
	case '[':
		tok = l.newToken(LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(RBRACKET, string(l.ch))
	case '@':
		tok = l.newToken(AT, string(l.ch))
	case '$':
		tok = l.newToken(DOLLAR, string(l.ch))
	case '#':
		tok.Type = COMMENT
		tok.Line = l.line
		tok.Column = l.column
		tok.Offset = l.position
		tok.Literal = l.readComment()
		return tok
	case '"', '\'':
		tok.Type = STRING
		tok.Line = l.line
		tok.Column = l.column
		tok.Offset = l.position
		tok.Literal = l.readString(l.ch)
		return tok
	case 'r':
		if l.peekChar() == '"' || l.peekChar() == '\'' {
			startCol := l.column
			startPos := l.position
			l.readChar() // consume 'r'
			quote := l.ch
			l.readChar() // consume quote
			str := l.readRawString(quote)
			tok.Type = RSTRING
			tok.Line = l.line
			tok.Column = startCol
			tok.Offset = startPos
			tok.Literal = "r" + string(quote) + str + string(quote)
			return tok
		}
		// If not a raw string, treat as identifier
		tok.Line = l.line
		tok.Column = l.column
		tok.Offset = l.position
		tok.Literal = l.readIdentifier()
		tok.Type = LookupIdent(tok.Literal)
		return tok
	case '\n':
		// Generate NL token and handle indentation
		tok = l.newToken(NL, "\n")
		l.readChar()

		// Count indentation after newline
		indent := 0
		for l.ch == ' ' || l.ch == '\t' {
			if l.ch == ' ' {
				indent++
			} else if l.ch == '\t' {
				indent += 4 // Assuming tab width of 4
			}
			l.readChar()
		}

		// Handle indentation changes
		l.handleIndentation(indent)
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
		tok.Offset = l.position
		return tok
	default:
		if isLetter(l.ch) || l.ch == '_' {
			tok.Line = l.line
			tok.Column = l.column
			tok.Offset = l.position
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Offset = l.position
			return l.readNumber()
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float)
func (l *Lexer) readNumber() Token {
	position := l.position
	isFloat := false

	// Check for hex or binary prefix
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar() // consume '0'
		l.readChar() // consume 'x'
		for isHexDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
		return l.newToken(HEX, l.input[position:l.position])
	} else if l.ch == '0' && (l.peekChar() == 'b' || l.peekChar() == 'B') {
		l.readChar() // consume '0'
		l.readChar() // consume 'b'
		for isBinDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
		return l.newToken(BIN, l.input[position:l.position])
	}

	// Regular number
	for isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume '.'
		for isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	}

	// Check for exponent
	if (l.ch == 'e' || l.ch == 'E') && (isDigit(l.peekChar()) || l.peekChar() == '+' || l.peekChar() == '-') {
		isFloat = true
		l.readChar() // consume 'e'
		if l.ch == '+' || l.ch == '-' {
			l.readChar() // consume '+' or '-'
		}
		for isDigit(l.ch) || l.ch == '_' {
			l.readChar()
		}
	}

	if isFloat {
		return l.newToken(FLOAT, l.input[position:l.position])
	}
	return l.newToken(INT, l.input[position:l.position])
}

// readString reads a string literal
func (l *Lexer) readString(quote rune) string {
	position := l.position
	l.readChar() // consume opening quote
	for {
		if l.ch == '\\' {
			l.readChar() // consume backslash
			l.readChar() // consume escaped character
			continue
		}
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 1
		}
		l.readChar()
	}
	if l.ch == quote {
		l.readChar() // consume closing quote
	}
	return l.input[position:l.position]
}

// readRawString reads a raw string literal (r"..." or r'...')
func (l *Lexer) readRawString(quote rune) string {
	position := l.position
	for {
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 1
		}
		l.readChar()
	}
	if l.ch == quote {
		l.readChar() // consume closing quote
	}
	return l.input[position : l.position-1]
}

// readComment reads a comment
func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

// skipWhitespace skips whitespace characters (except newlines which are significant)
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// handleIndentation handles indentation changes and generates INDENT/DEDENT tokens
func (l *Lexer) handleIndentation(indent int) {
	currentIndent := l.indentStack[len(l.indentStack)-1]

	if indent > currentIndent {
		// Indentation increased, push to stack and generate INDENT token
		l.indentStack = append(l.indentStack, indent)
		l.tokens = append(l.tokens, l.newToken(INDENT, ""))
	} else if indent < currentIndent {
		// Indentation decreased, pop from stack and generate DEDENT tokens
		for len(l.indentStack) > 0 && indent < l.indentStack[len(l.indentStack)-1] {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			l.tokens = append(l.tokens, l.newToken(DEDENT, ""))
		}

		// Ensure the indentation matches a previous level
		if len(l.indentStack) == 0 || indent != l.indentStack[len(l.indentStack)-1] {
			// Invalid indentation level
			l.tokens = append(l.tokens, l.newToken(ILLEGAL, "invalid indentation"))
		}
	}
}

// newToken creates a new token
func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
		Offset:  l.position,
	}
}

// Helper functions
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isHexDigit(ch rune) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isBinDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}
