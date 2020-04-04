package token

type TokenType string

const (
	LITERAL  = "LITERAL"
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	EXPONENT = "**"
	SLASH    = "/"
	LEFTP    = "("
	RIGHTP   = ")"
	EOF      = "EOF"
	ILLEGAL  = "ILLEGAL"
)

type Token struct {
	Type    TokenType
	Literal string
}
