package parser

import (
	"main/ast"
	"main/token"
)

// you can assume that a statement functions:
// - Must have the currToken be the first token in its syntax
// - When it's done the peekToken will be the first token of the next statement

func (p *Parser) parseLetStatement() (ast.LetStatement, []error) {
	letToken := p.currToken
	p.nextToken()

	if !p.currTokenIs(token.IDENTIFIER) {
		return ast.LetStatement{}, []error{p.badTokenTypeError(token.IDENTIFIER)}
	}
	identExp := ast.IdentifierExpression{Token: p.currToken}
	p.nextToken()

	if !p.currTokenIs(token.EQUAL) {
		return ast.LetStatement{}, []error{p.badTokenTypeError(token.EQUAL)}
	}
	p.nextToken()

	valueExp, errs := p.parseExpression(LOWEST)
	if len(errs) != 0 {
		return ast.LetStatement{}, errs
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return ast.LetStatement{
			Token:      letToken,
			Identifier: identExp,
			Expression: valueExp,
		},
		nil
}

func (p *Parser) parseReturnStatement() (ast.ReturnStatement, []error) {
	returnToken := p.currToken
	var exp ast.Expression = nil
	p.nextToken()

	if !p.currTokenIs(token.SEMICOLON) {
		errs := []error{}
		exp, errs = p.parseExpression(LOWEST)
		if len(errs) != 0 {
			return ast.ReturnStatement{}, errs
		}
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return ast.ReturnStatement{
			Token:      returnToken,
			Expression: exp,
		},
		nil
}

func (p *Parser) parseExpressionStatement() (ast.ExpressionStatement, []error) {
	firstToken := p.currToken

	exp, errs := p.parseExpression(LOWEST)
	if len(errs) != 0 {
		return ast.ExpressionStatement{}, errs
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return ast.ExpressionStatement{
		Token:      firstToken,
		Expression: exp,
	}, nil
}
