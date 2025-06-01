package evaluator

import (
	"main/ast"
	"main/object"
	"strconv"
)

func EvalExpression(n ast.Expression, env Environment) (object.Object, object.ErrorObj) {
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

func evalInteger(node ast.IntExpression) (object.Object, object.ErrorObj) {
	value, err := strconv.ParseInt(node.TokenLiteral(), 0, 64) // highly unlikely to fail
	if err != nil {
		return object.NullObj{}, object.NewErrorObj("failed to parse integer: " + err.Error())
	}

	return &object.IntegerObj{Value: value}, object.EmptyErrorObj()
}

func evalBoolean(node ast.BooleanExpression) (object.Object, object.ErrorObj) {
	if node.TokenLiteral() == "true" {
		return &object.BooleanObj{Value: true}, object.EmptyErrorObj()
	} else if node.TokenLiteral() == "false" {
		return &object.BooleanObj{Value: false}, object.EmptyErrorObj()
	}
	return object.NullObj{}, object.NewErrorObj("unknown boolean value: " + node.TokenLiteral())
}

func evalString(node ast.StringExpression) (object.Object, object.ErrorObj) {
	return &object.StringObj{Value: node.TokenLiteral()}, object.EmptyErrorObj()
}

func evalPrefix(node ast.PrefixExpression, env Environment) (object.Object, object.ErrorObj) {
	exp, err := EvalExpression(node.Expression, env)
	if !err.Ok() {
		return object.NullObj{}, object.NewErrorObj("failed to evaluate prefix expression", err)
	}

	if intexp, ok := exp.(*object.IntegerObj); ok {
		switch node.TokenLiteral() {
		case "-":
			intexp.Value = -intexp.Value
		case "++":
			intexp.Value += 1
		case "--":
			intexp.Value -= 1
		default:
			return object.NullObj{}, object.NewErrorObj("unknown int prefix operator: " + node.TokenLiteral())
		}

		return intexp, object.EmptyErrorObj()
	} else if boolExp, ok := exp.(*object.BooleanObj); ok {
		switch node.TokenLiteral() {
		case "!":
			boolExp.Value = !boolExp.Value
		default:
			return object.NullObj{}, object.NewErrorObj("unknown bool prefix operator: " + node.TokenLiteral())
		}

		return boolExp, object.EmptyErrorObj()
	} else {
		return object.NullObj{}, object.NewErrorObj("unknown prefix expression type: " + node.TokenLiteral())
	}
}

func evalInfix(node ast.InfixExpression, env Environment) (object.Object, object.ErrorObj) {
	left, err := EvalExpression(node.Left, env)
	if !err.Ok() {
		return object.NullObj{}, object.NewErrorObj("failed to evaluate left expression", err)
	}

	right, err := EvalExpression(node.Right, env)
	if !err.Ok() {
		return object.NullObj{}, object.NewErrorObj("failed to evaluate infix right expression", err)
	}

	leftInt, leftOk := left.(*object.IntegerObj)
	rightInt, rightOk := right.(*object.IntegerObj)
	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "+":
			return &object.IntegerObj{Value: leftInt.Value + rightInt.Value}, object.EmptyErrorObj()
		case "-":
			return &object.IntegerObj{Value: leftInt.Value - rightInt.Value}, object.EmptyErrorObj()
		case "*":
			return &object.IntegerObj{Value: leftInt.Value * rightInt.Value}, object.EmptyErrorObj()
		case "/":
			return &object.IntegerObj{Value: leftInt.Value / rightInt.Value}, object.EmptyErrorObj()
		case "%":
			return &object.IntegerObj{Value: leftInt.Value % rightInt.Value}, object.EmptyErrorObj()
		case "<":
			return &object.BooleanObj{Value: leftInt.Value < rightInt.Value}, object.EmptyErrorObj()
		case "<=":
			return &object.BooleanObj{Value: leftInt.Value <= rightInt.Value}, object.EmptyErrorObj()
		case ">":
			return &object.BooleanObj{Value: leftInt.Value > rightInt.Value}, object.EmptyErrorObj()
		case ">=":
			return &object.BooleanObj{Value: leftInt.Value >= rightInt.Value}, object.EmptyErrorObj()
		case "==":
			return &object.BooleanObj{Value: leftInt.Value == rightInt.Value}, object.EmptyErrorObj()
		case "!=":
			return &object.BooleanObj{Value: leftInt.Value != rightInt.Value}, object.EmptyErrorObj()
		case "&":
			return &object.IntegerObj{Value: leftInt.Value & rightInt.Value}, object.EmptyErrorObj()
		case "|":
			return &object.IntegerObj{Value: leftInt.Value | rightInt.Value}, object.EmptyErrorObj()
		default:
			return &object.NullObj{}, object.NewErrorObj("unknown int infix operator: " + node.TokenLiteral())
		}
	}

	leftBool, leftOk := left.(*object.BooleanObj)
	rightBool, rightOk := right.(*object.BooleanObj)
	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "==":
			return &object.BooleanObj{Value: leftBool.Value == rightBool.Value}, object.EmptyErrorObj()
		case "!=":
			return &object.BooleanObj{Value: leftBool.Value != rightBool.Value}, object.EmptyErrorObj()
		case "&&":
			return &object.BooleanObj{Value: leftBool.Value && rightBool.Value}, object.EmptyErrorObj()
		case "||":
			return &object.BooleanObj{Value: leftBool.Value || rightBool.Value}, object.EmptyErrorObj()
		default:
			return &object.NullObj{}, object.NewErrorObj("unknown bool infix operator: " + node.TokenLiteral())
		}
	}

	leftStr, leftOk := left.(*object.StringObj)
	rightStr, rightOk := right.(*object.StringObj)
	if leftOk && rightOk {
		switch node.TokenLiteral() {
		case "==":
			return &object.BooleanObj{Value: leftStr.Value == rightStr.Value}, object.EmptyErrorObj()
		case "!=":
			return &object.BooleanObj{Value: leftStr.Value != rightStr.Value}, object.EmptyErrorObj()
		case "+":
			return &object.StringObj{Value: leftStr.Value + rightStr.Value}, object.EmptyErrorObj()
		default:
			return &object.NullObj{}, object.NewErrorObj("unknown string infix operator: " + node.TokenLiteral())
		}
	}

	lts := string(left.Type())
	rts := string(right.Type())
	return &object.NullObj{}, object.NewErrorObj("unknown infix expression types: " +
		node.TokenLiteral() + " between " + lts + " and " + rts)
}

