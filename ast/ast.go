package ast

import (
	"main/token"
	"strings"
)

// this is atrocious
var INDENT_TRACKER = 0

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

type BlockStatement struct {
	Token      token.Token // token.LBracket
	Statements []Statement
}

func (bs BlockStatement) statementNode()       {}
func (bs BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs BlockStatement) String() string {
	var sb strings.Builder

	sb.WriteString("{\n")
	INDENT_TRACKER++
	for _, s := range bs.Statements {
		sb.WriteString(strings.Repeat("\t", INDENT_TRACKER) + s.String() + "\n")
	}
	INDENT_TRACKER--
	sb.WriteString(strings.Repeat("\t", INDENT_TRACKER) + "}")

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

type IfExpression struct {
	// Statement
	Token      token.Token // token.IF
	Conditions []Expression
	Blocks     []BlockStatement // len of Blocks should be Conditions+1 (in case last block is else)
}

func (ie IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IfExpression) expressionNode()      {}
func (ie IfExpression) String() string {
	var sb strings.Builder

	i := 0
	for ; i < len(ie.Conditions); i++ {
		if i != 0 {
			sb.WriteString(" else ")
		}
		sb.WriteString("if ")
		sb.WriteString(ie.Conditions[i].String())
		sb.WriteString(" ")
		sb.WriteString(ie.Blocks[i].String())
	}

	if len(ie.Conditions) < len(ie.Blocks) {
		sb.WriteString(" else ")
		sb.WriteString(ie.Blocks[i].String())
	}

	return sb.String()
}

type ExpressionStatement struct {
	// Statement
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es ExpressionStatement) statementNode()       {}
func (es ExpressionStatement) String() string       { return es.Expression.String() }

type IdentifierExpression struct {
	// Expression
	Token token.Token // token.Identifier + name
}

func (ie IdentifierExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IdentifierExpression) expressionNode()      {}
func (ie IdentifierExpression) String() string       { return ie.TokenLiteral() }

type BooleanExpression struct {
	// Expression
	Token token.Token // token.True or token.False
}

func (be BooleanExpression) TokenLiteral() string { return be.Token.Literal }
func (be BooleanExpression) expressionNode()      {}
func (be BooleanExpression) String() string       { return be.TokenLiteral() }

type IntExpression struct {
	// Expression
	Token token.Token // token.INT + value
}

func (ie IntExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IntExpression) expressionNode()      {}
func (ie IntExpression) String() string       { return ie.TokenLiteral() }

type PrefixExpression struct {
	// Expression
	Token      token.Token // operator toke (e.g. MINUS, EXCLAMATION)
	Expression Expression
}

type StringExpression struct {
	// Expression
	Token token.Token // token.STRING + value
}

func (se StringExpression) TokenLiteral() string { return se.Token.Literal }
func (se StringExpression) expressionNode()      {}
func (se StringExpression) String() string       { return se.TokenLiteral() }

func (pe PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe PrefixExpression) expressionNode()      {}
func (pe PrefixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(pe.TokenLiteral())
	sb.WriteString(pe.Expression.String())
	sb.WriteString(")")

	return sb.String()
}

type InfixExpression struct {
	// Expression
	Token token.Token // operator toke (e.g. MINUS, EXCLAMATION)
	Left  Expression
	Right Expression
}

func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) expressionNode()      {}
func (ie InfixExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ie.Left.String())
	sb.WriteString(" ")
	sb.WriteString(ie.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(ie.Right.String())
	sb.WriteString(")")

	return sb.String()
}

type CallExpression struct {
	// Expression
	Token      token.Token // the function identifier
	Identifier IdentifierExpression
	Args       []Expression
}

func (ce CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce CallExpression) expressionNode()      {}
func (ce CallExpression) String() string {
	var sb strings.Builder

	sb.WriteString(ce.TokenLiteral() + "(")
	for i, a := range ce.Args {
		sb.WriteString(a.String())
		if i != len(ce.Args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")

	return sb.String()
}

type ArrayExpression struct {
	// Expression
	Token token.Token // the [ token
	Elems []Expression
}

func (ae ArrayExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae ArrayExpression) expressionNode()      {}
func (ae ArrayExpression) String() string {
	var sb strings.Builder

	sb.WriteString("[")
	for i, a := range ae.Elems {
		sb.WriteString(a.String())
		if i != len(ae.Elems)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")

	return sb.String()
}

type IndexExpression struct {
	// Expression
	Token token.Token // the [ token
	Exp   Expression  // Expression attempting to be indexed
	Index Expression
}

func (ie IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IndexExpression) expressionNode()      {}
func (ie IndexExpression) String() string {
	var sb strings.Builder

	sb.WriteString("(")
	sb.WriteString(ie.Exp.String())
	sb.WriteString("[")
	sb.WriteString(ie.Index.String())
	sb.WriteString("])")

	return sb.String()
}

type FunctionExpression struct {
	// Expression
	Token token.Token // token.FUNCTION
	Args  []IdentifierExpression
	Body  BlockStatement
}

func (fe FunctionExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe FunctionExpression) expressionNode()      {}
func (fe FunctionExpression) String() string {
	var sb strings.Builder

	sb.WriteString("fn (")
	for i, a := range fe.Args {
		sb.WriteString(a.String())
		if i != len(fe.Args)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") ")
	sb.WriteString(fe.Body.String())

	return sb.String()
}

type KeyValuePair struct {
	// Expression
	Token token.Token // token.COLON
	Key   Expression
	Value Expression
}

func (kvp KeyValuePair) TokenLiteral() string { return kvp.Token.Literal }
func (kvp KeyValuePair) expressionNode()      {}
func (kvp KeyValuePair) String() string {
	var sb strings.Builder

	sb.WriteString(kvp.Key.String())
	sb.WriteString(": ")
	sb.WriteString(kvp.Value.String())

	return sb.String()
}

type HashExpression struct {
	// Expression
	Token token.Token // the { token
	Elems []KeyValuePair
}

func (he HashExpression) TokenLiteral() string { return he.Token.Literal }
func (he HashExpression) expressionNode()      {}
func (he HashExpression) String() string {
	var sb strings.Builder

	sb.WriteString("{")
	for i, a := range he.Elems {
		sb.WriteString(a.String())
		if i != len(he.Elems)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("}")

	return sb.String()
}
