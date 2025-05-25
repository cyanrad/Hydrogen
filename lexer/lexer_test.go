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
!-/*5;
5 < 10 > 5;
&|%;

if(x) {return true}
else {return false}
for

^
== != < <>>
===
+=-=++--&&||>=<=

"hello world"
[1, 2];
{"foo": "bar"}
 `

	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "five"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "ten"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "add"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "y"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACKET, Literal: "{"},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENTIFIER, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACKET, Literal: "}"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENTIFIER, Literal: "result"},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.IDENTIFIER, Literal: "add"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "five"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENTIFIER, Literal: "ten"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.BANG, Literal: "!"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "5"},
		{Type: token.LESS_THAN, Literal: "<"},
		{Type: token.INT, Literal: "10"},
		{Type: token.GREATER_THAN, Literal: ">"},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.AND, Literal: "&"},
		{Type: token.OR, Literal: "|"},
		{Type: token.MODULUS, Literal: "%"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IF, Literal: "if"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENTIFIER, Literal: "x"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACKET, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.BOOLEAN, Literal: "true"},
		{Type: token.RBRACKET, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.LBRACKET, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.BOOLEAN, Literal: "false"},
		{Type: token.RBRACKET, Literal: "}"},
		{Type: token.FOR, Literal: "for"},
		{Type: token.ILLEGAL, Literal: "^"},
		{Type: token.CONDITIONAL_EQUAL, Literal: "=="},
		{Type: token.CONDITIONAL_NOT_EQUAL, Literal: "!="},
		{Type: token.LESS_THAN, Literal: "<"},
		{Type: token.LESS_THAN, Literal: "<"},
		{Type: token.GREATER_THAN, Literal: ">"},
		{Type: token.GREATER_THAN, Literal: ">"},
		{Type: token.CONDITIONAL_EQUAL, Literal: "=="},
		{Type: token.EQUAL, Literal: "="},
		{Type: token.PLUS_EQUAL, Literal: "+="},
		{Type: token.MINUS_EQUAL, Literal: "-="},
		{Type: token.INCREMENT, Literal: "++"},
		{Type: token.DECREMENT, Literal: "--"},
		{Type: token.CONDITIONAL_AND, Literal: "&&"},
		{Type: token.CONDITIONAL_OR, Literal: "||"},
		{Type: token.GREATER_THAN_EQUAL, Literal: ">="},
		{Type: token.LESS_THAN_EQUAL, Literal: "<="},
		{Type: token.STRING, Literal: "hello world"},
		{Type: token.LSQPAREN, Literal: "["},
		{Type: token.INT, Literal: "1"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.INT, Literal: "2"},
		{Type: token.RSQPAREN, Literal: "]"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LBRACKET, Literal: "{"},
		{Type: token.STRING, Literal: "foo"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.STRING, Literal: "bar"},
		{Type: token.RBRACKET, Literal: "}"},
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
