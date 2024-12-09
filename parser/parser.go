package parser

import (
	"fmt"
	"main/ast"
	"main/lexer"
	"main/token"
)

var legalMathOperators = map[token.TokenType]struct{}{
	token.MODULUS:  {},
	token.ASTERISK: {},
	token.SLASH:    {},
	token.PLUS:     {},
	token.MINUS:    {},
}

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func CreateParser(l *lexer.Lexer) Parser {
	p := Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.GetNextToken()
}

func (p *Parser) ParseProgram() (ast.Program, []error) {
	prog := ast.Program{Statements: []ast.Statement{}}
	errors := []error{}

	for !p.currTokenIs(token.EOF) {
		var s ast.Statement
		var err error
		if p.currTokenIs(token.LET) {
			s, err = p.parseLetStatement()
		}

		if err != nil {
			errors = append(errors, err)
			p.skipToSemicolon()
		} else {
			prog.Statements = append(prog.Statements, s)
		}

		p.nextToken()
	}

	return prog, errors
}

// you can assume that the parse functions have the currToken as the first token in it
// you can assume that by the end that the currtoken should equal to ; if viable
func (p *Parser) parseLetStatement() (ast.LetStatement, error) {
	letToken := p.currToken

	p.nextToken()
	if !p.currTokenIs(token.IDENTIFIER) {
		return ast.LetStatement{}, p.badTokenTypeError(token.IDENTIFIER)
	}
	identExp := ast.IdentifierExpression{Token: p.currToken}

	p.nextToken()
	if !p.currTokenIs(token.EQUAL) {
		return ast.LetStatement{}, p.badTokenTypeError(token.EQUAL)
	}

	p.nextToken()
	var valueExp ast.IntExpression
	var err error
	if p.currTokenIs(token.INT) {
		valueExp, err = p.parseMathExpression()
		if err != nil {
			return ast.LetStatement{}, err
		}
	} else {
		return ast.LetStatement{}, fmt.Errorf("error - expected: expression - got: %s", p.currToken.Type)
	}

	p.nextToken()
	if !p.currTokenIs(token.SEMICOLON) {
		return ast.LetStatement{}, p.badTokenTypeError(token.SEMICOLON)
	}

	return ast.LetStatement{
			Token:      letToken,
			Identifier: identExp,
			Expression: valueExp,
		},
		nil
}

func (p *Parser) parseMathExpression() (ast.IntExpression, error) {
	intToken := p.currToken
	if intToken.Type != token.INT {
		return ast.IntExpression{}, fmt.Errorf("error - expected: int - got: %v", intToken)
	}

	return ast.IntExpression{Token: intToken}, nil
}
