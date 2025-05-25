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

	l.skipWhitespace()
	if isNumber(l.ch) {
		nextToken = l.NumberToken()
	} else if isLetter(l.ch) {
		nextToken = l.literalToken()
	} else if l.ch == '"' {
		nextToken = l.StringToken()
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
	if mt, ok := token.MapSourceToKeyword(literal); ok {
		t = mt
	} else {
		t = token.Token{Type: token.IDENTIFIER, Literal: literal}
	}

	return t
}

func (l *Lexer) StringToken() token.Token {
	l.readChar()
	p := l.position
	for l.readPosition <= len(l.source) && l.ch != '"' {
		l.readChar()

		if l.ch == '\n' {
			return token.Token{Type: token.ILLEGAL, Literal: ""}
		}
	}

	if l.ch == '"' {
		l.readChar()
	} else {
		return token.Token{Type: token.ILLEGAL, Literal: ""} // checking for proper string termination
	}

	return token.Token{Type: token.STRING, Literal: l.source[p : l.position-1]}
}

func (l *Lexer) specialToken() token.Token {
	t := token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
	if l.readPosition > len(l.source) {
		t = token.Token{Type: token.EOF, Literal: ""}
	}

	for p := l.position; l.readPosition <= len(l.source); l.readChar() {
		if mt, ok := token.MapSourceToSpecial(l.source[p:l.readPosition]); ok {
			t = mt
		} else {
			break
		}
	}

	// this is so stupid - to handle when pointer movement when an illegal character
	if t.Type == token.ILLEGAL {
		l.readChar()
	}
	return t
}

func (l *Lexer) NumberToken() token.Token {
	p := l.position
	for isNumber(l.ch) {
		l.readChar()
	}

	n := l.source[p:l.position]
	return token.Token{Type: token.INT, Literal: n}
}
