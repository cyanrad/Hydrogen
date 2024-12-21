package ast

import (
	"main/token"
)

var legalPrefixOperator = map[token.TokenType]struct{}{
	token.MINUS:     {},
	token.BANG:      {},
	token.INCREMENT: {},
	token.DECREMENT: {},
}

func IsLegalPrefixOperator(t token.TokenType) bool {
	_, ok := legalPrefixOperator[t]
	return ok
}

var legalInfexOperator = map[token.TokenType]struct{}{
	token.MINUS:                 {},
	token.PLUS:                  {},
	token.LESS_THAN:             {},
	token.GREATER_THAN:          {},
	token.MODULUS:               {},
	token.ASTERISK:              {},
	token.SLASH:                 {},
	token.AND:                   {},
	token.OR:                    {},
	token.CONDITIONAL_AND:       {},
	token.CONDITIONAL_OR:        {},
	token.CONDITIONAL_EQUAL:     {},
	token.CONDITIONAL_NOT_EQUAL: {},
	token.GREATER_THAN_EQUAL:    {},
	token.LESS_THAN_EQUAL:       {},
}

func IsLegalInfixOperator(t token.TokenType) bool {
	_, ok := legalInfexOperator[t]
	return ok
}
