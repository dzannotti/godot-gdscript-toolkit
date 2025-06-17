package ast

// Annotation represents a GDScript annotation (@annotation)
type Annotation struct {
	Pos  Position
	Name string
	Args []Expression
}

// Position returns the position of the annotation in the source code
func (a *Annotation) Position() Position {
	return a.Pos
}

// TokenLiteral returns the literal value of the token
func (a *Annotation) TokenLiteral() string {
	return a.Name
}

// NewAnnotation creates a new annotation
func NewAnnotation(name string, pos Position) *Annotation {
	return &Annotation{
		Pos:  pos,
		Name: name,
		Args: make([]Expression, 0),
	}
}

// AddArg adds an argument to the annotation
func (a *Annotation) AddArg(arg Expression) {
	a.Args = append(a.Args, arg)
}
