package lexer

import (
	"calculator/token"
	"testing"
)

func TestTokens(t *testing.T) {
	input := "3 + 2 / (1 -0) + 4.0 * 5**2"
	tests := []struct {
		expectedType  token.TokenType
		expectedValue string
	}{
		{token.LITERAL, "3"},
		{token.PLUS, "+"},
		{token.LITERAL, "2"},
		{token.SLASH, "/"},
		{token.LEFTP, "("},
		{token.LITERAL, "1"},
		{token.MINUS, "-"},
		{token.LITERAL, "0"},
		{token.RIGHTP, ")"},
		{token.PLUS, "+"},
		{token.LITERAL, "4.0"},
		{token.ASTERISK, "*"},
		{token.LITERAL, "5"},
		{token.EXPONENT, "**"},
		{token.LITERAL, "2"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for _, test := range tests {
		tok := l.NewToken()
		if tok.Type != test.expectedType {
			t.Fatalf("token type wrong, expected %s, actual %s", test.expectedType, tok.Type)
		}
		if tok.Literal != test.expectedValue {
			t.Fatalf("token value wrong, expected %s, actual %s", test.expectedValue, tok.Literal)
		}
	}
}
