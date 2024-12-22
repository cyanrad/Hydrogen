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
	prog := ast.Program{Statements: []ast.Statement{}}
	errors := []error{}

	for !p.currTokenIs(token.EOF) {
		var s ast.Statement
		var err error
		if p.currTokenIs(token.LET) {
			s, err = p.parseLetStatement()
		} else if p.currTokenIs(token.RETURN) {
			s, err = p.parseReturnStatement()
		} else {
			s, err = p.parseExpressionStatement()
		}

		if err != nil {
			errors = append(errors, err)
			p.skipToSemicolon()
		} else if s != nil {
			prog.Statements = append(prog.Statements, s)
		}

		// it is important for this to not be inside of the statement functions
		// as it is used to handle error cases
		p.nextToken()
	}

	return prog, errors
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
		exp = p.parseIdentifierExpression()
	} else if p.currTokenIs(token.INT) {
		exp, err = p.parseIntExpression()
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
		Token: token.Token{Type: token.IDENTIFIER, Literal: p.currToken.Literal},
	}
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
