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
	case ast.StringExpression:
		return evalString(exp)
	case ast.PrefixExpression:
		return evalPrefix(exp, env)
	case ast.InfixExpression:
		return evalInfix(exp, env)
	case ast.IfExpression:
		return evalIf(exp, env)
	case ast.IdentifierExpression:
		return evalIdentifier(exp, env)
	case ast.FunctionExpression:
		return evalFunction(exp)
	case ast.CallExpression:
		return evalCall(exp, env)
	case ast.ArrayExpression:
		return evalArray(exp, env)
	case ast.IndexExpression:
		return evalIndex(exp, env)
	case ast.HashExpression:
		return evalHash(exp, env)
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

func evalString(node ast.StringExpression) object.Object {
	return &object.StringObj{Value: node.TokenLiteral()}
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

	leftStr, leftOk := left.(*object.StringObj)
	rightStr, rightOk := right.(*object.StringObj)
	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "==":
			return &object.BooleanObj{Value: leftStr.Value == rightStr.Value}
		case "!=":
			return &object.BooleanObj{Value: leftStr.Value != rightStr.Value}
		case "+":
			return &object.StringObj{Value: leftStr.Value + rightStr.Value}
		default:
			panic("unknown infix operator")
		}
	}

	panic("failed to match left and right types")
}

func evalIf(node ast.IfExpression, env Environment) object.Object {
	tempEnv := NewEnclosedEnvironment(env)

	// looping over the conditions, return the block of the first condition that evaluates to true
	for i, condition := range node.Conditions {
		cond := EvalExpression(condition, env)
		if boolCond, ok := cond.(*object.BooleanObj); ok && boolCond.Value {
			return EvalStatement(node.Blocks[i], tempEnv)
		}
	}

	// the else condition
	if len(node.Blocks) > len(node.Conditions) {
		return EvalStatement(node.Blocks[len(node.Blocks)-1], tempEnv)
	}

	return &object.NullObj{}
}

func evalIdentifier(node ast.IdentifierExpression, env Environment) object.Object {
	// check if the identifier is a variable in the environment
	if obj := env.Get(node.TokenLiteral()); obj != nil {
		return obj
	}

	panic("unknown identifier: " + node.TokenLiteral())
}

func evalFunction(node ast.FunctionExpression) object.Object {
	args := []string{}
	for _, arg := range node.Args {
		args = append(args, arg.TokenLiteral())
	}

	return object.FunctionObj{
		Parameters: args,
		Body:       node.Body,
	}
}

func evalCall(node ast.CallExpression, env Environment) object.Object {
	obj := env.Get(node.Identifier.TokenLiteral())
	if obj == nil {
		panic("unknown function: " + node.Identifier.TokenLiteral())
	}

	switch funcObj := obj.(type) {
	case object.FunctionObj:
		funcEnv := NewEnclosedEnvironment(env)
		if len(node.Args) != len(funcObj.Parameters) {
			panic("incorrect count of arguments")
		}
		for i, arg := range node.Args {
			exp := EvalExpression(arg, env)
			funcEnv.Create(funcObj.Parameters[i], exp)
		}

		return EvalStatement(funcObj.Body, funcEnv)
	case *Builtin:
		args := []object.Object{}
		for _, arg := range node.Args {
			args = append(args, EvalExpression(arg, env))
		}
		return funcObj.Fn(env, args...)
	default:
		panic("unknown function type (What the shit?)")
	}

}

func evalArray(node ast.ArrayExpression, env Environment) object.Object {
	elems := []object.Object{}
	for _, e := range node.Elems {
		obj := EvalExpression(e, env)
		elems = append(elems, obj)
	}

	return &object.ArrayObj{Elements: elems}
}

func evalIndex(node ast.IndexExpression, env Environment) object.Object {
	exp := EvalExpression(node.Exp, env)
	index := EvalExpression(node.Index, env)

	switch indexObj := index.(type) {
	case *object.IntegerObj:
		return evalIntegerIndex(exp, indexObj)
	case *object.BooleanObj:
		return evalBoolIndex(exp, indexObj)
	case *object.StringObj:
		return evalStringIndex(exp, indexObj)
	default:
		panic("unsupported index data type")
	}
}

func evalIntegerIndex(exp object.Object, index *object.IntegerObj) object.Object {
	switch expObj := exp.(type) {
	case *object.ArrayObj:
		if index.Value >= int64(len(expObj.Elements)) {
			return &object.NullObj{}
		}
		return expObj.Elements[index.Value]
	case *object.StringObj:
		if index.Value >= int64(len(expObj.Value)) {
			return &object.NullObj{}
		}
		return &object.StringObj{Value: string(expObj.Value[index.Value])}
	case *object.HashObj:
		if pair, ok := expObj.Pairs[index.HashKey()]; ok {
			return pair.Value
		}
		return &object.NullObj{}
	default:
		panic("unindexable data type")
	}
}

func evalBoolIndex(exp object.Object, index *object.BooleanObj) object.Object {
	switch expObj := exp.(type) {
	case *object.HashObj:
		if pair, ok := expObj.Pairs[index.HashKey()]; ok {
			return pair.Value
		}
		return &object.NullObj{}
	default:
		panic("unindexable data type")
	}
}

func evalStringIndex(exp object.Object, index *object.StringObj) object.Object {
	switch expObj := exp.(type) {
	case *object.HashObj:
		if pair, ok := expObj.Pairs[index.HashKey()]; ok {
			return pair.Value
		}
		return &object.NullObj{}
	default:
		panic("unindexable data type")
	}
}

func evalHash(node ast.HashExpression, env Environment) object.Object {
	elems := map[object.HashKey]object.HashPair{}
	for _, kvp := range node.Elems {
		key := EvalExpression(kvp.Key, env)
		hashKey, ok := key.(object.Hashable)
		if !ok {
			panic("unhashable key type")
		}
		value := EvalExpression(kvp.Value, env)
		elems[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashObj{Pairs: elems}
}
