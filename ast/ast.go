package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ollybritton/monkey/token"
)

// Node is a node in the abstract syntax tree.
// TokenLiteral should return a string representing the literal value of the token associated with that Node,
// such as "5" or "if"
type Node interface {
	TokenLiteral() string
	String() string
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

// String returns all the statements in the program joined together.
func (p *Program) String() string {
	out := bytes.Buffer{}

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

// String returns the string representation of that let statement.
func (ls *LetStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	} else {
		out.WriteString("<nil>")
	}

	out.WriteString(";")

	return out.String()
}

// Identifier represents an identifier, which is something like a variable or function name.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral returns the literal value of the IDENT token. This will be the name of the variable.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns the name of the identifier.
func (i *Identifier) String() string { return i.Value }

// ReturnStatement represents a return statement, such as "return 0" or "return add(15)"
// The general form is "return <expression>"
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the literal value of the return token. This will always be "return".
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns the string representation of the return statement.
func (rs *ReturnStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	} else {
		out.WriteString("<nil>")
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement is an expresssion that is a statement. It does nothing, but is totally legal monkey code.
// For example, "x+1;" is valid but has no effect. The general form is "<expression>;"
type ExpressionStatement struct {
	Token      token.Token // The first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the literal value of the first token in the expression.
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns the string representation of the expression statement.
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return "<nil>"
}

// BlockStatement represents a set of statements surrounded by braces ("{" & "}").
type BlockStatement struct {
	Token      token.Token // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral is the token's literal. I'm sick of writing this.
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String is the string representation of the block.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// IntegerLiteral represents an integer in the AST, like "5".
type IntegerLiteral struct {
	Token token.Token // the token.INT type
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral returns the literal value of the first token in the expression, which in this case is the value of the number
// as a string.
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns the string representation of the integer.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// Boolean represents a boolean in the ast, either "true" or "false"
type Boolean struct {
	Token token.Token // the token.TRUE | token.FALSE
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral the literal value of the token.
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns the string representation of the boolean.
func (b *Boolean) String() string { return b.Token.Literal }

// PrefixExpression wraps an expression using a prefix, such as "-" or "!"
type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. - or !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the string representation of the prefix token, such as "-".
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns the string representation of the prefix expression.
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// InfixExpression represents an infix expression, such as "a+5"
type InfixExpression struct {
	Token    token.Token // The operator, such as +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral returns the operator value as a string.
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns the infix expression as a string, wrapping in brackets.
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// IfExpression represents an if-else statement in the AST.
type IfExpression struct {
	Token       token.Token // the 'if' token.
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ife *IfExpression) expressionNode() {}

// TokenLiteral returns the literal value of the 'if' token, which is always 'if'.
func (ife *IfExpression) TokenLiteral() string { return ife.Token.Literal }

// String returns the if statement represented as a string.
func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Condition.String())
	out.WriteString(" ")
	out.WriteString(ife.Consequence.String())

	if ife.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ife.Alternative.String())
	}

	return out.String()
}

// FunctionLiteral represents a function in the AST.
type FunctionLiteral struct {
	Token      token.Token // the 'fn' token.
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral is the token, literal.
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns the string representation of the function.
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents a function call inside the program.
// <expression>(<command seperated expressions>)
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral returns the literal value expression's token.
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns the string representation of the function.
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var arguments []string
	for _, a := range ce.Arguments {
		arguments = append(arguments, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(arguments, ", "))
	out.WriteString(")")

	return out.String()
}
