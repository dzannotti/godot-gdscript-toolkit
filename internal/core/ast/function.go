package ast

// Parameter represents a function parameter
type Parameter struct {
	Pos      Position
	Name     string
	TypeHint string
	Default  Expression
}

// Position returns the position of the parameter in the source code
func (p *Parameter) Position() Position {
	return p.Pos
}

// TokenLiteral returns the literal value of the token
func (p *Parameter) TokenLiteral() string {
	return p.Name
}

// Function represents a GDScript function
type Function struct {
	Pos           Position
	Name          string
	Parameters    []*Parameter
	ReturnType    string
	Statements    []Statement
	SubStatements []Statement
	Annotations   []*Annotation
	IsStatic      bool
}

// Position returns the position of the function in the source code
func (f *Function) Position() Position {
	return f.Pos
}

// TokenLiteral returns the literal value of the token
func (f *Function) TokenLiteral() string {
	return f.Name
}

func (f *Function) statementNode() {}

// NewFunction creates a new function
func NewFunction(name string, pos Position) *Function {
	return &Function{
		Pos:           pos,
		Name:          name,
		Parameters:    make([]*Parameter, 0),
		Statements:    make([]Statement, 0),
		SubStatements: make([]Statement, 0),
		Annotations:   make([]*Annotation, 0),
	}
}

// AddParameter adds a parameter to the function
func (f *Function) AddParameter(param *Parameter) {
	f.Parameters = append(f.Parameters, param)
}

// AddStatement adds a statement to the function
func (f *Function) AddStatement(statement Statement) {
	f.Statements = append(f.Statements, statement)
	f.SubStatements = append(f.SubStatements, statement)
}

// AddSubStatement adds a sub-statement to the function
func (f *Function) AddSubStatement(statement Statement) {
	f.SubStatements = append(f.SubStatements, statement)
}

// AddAnnotation adds an annotation to the function
func (f *Function) AddAnnotation(annotation *Annotation) {
	f.Annotations = append(f.Annotations, annotation)
}
