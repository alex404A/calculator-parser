package pratt

import (
	"calculator/lexer"
	"calculator/token"
	"fmt"
	"math"
	"strconv"
)

const (
	_        int = iota
	LOWEST       // lowest precedence
	SUM          // +
	PRODUCT      // *
	EXPONENT     // **
	PREFIX       // -X
	CALL         // LEFTP
)

var precedences = map[token.TokenType]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.EXPONENT: EXPONENT,
	token.LEFTP:    CALL,
}

type (
	prefixParseFn func() float64
	infixParseFn  func(float64) float64
)

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.LITERAL, p.parseLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.PLUS, p.parsePrefix)
	p.registerPrefix(token.LEFTP, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfix)
	p.registerInfix(token.MINUS, p.parseInfix)
	p.registerInfix(token.ASTERISK, p.parseInfix)
	p.registerInfix(token.SLASH, p.parseInfix)
	p.registerInfix(token.EXPONENT, p.parseInfix)
	p.NextToken()
	p.NextToken()
	return p
}

func (parser *Parser) NextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.l.NewToken()
}

func (parser *Parser) calc() float64 {
	return parser.parseExpression(LOWEST)
}

func (parser *Parser) parseExpression(precedence int) float64 {
	prefix, ok := parser.prefixParseFns[parser.curToken.Type]
	if !ok {
		panic(fmt.Sprintf("no prefix function for token type %s", parser.curToken.Type))
	}
	left := prefix()
	for parser.curToken.Type != token.EOF && precedence < precedences[parser.curToken.Type] {
		infix, ok := parser.infixParseFns[parser.curToken.Type]
		if !ok {
			panic(fmt.Sprintf("no infix function for token type %s", parser.curToken.Type))
		}
		left = infix(left)
	}
	return left
}

func (parser *Parser) parseLiteral() float64 {
	val, _ := strconv.ParseFloat(parser.curToken.Literal, 64)
	parser.NextToken()
	return val
}

func (parser *Parser) parsePrefix() float64 {
	t := parser.curToken
	precedence := precedences[parser.curToken.Type]
	parser.NextToken()
	val := parser.parseExpression(precedence)
	if t.Type == token.MINUS {
		return 0 - val
	} else {
		return val
	}
}

func (parser *Parser) parseGroupedExpression() float64 {
	parser.NextToken()

	val := parser.parseExpression(LOWEST)

	if parser.curToken.Type != token.RIGHTP {
		panic("no right parentheses found")
	}

	parser.NextToken()

	return val
}

func (parser *Parser) parseInfix(left float64) float64 {
	t := parser.curToken
	precedence, ok := precedences[t.Type]
	if !ok {
		panic(fmt.Sprintf("no precedence for token type %s", t.Type))
	}
	parser.NextToken()
	right := parser.parseExpression(precedence)
	if t.Type == token.PLUS {
		return left + right
	} else if t.Type == token.MINUS {
		return left - right
	} else if t.Type == token.ASTERISK {
		return left * right
	} else if t.Type == token.SLASH {
		return left / right
	} else if t.Type == token.EXPONENT {
		return math.Pow(left, right)
	} else {
		panic(fmt.Sprintf("unknown token type %s in infix", t.Type))
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
