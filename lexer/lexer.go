package lexer

import (
	"main/token"
	"strconv"
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

	l.skipWhitespace()
	if isNumber(l.ch) {
		nextToken = l.NumberToken()
	} else if isLetter(l.ch) {
		nextToken = l.literalToken()
	} else {
		nextToken = l.specialToken()
	}

	return nextToken
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) literalToken() token.Token {
	p := l.position
	for isLetter(l.ch) || isNumber(l.ch) {
		l.readChar()
	}

	literal := l.source[p:l.position]
	var t token.Token
	if mt, ok := token.MapSourceToToken(literal); ok {
		t = mt
	} else {
		t = token.Token{Type: token.IDENTIFIER, Literal: literal}
	}

	return t
}

func (l *Lexer) specialToken() token.Token {
	var t token.Token

	str := string(l.ch)
	switch l.ch {
	case '=':
		t = token.Token{Type: token.EQUAL, Literal: str}
	case '+':
		t = token.Token{Type: token.PLUS, Literal: str}
	case '<':
		t = token.Token{Type: token.LESS_THAN, Literal: str}
	case '>':
		t = token.Token{Type: token.GREATER_THAN, Literal: str}
	case '!':
		t = token.Token{Type: token.BANG, Literal: str}
	case '%':
		t = token.Token{Type: token.MODULUS, Literal: str}
	case '*':
		t = token.Token{Type: token.ASTERISK, Literal: str}
	case '-':
		t = token.Token{Type: token.MINUS, Literal: str}
	case '/':
		t = token.Token{Type: token.SLASH, Literal: str}
	case '&':
		t = token.Token{Type: token.AND, Literal: str}
	case '|':
		t = token.Token{Type: token.OR, Literal: str}
	case '(':
		t = token.Token{Type: token.LPAREN, Literal: str}
	case ')':
		t = token.Token{Type: token.RPAREN, Literal: str}
	case '{':
		t = token.Token{Type: token.LBRACKET, Literal: str}
	case '}':
		t = token.Token{Type: token.RBRACKET, Literal: str}
	case '[':
		t = token.Token{Type: token.LSQPAREN, Literal: str}
	case ']':
		t = token.Token{Type: token.RSQPAREN, Literal: str}
	case ',':
		t = token.Token{Type: token.COMMA, Literal: str}
	case ';':
		t = token.Token{Type: token.SEMICOLON, Literal: str}
	case 0:
		t = token.Token{Type: token.EOF, Literal: ""}
	default:
		t = token.Token{Type: token.ILLEGAL, Literal: ""}
	}

	l.readChar()
	return t
}

func (l *Lexer) NumberToken() token.Token {
	p := l.position
	for isNumber(l.ch) {
		l.readChar()
	}

	n := l.source[p:l.position]
	_, err := strconv.Atoi(n)
	if err != nil {
		return token.Token{Type: token.ILLEGAL, Literal: ""}
	}

	return token.Token{Type: token.INT, Literal: n}
}
