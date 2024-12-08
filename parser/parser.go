package parser

import (
	"fmt"
	"main/ast"
	"main/lexer"
	"main/token"
)

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

func (p *Parser) ParseProgram() (ast.Program, error) {
	prog := ast.Program{Statements: []ast.Statement{}}

	for p.currToken.Type != token.EOF {
		if p.currToken.Type == token.LET {
			s, err := p.parseLetStatement()
			if err != nil {
				return ast.Program{}, err
			}
			prog.Statements = append(prog.Statements, s)
		}

		if p.currToken.Type != token.SEMICOLON {
			return ast.Program{}, fmt.Errorf("error - expected: ; - found: %v", p.currToken)
		}
		p.nextToken()
	}

	return prog, nil
}

// you can assume that the parse functions have the currToken as the first token in it
// you can assume that by the end that the currtoken should equal to ;
func (p *Parser) parseLetStatement() (ast.LetStatement, error) {
	letToken := p.currToken

	p.nextToken()
	identToken := p.currToken
	if identToken.Type != token.IDENTIFIER {
		return ast.LetStatement{}, fmt.Errorf("error - expected: identifier token after let - found: %v", identToken)
	}

	p.nextToken()
	assignToken := p.currToken
	if assignToken.Type != token.EQUAL {
		return ast.LetStatement{}, fmt.Errorf("error - expected: = operator - found: %v", assignToken)
	}

	p.nextToken()
	exp, err := p.parseMathExpression()
	if err != nil {
		return ast.LetStatement{}, err
	}

	p.nextToken()
	return ast.LetStatement{
			Token:      letToken,
			Identifier: ast.IdentifierExpression{Token: identToken},
			Expression: exp,
		},
		nil
}

func (p *Parser) parseMathExpression() (ast.IntExpression, error) {
	intToken := p.currToken
	if intToken.Type != token.INT {
		return ast.IntExpression{}, fmt.Errorf("error - expected: int - found: %v", intToken)
	}

	return ast.IntExpression{Token: intToken}, nil
}