func evalIf(node ast.IfExpression, env Environment) (object.Object, object.ErrorObj) {
	tempEnv := NewEnclosedEnvironment(env)

	// looping over the conditions, return the block of the first condition that evaluates to true
	for i, condition := range node.Conditions {
		cond, err := EvalExpression(condition, env)
		if !err.Ok() {
			return object.NullObj{}, object.NewErrorObj("failed to evaluate if condition", err)
		}

		if boolCond, ok := cond.(*object.BooleanObj); ok && boolCond.Value {
			body, err := EvalStatement(node.Blocks[i], tempEnv)
			if !err.Ok() {
				return object.NullObj{}, object.NewErrorObj("failed to evaluate if block", err)
			}
			return body, object.EmptyErrorObj()
		}
	}

	// the else condition
	if len(node.Blocks) > len(node.Conditions) {
		body, err := EvalStatement(node.Blocks[len(node.Blocks)-1], tempEnv)
		if !err.Ok() {
			return object.NullObj{}, object.NewErrorObj("failed to evaluate else block", err)
		}
		return body, object.EmptyErrorObj()
	}

	return &object.NullObj{}, object.EmptyErrorObj()
}

func evalIdentifier(node ast.IdentifierExpression, env Environment) (object.Object, object.ErrorObj) {
	// check if the identifier is a variable in the environment
	if obj := env.Get(node.TokenLiteral()); obj != nil {
		return obj, object.EmptyErrorObj()
	}

	return &object.NullObj{}, object.NewErrorObj("unknown identifier: " + node.TokenLiteral())
}

func evalFunction(node ast.FunctionExpression) (object.Object, object.ErrorObj) {
	args := []string{}
	for _, arg := range node.Args {
		args = append(args, arg.TokenLiteral())
	}

	return object.FunctionObj{
		Parameters: args,
		Body:       node.Body,
	}, object.EmptyErrorObj()
}

func evalCall(node ast.CallExpression, env Environment) (object.Object, object.ErrorObj) {
	obj := env.Get(node.Identifier.TokenLiteral())
	if obj == nil {
		return &object.NullObj{}, object.NewErrorObj("unknown function: " + node.Identifier.TokenLiteral())
	}

	switch funcObj := obj.(type) {
	case object.FunctionObj:
		funcEnv := NewEnclosedEnvironment(env)

		if len(node.Args) != len(funcObj.Parameters) {
			return &object.NullObj{}, object.NewErrorObj(
				"wrong number of arguments: expected " + strconv.Itoa(len(funcObj.Parameters)) +
					", got " + strconv.Itoa(len(node.Args)),
			)
		}

		for i, arg := range node.Args {
			exp, err := EvalExpression(arg, env)
			if !err.Ok() {
				return &object.NullObj{}, object.NewErrorObj(
					"failed to evaluate argument for function '"+node.Identifier.TokenLiteral()+"'", err,
				)
			}

			funcEnv.Create(funcObj.Parameters[i], exp)
		}

		return EvalStatement(funcObj.Body, funcEnv)
	case *Builtin:
		args := []object.Object{}
		for _, arg := range node.Args {
			val, err := EvalExpression(arg, env)
			if !err.Ok() {
				return &object.NullObj{}, object.NewErrorObj(
					"failed to evaluate argument for builtin function '"+node.Identifier.TokenLiteral()+"'", err,
				)
			}

			args = append(args, val)
		}

		return funcObj.Fn(env, args...)
	default:
		return &object.NullObj{}, object.NewErrorObj("unknown function type (what the shit?)")
	}

}

