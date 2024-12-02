package token

const (
	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// types
	IDENTIFIER = "IDENTIFIER" // x, y, foo, variables, ...
	INT        = "INT"        // integers: 1,2,3,...

	// keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	FOR      = "FOR"

	// quotes
	SINGLE_QUOTE  = "'"
	DOUBLE_QUOTES = "\""

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

	// math operators
	EQUAL        = "="
	PLUS         = "+"
	MINUS        = "-"
	DIVIDE       = "/"
	MULTIPLY     = "*"
	MODULUS      = "%"
	GREATER_THAN = ">"
	LESS_THAN    = "<"

	// logic operators
	AND                   = "&"
	OR                    = "|"
	CONDITIONAL_AND       = "&&"
	CONDITIONAL_OR        = "||"
	BANG                  = "!"
	CONDITIONAL_EQUAL     = "=="
	CONDITIONAL_NOT_EQUAL = "!="
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var tokenSourceMapping map[string]Token = map[string]Token{
	"let": {Type: LET, Literal: "let"},
	"fn":  {Type: FUNCTION, Literal: "fn"},
	"if":  {Type: IF, Literal: "if"},
	"for": {Type: FOR, Literal: "for"},
}

func MapSourceToToken(sourceStr string) (Token, bool) {
	token, ok := tokenSourceMapping[sourceStr]
	return token, ok
}
