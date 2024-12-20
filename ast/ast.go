package ast

import (
	"main/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var sb strings.Builder
	for i, s := range p.Statements {
		sb.WriteString(s.String())

		if i < len(p.Statements)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

type LetStatement struct {
	// Statement
	Token      token.Token // token.LET
	Identifier IdentifierExpression
	Expression Expression
}

func (ls LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls LetStatement) statementNode()       {}
func (ls LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(ls.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(ls.Identifier.TokenLiteral())
	sb.WriteString(" = ")
	sb.WriteString(ls.Expression.TokenLiteral())
	sb.WriteString(";")

	return sb.String()
}

type ReturnStatement struct {
	// Statement
	Token      token.Token // token.RETURN
	Expression Expression
}

func (rs ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs ReturnStatement) statementNode()       {}
func (rs ReturnStatement) String() string {
	var sb strings.Builder

	sb.WriteString(rs.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(rs.Expression.String())
	sb.WriteString(";")

	return sb.String()
}

type IdentifierExpression struct {
	// Expression
	Token token.Token // token.Identifier + name
}

func (ie IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IdentifierExpression) expressionNode()      {}
func (ie IdentifierExpression) String() string       { return ie.TokenLiteral() }

type IntExpression struct {
	// Expression
	Token token.Token // token.INT + value
}

func (ie IntExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IntExpression) expressionNode()      {}
func (ie IntExpression) String() string       { return ie.TokenLiteral() }

type ExpressionStatement struct {
	// Statement
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es ExpressionStatement) statementNode()       {}
func (es ExpressionStatement) String() string       { return "" }
