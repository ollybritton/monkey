package lexer

import "github.com/ollybritton/monkey/token"

// Lexer represents a lexer for a monkey program.
// It acts on an ASCII string, not a unicode one for simplicity. If we wanted to use Unicode, we'd have to change l.ch
// to 'rune', and update the logic for l.input[l.readPosition] as we wouldn't be able to access single bytes now.
type Lexer struct {
	input        string
	position     int  // current position in input (index of current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a new lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}

	// Read one character so the lexer is fully initialised with values when returned.
	l.readChar()

	return l
}

// isLetter returns true if the character specified is a letter, and false if it is not (kind of self-explanatory if you ask me)
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns true if the character is a number.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isWhitespace returns true if the character is a type of whitespace (a space, a tab, a newline or a linefeed)
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// skipWhitespace consumes whitespace characters.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// readChar reads us the next character in the input. If there is no input left to read (i.e. the input is finished or the
// input is blank) then set the char value to ASCII NUL.
func (l *Lexer) readChar() {
	// Sets the current char under examination to the null char if there are no more chars left to read.
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	// Basically saying "set the current character to the next character"
	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns the next char in the input as a byte.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// readIdentifier reads a set of characters and returns the characters that it read as a string.
func (l *Lexer) readIdentifier() string {
	startPosition := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

// readNumber reads a set digits and returns the string representation of that number.
// At the moment, only integers are supported.
func (l *Lexer) readNumber() string {
	startPosition := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[startPosition:l.position]
}

// readString reads a string of characters.
func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

// newToken returns a new token from a specified token type and literal value, given as a byte.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// NextToken returns the next token in the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}

	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}

		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}
