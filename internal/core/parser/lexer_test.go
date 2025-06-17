package parser

import (
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	input := `
func test():
	var x = 5
	var y = 10
	var result = x + y
	return result
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{NL, "\n"},
		{FUNC, "func"},
		{IDENT, "test"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{COLON, ":"},
		{NL, "\n"},
		{INDENT, ""},
		{VAR, "var"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT, "5"},
		{NL, "\n"},
		{VAR, "var"},
		{IDENT, "y"},
		{ASSIGN, "="},
		{INT, "10"},
		{NL, "\n"},
		{VAR, "var"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{NL, "\n"},
		{RETURN, "return"},
		{IDENT, "result"},
		{NL, "\n"},
		{DEDENT, ""},
		{EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_ComplexTokens(t *testing.T) {
	input := `
class MyClass extends Node:
	var x: int = 5
	var y: float = 3.14
	var s = "hello world"
	var r = r"raw string"
	
	func _init():
		print("Constructor called")
		
	func calculate(a, b: int) -> int:
		if a > b:
			return a
		elif a < b:
			return b
		else:
			return a + b
`

	l := NewLexer(input)

	// Just test a few key tokens to ensure complex structures are lexed correctly
	expectedTokens := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{NL, "\n"},
		{CLASS, "class"},
		{IDENT, "MyClass"},
		{EXTENDS, "extends"},
		{IDENT, "Node"},
		{COLON, ":"},
		{NL, "\n"},
		{INDENT, ""},
		{VAR, "var"},
		{IDENT, "x"},
		{COLON, ":"},
		{IDENT, "int"},
		{ASSIGN, "="},
		{INT, "5"},
		{NL, "\n"},
		{VAR, "var"},
		{IDENT, "y"},
		{COLON, ":"},
		{IDENT, "float"},
		{ASSIGN, "="},
		{FLOAT, "3.14"},
		{NL, "\n"},
	}

	for i, tt := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral && tt.expectedLiteral != "" {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexer_Operators(t *testing.T) {
	input := `+ - * / % ** += -= *= /= %= **= & | ^ ~ << >> &= |= ^= <<= >>= && || ! == != < > <= >= . , : ; ( ) { } [ ] @ $ -> as is`

	expectedTokens := []TokenType{
		PLUS, MINUS, ASTERISK, SLASH, PERCENT, POWER,
		PLUSEQ, MINUSEQ, ASTERISKEQ, SLASHEQ, PERCENTEQ, POWEREQ,
		BITAND, BITOR, BITXOR, BITNOT, LTLT, GTGT,
		AMPEQ, PIPEEQ, CARETEQ, LTLTEQ, GTGTEQ,
		AMPAMP, PIPEPIPE, BANG, EQ, NOT_EQ,
		LT, GT, LTE, GTE,
		DOT, COMMA, COLON, SEMICOLON,
		LPAREN, RPAREN, LBRACE, RBRACE, LBRACKET, RBRACKET,
		AT, DOLLAR, ARROW, AS, IS,
		EOF,
	}

	l := NewLexer(input)

	for i, expected := range expectedTokens {
		tok := l.NextToken()

		if tok.Type != expected {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, expected, tok.Type)
		}
	}
}
