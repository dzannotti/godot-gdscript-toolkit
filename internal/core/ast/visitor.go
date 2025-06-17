package ast

// Visitor is the interface for the visitor pattern
type Visitor interface {
	// Visit is called for each node in the AST
	// If the returned visitor is not nil, it will be used to visit the children of the node
	// If the returned visitor is nil, the children of the node will not be visited
	Visit(node Node) Visitor
}

// Walk traverses the AST starting from the given node and calls the visitor for each node
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *AbstractSyntaxTree:
		if n.RootClass != nil {
			Walk(v, n.RootClass)
		}
		for _, class := range n.Classes {
			if class != n.RootClass {
				Walk(v, class)
			}
		}
		for _, function := range n.Functions {
			Walk(v, function)
		}

	case *Class:
		for _, annotation := range n.Annotations {
			Walk(v, annotation)
		}
		for _, statement := range n.Statements {
			Walk(v, statement)
		}
		for _, function := range n.Functions {
			Walk(v, function)
		}
		for _, subClass := range n.SubClasses {
			Walk(v, subClass)
		}

	case *Function:
		for _, annotation := range n.Annotations {
			Walk(v, annotation)
		}
		for _, param := range n.Parameters {
			Walk(v, param)
		}
		for _, statement := range n.Statements {
			Walk(v, statement)
		}

	case *Parameter:
		if n.Default != nil {
			Walk(v, n.Default)
		}

	case *Annotation:
		for _, arg := range n.Args {
			Walk(v, arg)
		}

	case *PassStatement, *BreakStatement, *ContinueStatement:
		// These statements have no children

	case *ReturnStatement:
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *ExpressionStatement:
		Walk(v, n.Expression)

	case *VarStatement:
		if n.Value != nil {
			Walk(v, n.Value)
		}

	case *IfStatement:
		Walk(v, n.Condition)
		for _, stmt := range n.Consequence {
			Walk(v, stmt)
		}
		for i, cond := range n.ElseCondition {
			Walk(v, cond)
			for _, stmt := range n.ElseBranches[i] {
				Walk(v, stmt)
			}
		}
		for _, stmt := range n.Alternative {
			Walk(v, stmt)
		}

	case *ForStatement:
		Walk(v, n.Collection)
		for _, stmt := range n.Body {
			Walk(v, stmt)
		}

	case *WhileStatement:
		Walk(v, n.Condition)
		for _, stmt := range n.Body {
			Walk(v, stmt)
		}

	case *MatchStatement:
		Walk(v, n.Value)
		for _, branch := range n.Branches {
			Walk(v, branch.Pattern)
			if branch.Guard != nil {
				Walk(v, branch.Guard)
			}
			for _, stmt := range branch.Body {
				Walk(v, stmt)
			}
		}

	case *Identifier, *StringLiteral, *NumberLiteral, *BooleanLiteral:
		// These expressions have no children

	case *ArrayLiteral:
		for _, element := range n.Elements {
			Walk(v, element)
		}

	case *DictionaryLiteral:
		for key, value := range n.Pairs {
			Walk(v, key)
			Walk(v, value)
		}

	case *PrefixExpression:
		Walk(v, n.Right)

	case *InfixExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *CallExpression:
		Walk(v, n.Function)
		for _, arg := range n.Arguments {
			Walk(v, arg)
		}

	case *IndexExpression:
		Walk(v, n.Left)
		Walk(v, n.Index)

	case *DotExpression:
		Walk(v, n.Left)

	case *AssignmentExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)
	}
}

// Inspector is a function that can be used to inspect nodes during traversal
type Inspector func(Node) bool

// Inspect traverses the AST starting from the given node and calls the inspector for each node
// If the inspector returns false, the children of the node will not be visited
func Inspect(node Node, f Inspector) {
	Walk(inspectorVisitor(f), node)
}

type inspectorVisitor func(Node) bool

func (v inspectorVisitor) Visit(node Node) Visitor {
	if v(node) {
		return v
	}
	return nil
}
