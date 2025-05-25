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
	} else if p.currTokenIs(token.STRING) {
		exp = p.parseStringExpression()
	} else if p.currTokenIs(token.LPAREN) {
		exp, errs = p.parseGroupedExpression()
	} else if p.currTokenIs(token.IF) {
		exp, errs = p.ParseIfExpression()
	} else if p.currTokenIs(token.FUNCTION) {
		exp, errs = p.parseFunctionExpression()
	} else if p.currTokenIs(token.LSQPAREN) {
		exp, errs = p.parseArrayExpression()
	} else if p.currTokenIs(token.LBRACKET) {
		exp, errs = p.ParseHashExpression()
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

	args, errs := p.parseExpressionList(token.RPAREN) // Fixed typo: Seperated -> Separated
	if len(errs) != 0 {
		return ast.CallExpression{}, errs
	}

	return ast.CallExpression{
		Token:      identifier.Token,
		Identifier: identifier,
		Args:       args,
	}, nil
}

func (p *Parser) parseIdentifierExpression() ast.IdentifierExpression {
	return ast.IdentifierExpression{Token: p.currToken}
}

func (p *Parser) parseBooleanExpression() ast.BooleanExpression {
	return ast.BooleanExpression{Token: p.currToken}
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

func (p *Parser) parseStringExpression() ast.StringExpression {
	return ast.StringExpression{Token: p.currToken}
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

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, []error) {
	p.nextToken()

	precedence := p.currPrecedence()
	operator := p.currToken

	var right ast.Expression
	var errs []error

	if operator.Type == token.LSQPAREN {
		right, errs = p.parseIndexExpression(left)
		if len(errs) != 0 {
			return ast.IndexExpression{}, errs
		}
		return right, nil
	}

	p.nextToken()
	right, errs = p.parseExpression(precedence)
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

func (p *Parser) parseArrayExpression() (ast.Expression, []error) {
	lqp := p.currToken
	p.nextToken()

	elems, errs := p.parseExpressionList(token.RSQPAREN)
	if len(errs) != 0 {
		return ast.ArrayExpression{}, errs
	}

	return ast.ArrayExpression{
		Token: lqp,
		Elems: elems,
	}, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.IndexExpression, []error) {
	lqp := p.currToken
	p.nextToken()

	indexExp, errs := p.parseExpression(LOWEST)
	if len(errs) != 0 {
		return ast.IndexExpression{}, errs
	}
	p.nextToken()

	return ast.IndexExpression{
		Token: lqp,
		Exp:   left,
		Index: indexExp,
	}, nil

}

func (p *Parser) parseExpressionList(terminationToken token.TokenType) ([]ast.Expression, []error) {
	args := []ast.Expression{}

	for !p.currTokenIs(terminationToken) {
		// parsing expression
		exp, errs := p.parseExpression(LOWEST)
		if len(errs) != 0 {
			return nil, errs
		}

		args = append(args, exp)
		p.nextToken()

		// parsing comma
		if !p.currTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
	}

	return args, nil
}

func (p *Parser) ParseHashExpression() (ast.HashExpression, []error) {
	lb := p.currToken
	p.nextToken()

	elems := []ast.KeyValuePair{}
	for !p.currTokenIs(token.RBRACKET) {
		// parsing expression
		exp, errs := p.parseKeyValuePairExpression()
		if len(errs) != 0 {
			return ast.HashExpression{}, errs
		}

		elems = append(elems, exp)
		p.nextToken()

		// parsing comma
		if !p.currTokenIs(token.COMMA) {
			break
		}
		p.nextToken()
	}

	return ast.HashExpression{
		Token: lb,
		Elems: elems,
	}, nil
}

func (p *Parser) parseKeyValuePairExpression() (ast.KeyValuePair, []error) {
	key, errs := p.parseExpression(LOWEST)
	if len(errs) != 0 {
		return ast.KeyValuePair{}, errs
	}

	p.nextToken()
	if !p.currTokenIs(token.COLON) {
		return ast.KeyValuePair{}, []error{p.badTokenTypeError(token.COLON)}
	}
	colon := p.currToken

	p.nextToken()
	value, errs := p.parseExpression(LOWEST)
	if len(errs) != 0 {
		return ast.KeyValuePair{}, errs
	}

	return ast.KeyValuePair{
		Token: colon,
		Key:   key,
		Value: value,
	}, nil
}
