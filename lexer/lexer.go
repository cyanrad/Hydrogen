package lexer

import (
	"main/token"
)

type Lexer struct {
	source       string
	position     int // pointer of current position in input
	readPosition int // pointer of char we're currently reading
	ch           byte
}

func CreateLexer(input string) *Lexer {
	l := &Lexer{source: input}
	l.readChar()
	return l

}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.source) {
		l.ch = 0
	} else {
		l.ch = l.source[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) GetNextToken() token.Token {
	var nextToken token.Token
	switch l.ch {
	case '=':
		nextToken = token.Token{Type: token.EQUAL, Literal: "="}
	case '+':
		nextToken = token.Token{Type: token.PLUS, Literal: "+"}
	case '(':
		nextToken = token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		nextToken = token.Token{Type: token.RPAREN, Literal: ")"}
	case '{':
		nextToken = token.Token{Type: token.LBRACKET, Literal: "{"}
	case '}':
		nextToken = token.Token{Type: token.RBRACKET, Literal: "}"}
	case '[':
		nextToken = token.Token{Type: token.LSQPAREN, Literal: "["}
	case ']':
		nextToken = token.Token{Type: token.RSQPAREN, Literal: "]"}
	case ',':
		nextToken = token.Token{Type: token.COMMA, Literal: ","}
	case ';':
		nextToken = token.Token{Type: token.SEMICOLON, Literal: ";"}
	case 0:
		nextToken = token.Token{Type: token.EOF, Literal: ""}
	default:
		nextToken = token.Token{Type: token.ILLEGAL, Literal: ""}
	}

	l.readChar()
	return nextToken
}
