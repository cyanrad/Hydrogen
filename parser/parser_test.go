package parser

import (
	"main/ast"
	"main/lexer"
	"main/token"
	"reflect"
	"testing"
)

func TestParserSimple(t *testing.T) {
	input := "let x = 10;"
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}

	expectedProg := ast.Program{
		Statements: []ast.Statement{
			ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Identifier: ast.IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "x"},
				},
				Expression: ast.IntExpression{
					Token: token.Token{Type: token.INT, Literal: "10"},
				},
			},
		},
	}

	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}
