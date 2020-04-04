package parser

import (
	"calculator/lexer"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue float64
	}{
		// {"3", 3.0},
		// {"3.0", 3.0},
		// {"+3", 3.0},
		{"-3", -3.0},
		{"3 + 2.0", 5.0},
		{"3 - 2.0", 1.0},
		{"3 - 2.0 * 2", -1.0},
		{"3 + 2.0 / 2", 4.0},
		{"3 + 2.0 / (2 - 1)", 5.0},
		{"3 + 4.0 / (2 - (1 - 3)) + 5.3", 9.3},
		{"3 + 2**3**(1+0 * 3) + 22", 33},
		{"3 + 2**-1", 3.5},
	}

	for _, test := range tests {
		l := lexer.NewLexer(test.input)
		p := NewParser(l)
		val := p.calculate()
		if val != test.expectedValue {
			t.Fatalf("input %s, expected %f, actual %f", test.input, test.expectedValue, val)
		}
	}
}
