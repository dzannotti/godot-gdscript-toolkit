package ast

import (
	"strconv"
)

// BaseExpression provides common functionality for all expressions
type BaseExpression struct {
	Pos Position
}

// Position returns the position of the expression in the source code
func (e *BaseExpression) Position() Position {
	return e.Pos
}

func (e *BaseExpression) expressionNode() {}

// Identifier represents a variable or function name
type Identifier struct {
	BaseExpression
	Value string
}

// TokenLiteral returns the literal value of the token
func (i *Identifier) TokenLiteral() string {
	return i.Value
}

// NewIdentifier creates a new identifier
func NewIdentifier(value string, pos Position) *Identifier {
	return &Identifier{
		BaseExpression: BaseExpression{Pos: pos},
		Value:          value,
	}
}

// StringLiteral represents a string literal
type StringLiteral struct {
	BaseExpression
	Value string
}

// TokenLiteral returns the literal value of the token
func (s *StringLiteral) TokenLiteral() string {
	return s.Value
}

// NewStringLiteral creates a new string literal
func NewStringLiteral(value string, pos Position) *StringLiteral {
	return &StringLiteral{
		BaseExpression: BaseExpression{Pos: pos},
		Value:          value,
	}
}

// NumberLiteral represents a numeric literal
type NumberLiteral struct {
	BaseExpression
	Value    float64
	IsInt    bool
	Original string
}

// TokenLiteral returns the literal value of the token
func (n *NumberLiteral) TokenLiteral() string {
	return n.Original
}

// NewIntLiteral creates a new integer literal
func NewIntLiteral(value int64, original string, pos Position) *NumberLiteral {
	return &NumberLiteral{
		BaseExpression: BaseExpression{Pos: pos},
		Value:          float64(value),
		IsInt:          true,
		Original:       original,
	}
}

// NewFloatLiteral creates a new float literal
func NewFloatLiteral(value float64, original string, pos Position) *NumberLiteral {
	return &NumberLiteral{
		BaseExpression: BaseExpression{Pos: pos},
		Value:          value,
		IsInt:          false,
		Original:       original,
	}
}

// NullLiteral represents a null literal
type NullLiteral struct {
	BaseExpression
}

// TokenLiteral returns the literal value of the token
func (n *NullLiteral) TokenLiteral() string {
	return "null"
}

// NewNullLiteral creates a new null literal
func NewNullLiteral(pos Position) *NullLiteral {
	return &NullLiteral{
		BaseExpression: BaseExpression{Pos: pos},
	}
}

type BooleanLiteral struct {
	BaseExpression
	Value bool
}

// TokenLiteral returns the literal value of the token
func (b *BooleanLiteral) TokenLiteral() string {
	return strconv.FormatBool(b.Value)
}

// NewBooleanLiteral creates a new boolean literal
func NewBooleanLiteral(value bool, pos Position) *BooleanLiteral {
	return &BooleanLiteral{
		BaseExpression: BaseExpression{Pos: pos},
		Value:          value,
	}
}

// ArrayLiteral represents an array literal
type ArrayLiteral struct {
	BaseExpression
	Elements []Expression
}

// TokenLiteral returns the literal value of the token
func (a *ArrayLiteral) TokenLiteral() string {
	return "array"
}

// NewArrayLiteral creates a new array literal
func NewArrayLiteral(pos Position) *ArrayLiteral {
	return &ArrayLiteral{
		BaseExpression: BaseExpression{Pos: pos},
		Elements:       make([]Expression, 0),
	}
}

// AddElement adds an element to the array
func (a *ArrayLiteral) AddElement(element Expression) {
	a.Elements = append(a.Elements, element)
}

// DictionaryLiteral represents a dictionary literal
type DictionaryLiteral struct {
	BaseExpression
	Pairs map[Expression]Expression
}

// TokenLiteral returns the literal value of the token
func (d *DictionaryLiteral) TokenLiteral() string {
	return "dictionary"
}

// NewDictionaryLiteral creates a new dictionary literal
func NewDictionaryLiteral(pos Position) *DictionaryLiteral {
	return &DictionaryLiteral{
		BaseExpression: BaseExpression{Pos: pos},
		Pairs:          make(map[Expression]Expression),
	}
}

