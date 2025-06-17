package ast

// BaseStatement provides common functionality for all statements
type BaseStatement struct {
	Pos         Position
	Kind        string
	Annotations []*Annotation
}

// Position returns the position of the statement in the source code
func (s *BaseStatement) Position() Position {
	return s.Pos
}

// TokenLiteral returns the literal value of the token
func (s *BaseStatement) TokenLiteral() string {
	return s.Kind
}

func (s *BaseStatement) statementNode() {}

// PassStatement represents a 'pass' statement
type PassStatement struct {
	BaseStatement
}

// NewPassStatement creates a new pass statement
func NewPassStatement(pos Position) *PassStatement {
	return &PassStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "pass_stmt",
			Annotations: make([]*Annotation, 0),
		},
	}
}

// ReturnStatement represents a 'return' statement
type ReturnStatement struct {
	BaseStatement
	Value Expression
}

// NewReturnStatement creates a new return statement
func NewReturnStatement(pos Position, value Expression) *ReturnStatement {
	return &ReturnStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "return_stmt",
			Annotations: make([]*Annotation, 0),
		},
		Value: value,
	}
}

// BreakStatement represents a 'break' statement
type BreakStatement struct {
	BaseStatement
}

// NewBreakStatement creates a new break statement
func NewBreakStatement(pos Position) *BreakStatement {
	return &BreakStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "break_stmt",
			Annotations: make([]*Annotation, 0),
		},
	}
}

// ContinueStatement represents a 'continue' statement
type ContinueStatement struct {
	BaseStatement
}

// NewContinueStatement creates a new continue statement
func NewContinueStatement(pos Position) *ContinueStatement {
	return &ContinueStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "continue_stmt",
			Annotations: make([]*Annotation, 0),
		},
	}
}

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	BaseStatement
	Expression Expression
}

// NewExpressionStatement creates a new expression statement
func NewExpressionStatement(pos Position, expr Expression) *ExpressionStatement {
	return &ExpressionStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "expr_stmt",
			Annotations: make([]*Annotation, 0),
		},
		Expression: expr,
	}
}

// VarStatement represents a variable declaration statement
type VarStatement struct {
	BaseStatement
	Name       string
	TypeHint   string
	Value      Expression
	IsTyped    bool
	IsConst    bool
	IsInferred bool
}

// NewVarStatement creates a new variable declaration statement
func NewVarStatement(pos Position, name string, isTyped bool) *VarStatement {
	kind := "func_var_stmt"
	if isTyped {
		kind = "func_var_typed"
	}

	return &VarStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        kind,
			Annotations: make([]*Annotation, 0),
		},
		Name:    name,
		IsTyped: isTyped,
	}
}

// SetValue sets the value of the variable
func (v *VarStatement) SetValue(value Expression) {
	v.Value = value
	if v.IsConst {
		if v.IsTyped {
			v.Kind = "const_typed_assgnd"
		} else {
			v.Kind = "const_assigned"
		}
	} else {
		if v.IsTyped {
			v.Kind = "func_var_typed_assgnd"
		} else {
			v.Kind = "func_var_assigned"
		}
	}
}

// SetTypeHint sets the type hint of the variable
func (v *VarStatement) SetTypeHint(typeHint string) {
	v.TypeHint = typeHint
	v.IsTyped = true
	if v.Value != nil {
		if v.IsConst {
			v.Kind = "const_typed_assgnd"
		} else {
			v.Kind = "func_var_typed_assgnd"
		}
	} else {
		if v.IsConst {
			v.Kind = "const_typed"
		} else {
			v.Kind = "func_var_typed"
		}
	}
}

// SetConst marks the variable as a constant
func (v *VarStatement) SetConst() {
	v.IsConst = true
	if v.IsTyped {
		if v.Value != nil {
			v.Kind = "const_typed_assgnd"
		} else {
			v.Kind = "const_typed"
		}
	} else {
		if v.Value != nil {
			v.Kind = "const_assigned"
		} else {
			v.Kind = "const_stmt"
		}
	}
}

// SetInferred marks the variable as type-inferred (:=)
func (v *VarStatement) SetInferred() {
	v.IsInferred = true
	v.IsTyped = true
	if v.IsConst {
		v.Kind = "const_typed_assgnd"
	} else {
		v.Kind = "func_var_assigned"
	}
}
