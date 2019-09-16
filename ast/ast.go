package ast

import "github.com/ollybritton/monkey/token"

// Node is a node in the abstract syntax tree.
// TokenLiteral should return a string representing the literal value of the token associated with that Node,
// such as "5" or "if"
type Node interface {
	TokenLiteral() string
}

// Statement is a node which holds a statement. A statement is that it doesn't produce a value, unlike expressions, which do.
// For example, "let i = 0" is a statement, whereas "i+1" is an expresison.
type Statement interface {
	Node
	statementNode()
}

// Expression is a node which represents an expression. An expression produces a value, unlike statements, which do not.
// For example, "(5/8)+4" is an expression, but "if (true) { a = a+1 }" is not.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node for every AST that the parser produces. Every valid Monkey program is a series of statements.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the token literal of the first statement in the tree.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}

	return p.Statements[0].TokenLiteral()
}

// LetStatement represents a let statement, such as "let a = 1" or "let q = 5 * add(1,2)"
// The general form is "let <ident> = <expression>"
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral returns the literal value of the LET token. This is will always be "let"
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier represents an identifier, which is something like a variable or function name.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the literal value of the IDENT token. This will be the name of the variable.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// ReturnStatement represents a return statement, such as "return 0" or "return add(15)"
// The general form is "return <expression>"
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal value of the return token. This will always be "return".
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
