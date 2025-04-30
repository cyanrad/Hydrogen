package parser

import (
	"fmt"
	"main/ast"
	"main/token"
	"strconv"
)

func (p *Parser) parseExpression(precedence int) (ast.Expression, []error) {
	var exp ast.Expression
	var errs []error

	if p.currTokenIsLegalPrefix() {
		exp, errs = p.parsePrefixExpression()
	} else if p.currTokenIs(token.IDENTIFIER) && p.peekTokenIs(token.LPAREN) {
		exp, errs = p.parseCallExpression()
	} else if p.currTokenIs(token.IDENTIFIER) {
		exp = p.parseIdentifierExpression()
	} else if p.currTokenIs(token.BOOLEAN) {
		exp = p.parseBooleanExpression()
	} else if p.currTokenIs(token.INT) {
		exp, errs = p.parseIntExpression()
	} else if p.currTokenIs(token.LPAREN) {
		exp, errs = p.parseGroupedExpression()
	} else if p.currTokenIs(token.IF) {
		exp, errs = p.ParseIfExpression()
	} else if p.currTokenIs(token.FUNCTION) {
		exp, errs = p.parseFunctionExpression()
	} else {
		return nil, []error{fmt.Errorf("error - expected: expression - got: %s", p.currToken.Type)}
	}
	if len(errs) != 0 {
		return nil, errs
	}

	for IsLegalInfixOperator(p.peekToken.Type) && precedence < p.peekPrecedence() {
		exp, errs = p.parseInfixExpression(exp)
		if len(errs) != 0 {
			return nil, errs
		}
	}

	return exp, nil
}

func (p *Parser) parsePrefixExpression() (ast.PrefixExpression, []error) {
	operator := p.currToken

	p.nextToken()
	exp, errs := p.parseExpression(PREFIX)
	if len(errs) != 0 {
		return ast.PrefixExpression{}, errs
	}

	return ast.PrefixExpression{
		Token:      operator,
		Expression: exp,
	}, nil
}

func (p *Parser) parseCallExpression() (ast.CallExpression, []error) {
	identifier := p.parseIdentifierExpression()
	p.nextToken()
	p.nextToken()

	args := []ast.Expression{}
	// this whole thing (parsing comma seperated expressions) should be abstracted out
	if !p.currTokenIs(token.RPAREN) {
		for {
			exp, errs := p.parseExpression(LOWEST)
			if len(errs) != 0 {
				return ast.CallExpression{}, errs
			}

			args = append(args, exp)
			p.nextToken()

			if !p.currTokenIs(token.COMMA) {
				break
			}
			p.nextToken()
		}

		if !p.currTokenIs(token.RPAREN) {
			return ast.CallExpression{}, []error{p.badTokenTypeError(token.RPAREN)}
		}
	}

	return ast.CallExpression{
		Token:      identifier.Token,
		Identifier: identifier,
		Args:       args,
	}, nil
}

func (p *Parser) parseIdentifierExpression() ast.IdentifierExpression {
	return ast.IdentifierExpression{
		Token: p.currToken,
	}
}

func (p *Parser) parseBooleanExpression() ast.BooleanExpression {
	return ast.BooleanExpression{
		Token: p.currToken,
	}
}

func (p *Parser) parseIntExpression() (ast.IntExpression, []error) {
	// checking if it's parsable first
	_, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		return ast.IntExpression{}, []error{fmt.Errorf("error - could not parse %q as integer", p.currToken.Literal)}
	}

	return ast.IntExpression{
		Token: token.Token{Type: token.INT, Literal: p.currToken.Literal},
	}, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, []error) {
	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	p.nextToken()

	if !p.currTokenIs(token.RPAREN) {
		return nil, []error{p.badTokenTypeError(token.RPAREN)}
	}

	return exp, nil
}

func (p *Parser) ParseIfExpression() (ast.IfExpression, []error) {
	blocks := []ast.BlockStatement{}
	conditions := []ast.Expression{}
	t := p.currToken

	// if & else if
	for {
		p.nextToken()

		exp, errs := p.parseExpression(LOWEST)
		if len(errs) != 0 {
			return ast.IfExpression{}, errs
		}
		conditions = append(conditions, exp)
		p.nextToken()

		b, errs := p.ParseBlockStatement()
		if len(errs) != 0 {
			return ast.IfExpression{}, errs
		}
		blocks = append(blocks, b)
		p.nextToken()

		if !(p.currTokenIs(token.ELSE) && p.peekTokenIs(token.IF)) {
			break
		}
		p.nextToken()
	}

	if p.currTokenIs(token.ELSE) {
		p.nextToken()
		b, errs := p.ParseBlockStatement()
		if len(errs) != 0 {
			return ast.IfExpression{}, errs
		}

		blocks = append(blocks, b)
	}

	return ast.IfExpression{
		Token:      t,
		Blocks:     blocks,
		Conditions: conditions,
	}, []error{}
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.InfixExpression, []error) {
	p.nextToken()

	precedence := p.currPrecedence()
	operator := p.currToken
	p.nextToken()

	right, errs := p.parseExpression(precedence)
	if len(errs) != 0 {
		return ast.InfixExpression{}, errs
	}

	return ast.InfixExpression{
		Token: operator,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseFunctionExpression() (ast.FunctionExpression, []error) {
	fn := p.currToken
	p.nextToken()

	if !p.currTokenIs(token.LPAREN) {
		return ast.FunctionExpression{}, []error{p.badTokenTypeError(token.LPAREN)}
	}
	p.nextToken()

	args := []ast.IdentifierExpression{}
	if !p.currTokenIs(token.RPAREN) {
		for {
			if !p.currTokenIs(token.IDENTIFIER) {
				return ast.FunctionExpression{}, []error{p.badTokenTypeError(token.IDENTIFIER)}
			}
			ident := p.parseIdentifierExpression()

			args = append(args, ident)
			p.nextToken()

			if p.currTokenIs(token.COMMA) {
				p.nextToken()
			} else if p.currTokenIs(token.RPAREN) {
				break
			} else {
				return ast.FunctionExpression{}, []error{p.badTokenTypeError(token.COMMA)}
			}
		}
	}
	p.nextToken()

	body, err := p.ParseBlockStatement()
	p.nextToken()
	if len(err) != 0 {
		return ast.FunctionExpression{}, err
	}

	return ast.FunctionExpression{
		Token: fn,
		Args:  args,
		Body:  body,
	}, nil
}
