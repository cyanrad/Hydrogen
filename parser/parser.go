package parser

import (
	"fmt"
	"main/ast"
	"main/lexer"
	"main/token"
	"strconv"
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

func (p *Parser) ParseProgram() (ast.Program, []error) {
	statements := []ast.Statement{}
	errors := []error{} // it's better to have a general abstracted struct for this

	for !p.currTokenIs(token.EOF) {
		s, errs := p.ParseStatement()
		if s != nil {
			statements = append(statements, s)
		}
		errors = append(errors, errs...)
	}
	return ast.Program{Statements: statements}, errors
}

func (p *Parser) ParseBlockStatement() (ast.BlockStatement, []error) {
	// this function checks for starting and ending braces as such checks are common and not handling them here will lead to repeated code
	if !p.currTokenIs(token.LBRACKET) {
		return ast.BlockStatement{}, []error{p.badTokenTypeError(token.LBRACKET)}
	}
	t := p.currToken
	p.nextToken()

	statements := []ast.Statement{}
	errors := []error{} // it's better to have a general abstracted struct for this
	for !(p.currTokenIs(token.EOF) || p.currTokenIs(token.RBRACKET)) {
		s, errs := p.ParseStatement()
		if s != nil {
			statements = append(statements, s)
		}
		errors = append(errors, errs...)
	}

	return ast.BlockStatement{
		Token:      t,
		Statements: statements,
	}, errors
}

func (p *Parser) ParseStatement() (ast.Statement, []error) {
	var s ast.Statement
	var err error
	var errs []error // this is bad
	if p.currTokenIs(token.LET) {
		s, err = p.parseLetStatement()
	} else if p.currTokenIs(token.RETURN) {
		s, err = p.parseReturnStatement()
	} else if p.currTokenIs(token.IF) {
		s, errs = p.ParseIfExpression()
	} else {
		s, err = p.parseExpressionStatement()
	}

	if err != nil {
		errs = append(errs, err)
		s = nil
		p.skipToSemicolon()
	} else if len(errs) != 0 {
		s = nil
		p.skipToSemicolon()
	}
	// it is important for this to not be inside of the statement functions
	// as it is used to handle error cases
	p.nextToken()

	return s, errs
}

// you can assume that a statement functions have the currToken as the first token in it
// and that when it's done the currentToken is either EoF or the first token of the next statement
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

	valueExp, err := p.parseExpression(LOWEST)
	if err != nil {
		return ast.LetStatement{}, err
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

func (p *Parser) parseReturnStatement() (ast.ReturnStatement, error) {
	returnToken := p.currToken
	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return ast.ReturnStatement{}, err
	}
	p.nextToken()

	if !p.currTokenIs(token.SEMICOLON) {
		return ast.ReturnStatement{}, p.badTokenTypeError(token.SEMICOLON)
	}

	return ast.ReturnStatement{
			Token:      returnToken,
			Expression: exp,
		},
		nil
}

func (p *Parser) ParseIfExpression() (ast.IfExpression, []error) {
	blocks := []ast.BlockStatement{}
	conditions := []ast.Expression{}
	t := p.currToken

	// if & else if
	for {
		p.nextToken()

		exp, err := p.parseExpression(LOWEST)
		if err != nil {
			return ast.IfExpression{}, []error{err}
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
		p.nextToken()
	}

	return ast.IfExpression{
		Token:      t,
		Blocks:     blocks,
		Conditions: conditions,
	}, []error{}
}

func (p *Parser) parseExpressionStatement() (ast.ExpressionStatement, error) {
	firstToken := p.currToken

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return ast.ExpressionStatement{}, err
	}
	p.nextToken()

	if !p.currTokenIs(token.SEMICOLON) {
		return ast.ExpressionStatement{}, p.badTokenTypeError(token.SEMICOLON)
	}

	return ast.ExpressionStatement{
		Token:      firstToken,
		Expression: exp,
	}, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	var exp ast.Expression
	var err error

	if p.currTokenIsLegalPrefix() {
		exp, err = p.parsePrefixExpression()
	} else if p.currTokenIs(token.IDENTIFIER) {
		if p.peekTokenIs(token.LPAREN) {
			exp, err = p.parseFunctionExpression()
		} else {
			exp = p.parseIdentifierExpression()
		}
	} else if p.currTokenIs(token.BOOLEAN) {
		exp = p.parseBooleanExpression()
	} else if p.currTokenIs(token.INT) {
		exp, err = p.parseIntExpression()
	} else if p.currTokenIs(token.LPAREN) {
		exp, err = p.parseGroupedExpression()
	} else {
		return nil, fmt.Errorf("error - expected: expression - got: %s", p.currToken.Type)
	}
	if err != nil {
		return nil, err
	}

	for IsLegalInfixOperator(p.peekToken.Type) && precedence < p.peekPrecedence() {
		exp, err = p.parseInfixExpression(exp)
		if err != nil {
			return nil, err
		}
	}

	return exp, nil
}

func (p *Parser) parseIdentifierExpression() ast.IdentifierExpression {
	return ast.IdentifierExpression{
		Token: p.currToken,
	}
}

func (p *Parser) parseFunctionExpression() (ast.CallExpression, error) {
	identifier := p.parseIdentifierExpression()
	p.nextToken()
	p.nextToken()

	args := []ast.Expression{}
	// this whole thing (parsing comma seperated expressions) should be abstracted out
	if !p.peekTokenIs(token.RPAREN) {
		for {
			exp, err := p.parseExpression(LOWEST)
			if err != nil {
				return ast.CallExpression{}, err
			}

			args = append(args, exp)
			p.nextToken()

			if !p.currTokenIs(token.COMMA) {
				break
			}
			p.nextToken()
		}

		if !p.currTokenIs(token.RPAREN) {
			return ast.CallExpression{}, p.badTokenTypeError(token.RPAREN)
		}
	}

	return ast.CallExpression{
		Token:      identifier.Token,
		Identifier: identifier,
		Args:       args,
	}, nil
}

func (p *Parser) parseBooleanExpression() ast.BooleanExpression {
	return ast.BooleanExpression{
		Token: p.currToken,
	}
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	p.nextToken()

	if !p.currTokenIs(token.RPAREN) {
		return nil, p.badTokenTypeError(token.RPAREN)
	}

	return exp, nil
}

func (p *Parser) parseIntExpression() (ast.IntExpression, error) {
	// checking if it's parsable first
	_, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		return ast.IntExpression{}, fmt.Errorf("error - could not parse %q as integer", p.currToken.Literal)
	}

	return ast.IntExpression{
		Token: token.Token{Type: token.INT, Literal: p.currToken.Literal},
	}, nil
}

func (p *Parser) parsePrefixExpression() (ast.PrefixExpression, error) {
	operator := p.currToken

	p.nextToken()
	exp, err := p.parseExpression(PREFIX)
	if err != nil {
		return ast.PrefixExpression{}, err
	}

	return ast.PrefixExpression{
		Token:      operator,
		Expression: exp,
	}, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.InfixExpression, error) {
	p.nextToken()

	precedence := p.currPrecedence()
	operator := p.currToken
	p.nextToken()

	right, err := p.parseExpression(precedence)
	if err != nil {
		return ast.InfixExpression{}, err
	}

	return ast.InfixExpression{
		Token: operator,
		Left:  left,
		Right: right,
	}, nil
}
