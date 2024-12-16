package ast

import (
	"main/token"
	"testing"
)

func TestString(t *testing.T) {
	expected := "let x = 10;"

	letToken, _ := token.MapSourceToKeyword("let")
	p := Program{Statements: []Statement{
		LetStatement{
			Token:      letToken,
			Identifier: IdentifierExpression{token.Token{Type: token.LET, Literal: "x"}},
			Expression: IntExpression{Token: token.Token{Type: token.INT, Literal: "10"}},
		},
	}}

	if p.String() != expected {
		t.Fatalf("error - expected: %s - actual: %s", expected, p.String())
	}
}
