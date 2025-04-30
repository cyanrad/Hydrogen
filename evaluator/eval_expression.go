package evaluator

import (
	"main/ast"
	"main/object"
	"strconv"
)

func EvalExpression(n ast.Expression, env Environment) object.Object {
	switch exp := n.(type) {
	case ast.IntExpression:
		return evalInteger(exp)
	case ast.BooleanExpression:
		return evalBoolean(exp)
	case ast.PrefixExpression:
		return evalPrefix(exp, env)
	case ast.InfixExpression:
		return evalInfix(exp, env)
	case ast.IfExpression:
		return evalIf(exp, env)
	case ast.IdentifierExpression:
		return evalIdentifier(exp, env)
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

func evalPrefix(node ast.PrefixExpression, env Environment) object.Object {
	exp := EvalExpression(node.Expression, env)

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

func evalInfix(node ast.InfixExpression, env Environment) object.Object {
	left := EvalExpression(node.Left, env)
	right := EvalExpression(node.Right, env)

	leftInt, leftOk := left.(*object.IntegerObj)
	rightInt, rightOk := right.(*object.IntegerObj)
	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "+":
			return &object.IntegerObj{Value: leftInt.Value + rightInt.Value}
		case "-":
			return &object.IntegerObj{Value: leftInt.Value - rightInt.Value}
		case "*":
			return &object.IntegerObj{Value: leftInt.Value * rightInt.Value}
		case "/":
			return &object.IntegerObj{Value: leftInt.Value / rightInt.Value}
		case "%":
			return &object.IntegerObj{Value: leftInt.Value % rightInt.Value}
		case "<":
			return &object.BooleanObj{Value: leftInt.Value < rightInt.Value}
		case "<=":
			return &object.BooleanObj{Value: leftInt.Value <= rightInt.Value}
		case ">":
			return &object.BooleanObj{Value: leftInt.Value > rightInt.Value}
		case ">=":
			return &object.BooleanObj{Value: leftInt.Value >= rightInt.Value}
		case "==":
			return &object.BooleanObj{Value: leftInt.Value == rightInt.Value}
		case "!=":
			return &object.BooleanObj{Value: leftInt.Value != rightInt.Value}
		case "&":
			return &object.IntegerObj{Value: leftInt.Value & rightInt.Value}
		case "|":
			return &object.IntegerObj{Value: leftInt.Value | rightInt.Value}
		default:
			panic("unknown infix operator: " + node.TokenLiteral())
		}
	}

	leftBool, leftOk := left.(*object.BooleanObj)
	rightBool, rightOk := right.(*object.BooleanObj)
	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "==":
			return &object.BooleanObj{Value: leftBool.Value == rightBool.Value}
		case "!=":
			return &object.BooleanObj{Value: leftBool.Value != rightBool.Value}
		case "&&":
			return &object.BooleanObj{Value: leftBool.Value && rightBool.Value}
		case "||":
			return &object.BooleanObj{Value: leftBool.Value || rightBool.Value}
		default:
			panic("unknown infix operator")
		}
	}

	panic("failed to match left and right types")
}

func evalIf(node ast.IfExpression, env Environment) object.Object {
	// looping over the conditions, return the block of the first condition that evaluates to true
	for i, condition := range node.Conditions {
		cond := EvalExpression(condition, env)
		if boolCond, ok := cond.(*object.BooleanObj); ok && boolCond.Value {
			return EvalStatement(node.Blocks[i], env)
		}
	}

	// the else condition
	if len(node.Blocks) > len(node.Conditions) {
		return EvalStatement(node.Blocks[len(node.Blocks)-1], env)
	}

	return &object.NullObj{}
}

func evalIdentifier(node ast.IdentifierExpression, env Environment) object.Object {
	// check if the identifier is a variable in the environment
	if obj := env.GetVariable(node.TokenLiteral()); obj != nil {
		return obj
	}

	panic("unknown identifier: " + node.TokenLiteral())
}
