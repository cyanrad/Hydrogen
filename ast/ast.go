package ast

import "main/token"

type Node interface {
	TokenLiteral() string
}
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// returns the token literal of the first statement
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	// Statement
	Token      token.Token // token.LET
	Identifier IdentifierExpression
	Expression Expression
}

func (ls LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls LetStatement) statementNode()       {}

type ReturnStatement struct {
	// Statement
	Token      token.Token // token.RETURN
	Expression Expression
}

func (rs ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs ReturnStatement) statementNode()       {}

type IdentifierExpression struct {
	// Expression
	Token token.Token // token.Identifier + name
}

func (ie IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IdentifierExpression) expressionNode()      {}

type IntExpression struct {
	// Expression
	Token token.Token // token.INT + value
}

func (ie IntExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IntExpression) expressionNode()      {}
