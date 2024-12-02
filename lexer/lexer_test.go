package lexer

import (
	"main/token"
	"testing"
)

func TestGetNextTokenSpecialCharacters(t *testing.T) {
	input := "=+(){}[],;"

	expected := []token.Token{
		{Type: token.EQUAL, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACKET, Literal: "{"},
		{Type: token.RBRACKET, Literal: "}"},
		{Type: token.LSQPAREN, Literal: "["},
		{Type: token.RSQPAREN, Literal: "]"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}

	l := CreateLexer(input)
	for i, et := range expected {
		nt := l.GetNextToken()
		// t.Log(nt)
		if !(nt.Type == et.Type && nt.Literal == et.Literal) {
			t.Fatalf("test[%d] - mismatch between expected and actual token - expected: %s, %s - actual: %s, %s", i, et.Type, et.Literal, nt.Type, nt.Literal)
		}
	}
}

func TestGetNextTokenCode(t *testing.T) {
	input := `let    five = 5;
 let ten = 10;
 let add = fn(x, y) {
 x + y;
 };
 let result = add(five, ten);
 `

	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACKET, Literal: "{"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACKET, Literal: "}"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "result"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}

	l := CreateLexer(input)
	for i, et := range expected {
		nt := l.GetNextToken()
		// t.Log(nt)
		if !(nt.Type == et.Type && nt.Literal == et.Literal) {
			t.Fatalf("test[%d] - mismatch between expected and actual token - expected: %s, %s - actual: %s, %s", i, et.Type, et.Literal, nt.Type, nt.Literal)
		}
	}
}
