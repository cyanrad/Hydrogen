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
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) currTokenIsLegalPrefix() bool {
	return IsLegalPrefixOperator(p.currToken.Type)
}

func (p *Parser) skipToSemicolon() {
	for !p.currTokenIs(token.SEMICOLON) && !p.currTokenIs(token.EOF) {
		p.nextToken()
	}
}

func (p *Parser) currPrecedence() int {
	if precedence, ok := precedences[p.currToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

var legalPrefixOperator = map[token.TokenType]struct{}{
	token.MINUS:     {},
	token.BANG:      {},
	token.INCREMENT: {},
	token.DECREMENT: {},
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func IsLegalPrefixOperator(t token.TokenType) bool {
	_, ok := legalPrefixOperator[t]
	return ok
}

var legalInfexOperator = map[token.TokenType]struct{}{
	token.MINUS:                 {},
	token.PLUS:                  {},
	token.LESS_THAN:             {},
	token.GREATER_THAN:          {},
	token.MODULUS:               {},
	token.ASTERISK:              {},
	token.SLASH:                 {},
	token.AND:                   {},
	token.OR:                    {},
	token.CONDITIONAL_AND:       {},
	token.CONDITIONAL_OR:        {},
	token.CONDITIONAL_EQUAL:     {},
	token.CONDITIONAL_NOT_EQUAL: {},
	token.GREATER_THAN_EQUAL:    {},
	token.LESS_THAN_EQUAL:       {},
	token.LSQPAREN:              {},
}

func IsLegalInfixOperator(t token.TokenType) bool {
	_, ok := legalInfexOperator[t]
	return ok
}

const (
	_           int = iota
	LOWEST          // _ (black identifier)
	EQUALS          // ==
	LESSGREATER     // > or <
	AND             // &&
	OR              // ||
	BITWISE         // & |
	SUM             // +
	PRODUCT         // *
	PREFIX          // -x or !x
	CALL            // myFunc(x)
	INDEX           // myArray[x]
)

var precedences = map[token.TokenType]int{
	token.CONDITIONAL_EQUAL:     EQUALS,
	token.CONDITIONAL_NOT_EQUAL: EQUALS,
	token.LESS_THAN:             LESSGREATER,
	token.LESS_THAN_EQUAL:       LESSGREATER,
	token.GREATER_THAN:          LESSGREATER,
	token.GREATER_THAN_EQUAL:    LESSGREATER,
	token.CONDITIONAL_AND:       AND,
	token.CONDITIONAL_OR:        OR,
	token.AND:                   BITWISE,
	token.OR:                    BITWISE,
	token.PLUS:                  SUM,
	token.MINUS:                 SUM,
	token.SLASH:                 PRODUCT,
	token.ASTERISK:              PRODUCT,
	token.MODULUS:               PRODUCT,
	token.LSQPAREN:              INDEX,
}
