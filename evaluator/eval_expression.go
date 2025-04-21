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
	case ast.PrefixExpression:
		return evalPrefix(exp)
	case ast.InfixExpression:
		return evalInfix(exp)
	default:
		panic("unknown expression type")
	}
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

func evalPrefix(node ast.PrefixExpression) object.Object {
	exp := EvalExpression(node.Expression)

	if intexp, ok := exp.(*object.IntegerObj); ok {
		switch node.TokenLiteral() {
		case "-":
			intexp.Value = -intexp.Value
		case "++":
			intexp.Value += 1
		case "--":
			intexp.Value -= 1
		default:
			panic("unknown prefix operator")
		}

		return intexp
	} else if boolExp, ok := exp.(*object.BooleanObj); ok {
		switch node.TokenLiteral() {
		case "!":
			boolExp.Value = !boolExp.Value
		default:
			panic("unknown prefix operator")
		}

		return boolExp
	} else {
		panic("unknown expression type")
	}
}

func evalInfix(node ast.InfixExpression) object.Object {
	left := EvalExpression(node.Left)
	right := EvalExpression(node.Right)

	leftType, leftOk := left.(*object.IntegerObj)
	rightType, rightOk := right.(*object.IntegerObj)

	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "+":
			return &object.IntegerObj{Value: leftType.Value + rightType.Value}
		case "-":
			return &object.IntegerObj{Value: leftType.Value - rightType.Value}
		case "*":
			return &object.IntegerObj{Value: leftType.Value * rightType.Value}
		case "/":
			return &object.IntegerObj{Value: leftType.Value / rightType.Value}
		case "%":
			return &object.IntegerObj{Value: leftType.Value % rightType.Value}
		default:
			panic("unknown infix operator")
		}
	}

	panic("failed to match left and right types")
}
