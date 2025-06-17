package ast

// IfStatement represents an if statement
type IfStatement struct {
	BaseStatement
	Condition     Expression
	Consequence   []Statement
	ElseCondition []Expression
	ElseBranches  [][]Statement
	Alternative   []Statement
}

// NewIfStatement creates a new if statement
func NewIfStatement(pos Position, condition Expression) *IfStatement {
	return &IfStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "if_stmt",
			Annotations: make([]*Annotation, 0),
		},
		Condition:     condition,
		Consequence:   make([]Statement, 0),
		ElseCondition: make([]Expression, 0),
		ElseBranches:  make([][]Statement, 0),
		Alternative:   make([]Statement, 0),
	}
}

// AddConsequenceStatement adds a statement to the if branch
func (i *IfStatement) AddConsequenceStatement(statement Statement) {
	i.Consequence = append(i.Consequence, statement)
}

// AddElseIfBranch adds an else-if branch
func (i *IfStatement) AddElseIfBranch(condition Expression, statements []Statement) {
	i.ElseCondition = append(i.ElseCondition, condition)
	i.ElseBranches = append(i.ElseBranches, statements)
}

// SetAlternative sets the else branch
func (i *IfStatement) SetAlternative(statements []Statement) {
	i.Alternative = statements
}

// ForStatement represents a for loop
type ForStatement struct {
	BaseStatement
	Iterator    string
	IteratorPos Position
	Collection  Expression
	Body        []Statement
	TypeHint    string
	IsTyped     bool
}

// NewForStatement creates a new for statement
func NewForStatement(pos Position, iterator string, iteratorPos Position, collection Expression, isTyped bool) *ForStatement {
	kind := "for_stmt"
	if isTyped {
		kind = "for_stmt_typed"
	}

	return &ForStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        kind,
			Annotations: make([]*Annotation, 0),
		},
		Iterator:    iterator,
		IteratorPos: iteratorPos,
		Collection:  collection,
		Body:        make([]Statement, 0),
		IsTyped:     isTyped,
	}
}

// AddBodyStatement adds a statement to the for loop body
func (f *ForStatement) AddBodyStatement(statement Statement) {
	f.Body = append(f.Body, statement)
}

// SetTypeHint sets the type hint for the iterator
func (f *ForStatement) SetTypeHint(typeHint string) {
	f.TypeHint = typeHint
	f.IsTyped = true
	f.Kind = "for_stmt_typed"
}

// WhileStatement represents a while loop
type WhileStatement struct {
	BaseStatement
	Condition Expression
	Body      []Statement
}

// NewWhileStatement creates a new while statement
func NewWhileStatement(pos Position, condition Expression) *WhileStatement {
	return &WhileStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "while_stmt",
			Annotations: make([]*Annotation, 0),
		},
		Condition: condition,
		Body:      make([]Statement, 0),
	}
}

// AddBodyStatement adds a statement to the while loop body
func (w *WhileStatement) AddBodyStatement(statement Statement) {
	w.Body = append(w.Body, statement)
}

// MatchStatement represents a match statement
type MatchStatement struct {
	BaseStatement
	Value    Expression
	Branches []*MatchBranch
}

// NewMatchStatement creates a new match statement
func NewMatchStatement(pos Position, value Expression) *MatchStatement {
	return &MatchStatement{
		BaseStatement: BaseStatement{
			Pos:         pos,
			Kind:        "match_stmt",
			Annotations: make([]*Annotation, 0),
		},
		Value:    value,
		Branches: make([]*MatchBranch, 0),
	}
}

// AddBranch adds a branch to the match statement
func (m *MatchStatement) AddBranch(branch *MatchBranch) {
	m.Branches = append(m.Branches, branch)
}

// MatchBranch represents a branch in a match statement
type MatchBranch struct {
	Pos       Position
	Pattern   Expression
	Guard     Expression
	Body      []Statement
	IsGuarded bool
}

// NewMatchBranch creates a new match branch
func NewMatchBranch(pos Position, pattern Expression) *MatchBranch {
	return &MatchBranch{
		Pos:     pos,
		Pattern: pattern,
		Body:    make([]Statement, 0),
	}
}

// SetGuard sets the guard condition for the branch
func (m *MatchBranch) SetGuard(guard Expression) {
	m.Guard = guard
	m.IsGuarded = true
}

// AddBodyStatement adds a statement to the branch body
func (m *MatchBranch) AddBodyStatement(statement Statement) {
	m.Body = append(m.Body, statement)
}

// Position returns the position of the branch in the source code
func (m *MatchBranch) Position() Position {
	return m.Pos
}
