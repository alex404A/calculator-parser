package lexer

import (
	"calculator/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func (l *Lexer) NewToken() token.Token {
	l.skipWhitespace()
	var tok token.Token
	switch l.ch {
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		if l.peekNext() == '*' {
			tok.Type = token.EXPONENT
			tok.Literal = "**"
			l.readChar()
		} else {
			tok = newToken(token.ASTERISK, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '(':
		tok = newToken(token.LEFTP, l.ch)
	case ')':
		tok = newToken(token.RIGHTP, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.LITERAL
			return tok
		} else {
			// token = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekNext() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || ch == '.'
}
