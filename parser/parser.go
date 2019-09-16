package parser

import (
	"fmt"

	"github.com/ollybritton/monkey/ast"
	"github.com/ollybritton/monkey/lexer"
	"github.com/ollybritton/monkey/token"
)

// Parser represents a parser that will be used to parse the program.
type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

// New returns a new parser from a lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens so that curToken and peekTOken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// Errors returns the errors encountered while parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError creates a new error that says that the peeked token was expected to be something else.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken gets the next token from the lexer.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses the program into an abstract syntax tree, the root node being an *ast.Program struct.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

// parseStatement parses a single statement into an ast.Statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatment parses a let statement into an ast.LetStatement.
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// Check that the next token is an identifier, and move it to p.curToken if it is.
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Value: p.curToken.Literal, Token: p.curToken}

	// Check that the next token is an assignment (=), and move to p.curToken if it is.
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we encounter a semicolon.
	// We haven't written the expression parsing code yet!
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parsers a return statement into an ast.ReturnStatment.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: We're skipping the expressions until we encounter a semicolon.
	// We haven't written the expression parsing code yet!
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// curTokenIs checks if the current token is a specific type of token.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs chekcs if the next token is a specific type of token.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the peaked token is a specific type of token. If it is, it will read another token and return true.
// otherwise, it returns false.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}
