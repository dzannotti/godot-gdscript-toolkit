package parser

// TokenType represents the type of a token
type TokenType string

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	Offset  int
}

// NewToken creates a new token
func NewToken(tokenType TokenType, literal string, line, column, offset int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
		Column:  column,
		Offset:  offset,
	}
}

// Token types
const (
	// Special tokens
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	COMMENT TokenType = "COMMENT"
	NL      TokenType = "NL"
	INDENT  TokenType = "INDENT"
	DEDENT  TokenType = "_DEDENT"

	// Identifiers and literals
	IDENT   TokenType = "IDENT"
	INT     TokenType = "INT"
	FLOAT   TokenType = "FLOAT"
	STRING  TokenType = "STRING"
	RSTRING TokenType = "RSTRING"
	HEX     TokenType = "HEX"
	BIN     TokenType = "BIN"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"
	POWER    TokenType = "**"

	LT       TokenType = "<"
	GT       TokenType = ">"
	EQ       TokenType = "=="
	NOT_EQ   TokenType = "!="
	LTE      TokenType = "<="
	GTE      TokenType = ">="
	AND      TokenType = "and"
	OR       TokenType = "or"
	NOT      TokenType = "not"
	AMPAMP   TokenType = "&&"
	PIPEPIPE TokenType = "||"

	PLUSEQ     TokenType = "+="
	MINUSEQ    TokenType = "-="
	ASTERISKEQ TokenType = "*="
	SLASHEQ    TokenType = "/="
	PERCENTEQ  TokenType = "%="
	AMPEQ      TokenType = "&="
	PIPEEQ     TokenType = "|="
	CARETEQ    TokenType = "^="
	LTLTEQ     TokenType = "<<="
	GTGTEQ     TokenType = ">>="
	POWEREQ    TokenType = "**="

	BITAND TokenType = "&"
	BITOR  TokenType = "|"
	BITXOR TokenType = "^"
	BITNOT TokenType = "~"
	LTLT   TokenType = "<<"
	GTGT   TokenType = ">>"

	DOT         TokenType = "."
	COMMA       TokenType = ","
	COLON       TokenType = ":"
	SEMICOLON   TokenType = ";"
	LPAREN      TokenType = "("
	RPAREN      TokenType = ")"
	LBRACE      TokenType = "{"
	RBRACE      TokenType = "}"
	LBRACKET    TokenType = "["
	RBRACKET    TokenType = "]"
	AT          TokenType = "@"
	DOLLAR      TokenType = "$"
	CARET       TokenType = "^"
	ARROW       TokenType = "->"
	COLONASSIGN TokenType = ":="

	// Keywords
	FUNC       TokenType = "func"
	CLASS      TokenType = "class"
	EXTENDS    TokenType = "extends"
	CLASS_NAME TokenType = "class_name"
	VAR        TokenType = "var"
	CONST      TokenType = "const"
	STATIC     TokenType = "static"
	ENUM       TokenType = "enum"
	SIGNAL     TokenType = "signal"
	PASS       TokenType = "pass"
	RETURN     TokenType = "return"
	IF         TokenType = "if"
	ELIF       TokenType = "elif"
	ELSE       TokenType = "else"
	FOR        TokenType = "for"
	WHILE      TokenType = "while"
	BREAK      TokenType = "break"
	CONTINUE   TokenType = "continue"
	MATCH      TokenType = "match"
	WHEN       TokenType = "when"
	IN         TokenType = "in"
	IS         TokenType = "is"
	AS         TokenType = "as"
	AWAIT      TokenType = "await"
	BREAKPOINT TokenType = "breakpoint"
	TRUE       TokenType = "true"
	FALSE      TokenType = "false"
	NULL       TokenType = "null"
	SELF       TokenType = "self"
	GET        TokenType = "get"
	SET        TokenType = "set"
)

// Keywords map
var keywords = map[string]TokenType{
	"func":       FUNC,
	"class":      CLASS,
	"extends":    EXTENDS,
	"class_name": CLASS_NAME,
	"var":        VAR,
	"const":      CONST,
	"static":     STATIC,
	"enum":       ENUM,
	"signal":     SIGNAL,
	"pass":       PASS,
	"return":     RETURN,
	"if":         IF,
	"elif":       ELIF,
	"else":       ELSE,
	"for":        FOR,
	"while":      WHILE,
	"break":      BREAK,
	"continue":   CONTINUE,
	"match":      MATCH,
	"when":       WHEN,
	"in":         IN,
	"is":         IS,
	"as":         AS,
	"await":      AWAIT,
	"breakpoint": BREAKPOINT,
	"true":       TRUE,
	"false":      FALSE,
	"null":       NULL,
	"self":       SELF,
	"get":        GET,
	"set":        SET,
	"and":        AND,
	"or":         OR,
	"not":        NOT,
}

// LookupIdent checks if the given identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
