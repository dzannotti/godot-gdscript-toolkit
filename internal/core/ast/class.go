package ast

// Class represents a GDScript class
type Class struct {
	Pos         Position
	Name        string
	Extends     string
	SubClasses  []*Class
	Functions   []*Function
	Statements  []Statement
	Annotations []*Annotation
}

// Position returns the position of the class in the source code
func (c *Class) Position() Position {
	return c.Pos
}

// TokenLiteral returns the literal value of the token
func (c *Class) TokenLiteral() string {
	return c.Name
}

func (c *Class) statementNode() {}

// NewClass creates a new class
func NewClass(name string, pos Position) *Class {
	return &Class{
		Pos:         pos,
		Name:        name,
		SubClasses:  make([]*Class, 0),
		Functions:   make([]*Function, 0),
		Statements:  make([]Statement, 0),
		Annotations: make([]*Annotation, 0),
	}
}

// AddSubClass adds a subclass to the class
func (c *Class) AddSubClass(class *Class) {
	c.SubClasses = append(c.SubClasses, class)
}

// AddFunction adds a function to the class
func (c *Class) AddFunction(function *Function) {
	c.Functions = append(c.Functions, function)
}

// AddStatement adds a statement to the class
func (c *Class) AddStatement(statement Statement) {
	c.Statements = append(c.Statements, statement)
}

// AddAnnotation adds an annotation to the class
func (c *Class) AddAnnotation(annotation *Annotation) {
	c.Annotations = append(c.Annotations, annotation)
}