// AddPair adds a key-value pair to the dictionary
func (d *DictionaryLiteral) AddPair(key, value Expression) {
	d.Pairs[key] = value
}

// PrefixExpression represents a prefix operator expression
type PrefixExpression struct {
	BaseExpression
	Operator string
	Right    Expression
}

// TokenLiteral returns the literal value of the token
func (p *PrefixExpression) TokenLiteral() string {
	return p.Operator
}

// NewPrefixExpression creates a new prefix expression
func NewPrefixExpression(operator string, right Expression, pos Position) *PrefixExpression {
	return &PrefixExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Operator:       operator,
		Right:          right,
	}
}

// InfixExpression represents an infix operator expression
type InfixExpression struct {
	BaseExpression
	Left     Expression
	Operator string
	Right    Expression
}

// TokenLiteral returns the literal value of the token
func (i *InfixExpression) TokenLiteral() string {
	return i.Operator
}

// NewInfixExpression creates a new infix expression
func NewInfixExpression(left Expression, operator string, right Expression, pos Position) *InfixExpression {
	return &InfixExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Left:           left,
		Operator:       operator,
		Right:          right,
	}
}

// CallExpression represents a function call
type CallExpression struct {
	BaseExpression
	Function  Expression
	Arguments []Expression
}

// TokenLiteral returns the literal value of the token
func (c *CallExpression) TokenLiteral() string {
	return "call"
}

// NewCallExpression creates a new call expression
func NewCallExpression(function Expression, pos Position) *CallExpression {
	return &CallExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Function:       function,
		Arguments:      make([]Expression, 0),
	}
}

// AddArgument adds an argument to the call
func (c *CallExpression) AddArgument(arg Expression) {
	c.Arguments = append(c.Arguments, arg)
}

// IndexExpression represents an array/dictionary index expression
type IndexExpression struct {
	BaseExpression
	Left  Expression
	Index Expression
}

// TokenLiteral returns the literal value of the token
func (i *IndexExpression) TokenLiteral() string {
	return "index"
}

// NewIndexExpression creates a new index expression
func NewIndexExpression(left, index Expression, pos Position) *IndexExpression {
	return &IndexExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Left:           left,
		Index:          index,
	}
}

// DotExpression represents a dot operator expression (obj.property)
type DotExpression struct {
	BaseExpression
	Left     Expression
	Property string
}

// TokenLiteral returns the literal value of the token
func (d *DotExpression) TokenLiteral() string {
	return "."
}

// NewDotExpression creates a new dot expression
func NewDotExpression(left Expression, property string, pos Position) *DotExpression {
	return &DotExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Left:           left,
		Property:       property,
	}
}

// AssignmentExpression represents an assignment expression
type AssignmentExpression struct {
	BaseExpression
	Left     Expression
	Operator string
	Right    Expression
}

// TokenLiteral returns the literal value of the token
func (a *AssignmentExpression) TokenLiteral() string {
	return a.Operator
}

// NewAssignmentExpression creates a new assignment expression
func NewAssignmentExpression(left Expression, operator string, right Expression, pos Position) *AssignmentExpression {
	return &AssignmentExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Left:           left,
		Operator:       operator,
		Right:          right,
	}
}

// ConditionalExpression represents a conditional expression (ternary operator)
// Format: value_if_true if condition else value_if_false
type ConditionalExpression struct {
	BaseExpression
	Condition    Expression
	ValueIfTrue  Expression
	ValueIfFalse Expression
}

// TokenLiteral returns the literal value of the token
func (c *ConditionalExpression) TokenLiteral() string {
	return "conditional"
}

// NewConditionalExpression creates a new conditional expression
func NewConditionalExpression(condition, valueIfTrue, valueIfFalse Expression, pos Position) *ConditionalExpression {
	return &ConditionalExpression{
		BaseExpression: BaseExpression{Pos: pos},
		Condition:      condition,
		ValueIfTrue:    valueIfTrue,
		ValueIfFalse:   valueIfFalse,
	}
}
