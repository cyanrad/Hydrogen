package token

const (
	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// types
	IDENT = "IDENT" // x, y, foo, variables, ...
	INT   = "INT"   // integers: 1,2,3,...

	// keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// brackets
	LPAREN   = "("
	RPAREN   = ")"
	LBRACKET = "{"
	RBRACKET = "}"
	LSQPAREN = "["
	RSQPAREN = "]"

	// operators
	EQUAL = "="
	PLUS  = "+"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
