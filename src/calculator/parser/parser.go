package parser

import (
	"calculator/lexer"
	"calculator/token"
	"fmt"
	"math"
	"strconv"
)

/*
  simple calculator implemented via recursive descent
  add_op := + | -
  mul_op := * | /
  expr   := term {add_op term}
  term   := factor {mul_op factor}
  factor := element | '(' expr ')'
  elemetn := digits { exp }
	digits := {+|-} [0..9] {[0..9]}
	exp := ** digits { ** digits }
*/

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	parser := &Parser{l: l}
	parser.NextToken()
	parser.NextToken()
	return parser
}

func (parser *Parser) NextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.l.NewToken()
}

func (parser *Parser) calculate() float64 {
	return parser.findExpr()
}

func (parser *Parser) findExpr() float64 {
	val := parser.findTerm()
	for parser.curToken.Type == token.PLUS || parser.curToken.Type == token.MINUS {
		op := parser.curToken.Type
		parser.NextToken()
		term := parser.findTerm()
		if op == token.PLUS {
			val += term
		} else {
			val -= term
		}
	}
	if parser.curToken.Type != token.EOF && parser.curToken.Type != token.RIGHTP {
		panic(fmt.Sprintf("the last token type is %s, token value is %s", parser.curToken.Type, parser.curToken.Literal))
	} else {
		if parser.curToken.Type == token.RIGHTP {
			parser.NextToken()
		}
		return val
	}
}

func (parser *Parser) findTerm() float64 {
	val := parser.findFactor()
	for parser.curToken.Type == token.ASTERISK || parser.curToken.Type == token.SLASH {
		op := parser.curToken.Type
		parser.NextToken()
		factor := parser.findFactor()
		if op == token.ASTERISK {
			val *= factor
		} else {
			val /= factor
		}
	}
	return val
}

func (parser *Parser) findFactor() float64 {
	if parser.curToken.Type == token.LEFTP {
		parser.NextToken()
		return parser.findExpr()
	} else {
		return parser.findElement()
	}
}

func (parser *Parser) findElement() float64 {
	base := parser.findDigit()
	if parser.curToken.Type == token.EXPONENT {
		exp := parser.findExponent()
		return math.Pow(base, exp)
	} else {
		return base
	}
}

func (parser *Parser) findDigit() float64 {
	var result float64
	if parser.curToken.Type == token.PLUS || parser.curToken.Type == token.MINUS {
		if parser.peekToken.Type != token.LITERAL {
			panic(fmt.Sprintf("fail to find literal after op in factor, token type is %s, token value is %s", parser.peekToken.Type, parser.peekToken.Literal))
		}
		op := parser.curToken.Type
		parser.NextToken()
		val, _ := strconv.ParseFloat(parser.curToken.Literal, 64)
		if op == token.PLUS {
			result = val
		} else {
			result = 0 - val
		}
	} else {
		result, _ = strconv.ParseFloat(parser.curToken.Literal, 64)
	}
	parser.NextToken()
	return result
}

func (parser *Parser) findExponent() float64 {
	arr := []float64{}
	var val float64
	for {
		if parser.curToken.Type != token.EXPONENT {
			break
		}
		parser.NextToken()
		if parser.curToken.Type == token.LEFTP {
			parser.NextToken()
			val = parser.findExpr()
		} else {
			val = parser.findDigit()
		}
		arr = append(arr, val)
	}
	cur := 1.0
	for i := len(arr) - 1; i >= 0; i-- {
		cur = math.Pow(arr[i], cur)
	}
	return cur
}