func evalArray(node ast.ArrayExpression, env Environment) (object.Object, object.ErrorObj) {
	elems := []object.Object{}
	for _, e := range node.Elems {
		obj, err := EvalExpression(e, env)
		if !err.Ok() {
			return &object.NullObj{}, object.NewErrorObj("failed to evaluate array element", err)
		}
		elems = append(elems, obj)
	}

	return &object.ArrayObj{Elements: elems}, object.EmptyErrorObj()
}

func evalIndex(node ast.IndexExpression, env Environment) (object.Object, object.ErrorObj) {
	exp, err := EvalExpression(node.Exp, env)
	if !err.Ok() {
		return &object.NullObj{}, object.NewErrorObj("failed to evaluate container identifier", err)
	}

	index, err := EvalExpression(node.Index, env)
	if !err.Ok() {
		return &object.NullObj{}, object.NewErrorObj("failed to evaluate container index", err)
	}

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

func evalIntegerIndex(exp object.Object, index *object.IntegerObj) (object.Object, object.ErrorObj) {
	switch expObj := exp.(type) {
	case *object.ArrayObj:
		if index.Value >= int64(len(expObj.Elements)) {
			return &object.NullObj{}, object.NewErrorObj(
				"index out of bounds, attempted to access " + index.Inspect() +
					" in array of length " + strconv.Itoa(len(expObj.Elements)),
			)
		}
		return expObj.Elements[index.Value], object.EmptyErrorObj()
	case *object.StringObj:
		if index.Value >= int64(len(expObj.Value)) {
			return &object.NullObj{}, object.NewErrorObj(
				"index out of bounds, attempted to access " + index.Inspect() +
					" in a string of length " + strconv.Itoa(len(expObj.Value)),
			)
		}
		return &object.StringObj{Value: string(expObj.Value[index.Value])}, object.EmptyErrorObj()
	case *object.HashObj:
		if pair, ok := expObj.Pairs[index.HashKey()]; ok {
			return pair.Value, object.EmptyErrorObj()
		}
		return &object.NullObj{}, object.NewErrorObj("key " + index.Inspect() + " not found in hash")

	default:
		return &object.NullObj{}, object.NewErrorObj("unindexable data type using int: " + string(exp.Type()))
	}
}

func evalBoolIndex(exp object.Object, index *object.BooleanObj) (object.Object, object.ErrorObj) {
	switch expObj := exp.(type) {
	case *object.HashObj:
		if pair, ok := expObj.Pairs[index.HashKey()]; ok {
			return pair.Value, object.EmptyErrorObj()
		}
		return &object.NullObj{}, object.NewErrorObj("key " + index.Inspect() + " not found in hash")
	default:
		return &object.NullObj{}, object.NewErrorObj("unindexable data type using bool: " + string(exp.Type()))
	}
}

func evalStringIndex(exp object.Object, index *object.StringObj) (object.Object, object.ErrorObj) {
	switch expObj := exp.(type) {
	case *object.HashObj:
		if pair, ok := expObj.Pairs[index.HashKey()]; ok {
			return pair.Value, object.EmptyErrorObj()
		}
		return &object.NullObj{}, object.NewErrorObj("key '" + index.Inspect() + "' not found in hash")
	default:
		return &object.NullObj{}, object.NewErrorObj("unindexable data type using string: " + string(exp.Type()))
	}
}

func evalHash(node ast.HashExpression, env Environment) (object.Object, object.ErrorObj) {
	elems := map[object.HashKey]object.HashPair{}
	for _, kvp := range node.Elems {
		key, err := EvalExpression(kvp.Key, env)
		if !err.Ok() {
			return &object.NullObj{}, object.NewErrorObj("failed to evaluate hash key", err)
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			panic("unhashable key type")
		}

		value, err := EvalExpression(kvp.Value, env)
		if !err.Ok() {
			return &object.NullObj{}, object.NewErrorObj("failed to evaluate hash value", err)
		}

		elems[hashKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashObj{Pairs: elems}, object.EmptyErrorObj()
}
