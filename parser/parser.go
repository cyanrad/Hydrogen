package parser

import (
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

	// loading tokens into curr and peek
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
	errors := []error{}

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
	if !p.currTokenIs(token.LBRACKET) {
		return ast.BlockStatement{}, []error{p.badTokenTypeError(token.LBRACKET)}
	}
	t := p.currToken
	p.nextToken()

	statements := []ast.Statement{}
	errors := []error{}
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
	var errs []error // this is bad

	// Parsing the statement
	if p.currTokenIs(token.LET) {
		s, errs = p.parseLetStatement()
	} else if p.currTokenIs(token.RETURN) {
		s, errs = p.parseReturnStatement()
	} else {
		s, errs = p.parseExpressionStatement()
	}

	// handling error cases
	if len(errs) != 0 {
		s = nil
		p.skipToSemicolon()
	}

	// it is important for this to be outside of statement functions
	// as it is used to handle error cases
	p.nextToken()

	return s, errs
}
