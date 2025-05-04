package token

const (
	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// types
	IDENTIFIER = "IDENTIFIER" // x, y, foo, variables, ...
	INT        = "INT"        // integers: 1,2,3,...
	BOOLEAN    = "BOOLEAN"    // true or false
	STRING     = "STRING"     // string literals: "hello"

	// keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	RETURN   = "RETURN"

	// quotes
	SINGLE_QUOTE  = "'"
	DOUBLE_QUOTES = "\""

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// brackets
	LPAREN   = "("
	RPAREN   = ")"
	LBRACKET = "{"
	RBRACKET = "}"
	LSQPAREN = "["
	RSQPAREN = "]"

	// operators
	EQUAL       = "="
	PLUS_EQUAL  = "+="
	MINUS_EQUAL = "-="
	PLUS        = "+"
	MINUS       = "-"
	INCREMENT   = "++"
	DECREMENT   = "--"
	SLASH       = "/"
	ASTERISK    = "*"
	MODULUS     = "%"

	// logic operators
	AND                   = "&"
	OR                    = "|"
	CONDITIONAL_AND       = "&&"
	CONDITIONAL_OR        = "||"
	BANG                  = "!"
	CONDITIONAL_EQUAL     = "=="
	CONDITIONAL_NOT_EQUAL = "!="
	GREATER_THAN          = ">"
	LESS_THAN             = "<"
	GREATER_THAN_EQUAL    = ">="
	LESS_THAN_EQUAL       = "<="
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywordTokenMap map[string]Token = map[string]Token{
	"let":    {Type: LET, Literal: "let"},
	"fn":     {Type: FUNCTION, Literal: "fn"},
	"if":     {Type: IF, Literal: "if"},
	"else":   {Type: ELSE, Literal: "else"},
	"for":    {Type: FOR, Literal: "for"},
	"true":   {Type: BOOLEAN, Literal: "true"},
	"false":  {Type: BOOLEAN, Literal: "false"},
	"return": {Type: RETURN, Literal: "return"},
}

var specialTokenMap map[string]Token = map[string]Token{
	// single char
	"=": {Type: EQUAL, Literal: "="},
	"+": {Type: PLUS, Literal: "+"},
	"<": {Type: LESS_THAN, Literal: "<"},
	">": {Type: GREATER_THAN, Literal: ">"},
	"!": {Type: BANG, Literal: "!"},
	"%": {Type: MODULUS, Literal: "%"},
	"*": {Type: ASTERISK, Literal: "*"},
	"-": {Type: MINUS, Literal: "-"},
	"/": {Type: SLASH, Literal: "/"},
	"&": {Type: AND, Literal: "&"},
	"|": {Type: OR, Literal: "|"},
	"(": {Type: LPAREN, Literal: "("},
	")": {Type: RPAREN, Literal: ")"},
	"{": {Type: LBRACKET, Literal: "{"},
	"}": {Type: RBRACKET, Literal: "}"},
	"[": {Type: LSQPAREN, Literal: "["},
	"]": {Type: RSQPAREN, Literal: "]"},
	",": {Type: COMMA, Literal: ","},
	";": {Type: SEMICOLON, Literal: ";"},
	":": {Type: COLON, Literal: ":"},

	// double char
	"+=": {Type: PLUS_EQUAL, Literal: "+="},
	"-=": {Type: MINUS_EQUAL, Literal: "-="},
	"++": {Type: INCREMENT, Literal: "++"},
	"--": {Type: DECREMENT, Literal: "--"},
	"&&": {Type: CONDITIONAL_AND, Literal: "&&"},
	"||": {Type: CONDITIONAL_OR, Literal: "||"},
	"==": {Type: CONDITIONAL_EQUAL, Literal: "=="},
	"!=": {Type: CONDITIONAL_NOT_EQUAL, Literal: "!="},
	">=": {Type: GREATER_THAN_EQUAL, Literal: ">="},
	"<=": {Type: LESS_THAN_EQUAL, Literal: "<="},
}

func MapSourceToKeyword(sourceStr string) (Token, bool) {
	token, ok := keywordTokenMap[sourceStr]
	return token, ok
}

func MapSourceToSpecial(sourceStr string) (Token, bool) {
	token, ok := specialTokenMap[sourceStr]
	return token, ok
}
