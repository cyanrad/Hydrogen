package evaluator

import (
	"main/ast"
	"main/object"
	"strconv"
)

func EvalExpression(n ast.Expression) object.Object {
	switch exp := n.(type) {
	case ast.IntExpression:
		return evalInteger(exp)
	case ast.BooleanExpression:
		return evalBoolean(exp)
	}

	return nil
}

func evalInteger(node ast.IntExpression) object.Object {
	value, err := strconv.ParseInt(node.TokenLiteral(), 0, 64) // highly unlikely to fail
	if err != nil {
		panic("could not parse integer")
	}

	return &object.IntegerObj{Value: value}
}

func evalBoolean(node ast.BooleanExpression) object.Object {
	if node.TokenLiteral() == "true" {
		return &object.BooleanObj{Value: true}
	} else if node.TokenLiteral() == "false" {
		return &object.BooleanObj{Value: false}
	}
	panic("unknown boolean value")
}
