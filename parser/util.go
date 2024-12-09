package parser

import (
	"fmt"
	"main/token"
)

func (p *Parser) badTokenTypeError(expected token.TokenType) error {
	return fmt.Errorf("error - expected: %s - got: %s", expected, p.currToken.Type)
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

//	func (p *Parser) peekTokenIs(t token.TokenType) bool {
//		return p.peekToken.Type == t
//	}
func (p *Parser) skipToSemicolon() {
	for !p.currTokenIs(token.SEMICOLON) && !p.currTokenIs(token.EOF) {
		p.nextToken()
	}
}
