package evaluator

import (
	"fmt"
	"main/object"
)

type BuiltinFunction func(env Environment, args ...object.Object) object.Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() object.ObjectType { return object.BUILTIN_OBJ }
func (b *Builtin) Inspect() string         { return "builtin function" }

var builtins map[string]*Builtin

func InitBuiltins() {
	builtins = map[string]*Builtin{
		"len":    {Fn: builtin_len},
		"push":   {Fn: builtin_push},
		"print":  {Fn: builtin_print},
		"filter": {Fn: builtin_filter},
		"map":    {Fn: builtin_map},
		"reduce": {Fn: builtin_reduce},
	}
}

func builtin_len(_ Environment, args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("wrong number of arguments in len()")
	}

	switch obj := args[0].(type) {
	case *object.StringObj:
		return &object.IntegerObj{Value: int64(len(obj.Value))}
	case *object.ArrayObj:
		return &object.IntegerObj{Value: int64(len(obj.Elements))}
	case *object.HashObj:
		return &object.IntegerObj{Value: int64(len(obj.Pairs))}
	}

	panic("argument type to `len` not supported, got " + args[0].Type())
}

func builtin_push(_ Environment, args ...object.Object) object.Object {
	if len(args) != 3 && len(args) != 2 {
		panic("wrong number of arguments in push()")
	}

	switch obj := args[0].(type) {
	case *object.ArrayObj:
		if len(args) != 2 {
			panic("wrong number of arguments in push() to an array")
		}
		obj.Elements = append(obj.Elements, args[1])
		return obj
	case *object.HashObj:
		if len(args) != 3 {
			panic("wrong number of arguments in push() to a hash")
		}
		key, ok := args[1].(object.Hashable)
		if !ok {
			panic("key to push() to a hash must be hashable")
		}
		obj.Pairs[key.HashKey()] = object.HashPair{Key: args[1], Value: args[2]}
		return obj
	}

	panic("argument type to `push` not supported, got " + args[0].Type())
}

func builtin_print(_ Environment, args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect() + " ")
	}
	fmt.Print("\n")
	return &object.NullObj{}
}

func builtin_filter(env Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		panic("wrong number of arguments in filter()")
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		panic("first argument to filter() must be an array")
	}

	fn, ok := args[1].(object.FunctionObj)
	if !ok {
		panic("second argument to filter() must be a function")
	} else if len(fn.Parameters) != 1 {
		panic("filter function must take exactly one argument")
	}

	result := []object.Object{}
	for _, elem := range arr.Elements {
		funcEnv := NewEnclosedEnvironment(env)
		funcEnv.Create(fn.Parameters[0], elem)
		bool := EvalStatement(fn.Body, funcEnv)
		if boolObj, ok := bool.(*object.BooleanObj); ok && boolObj.Value {
			result = append(result, elem)
		}
	}

	return &object.ArrayObj{Elements: result}
}

func builtin_map(env Environment, args ...object.Object) object.Object {
	if len(args) != 2 {
		panic("wrong number of arguments in map()")
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		panic("first argument to map() must be an array")
	}

	fn, ok := args[1].(object.FunctionObj)
	if !ok {
		panic("second argument to map() must be a function")
	} else if len(fn.Parameters) != 1 {
		panic("map function must take exactly one argument")
	}

	result := []object.Object{}
	for _, elem := range arr.Elements {
		funcEnv := NewEnclosedEnvironment(env)
		funcEnv.Create(fn.Parameters[0], elem)
		mapped := EvalStatement(fn.Body, funcEnv)
		result = append(result, mapped)
	}

	return &object.ArrayObj{Elements: result}
}

func builtin_reduce(env Environment, args ...object.Object) object.Object {
	if len(args) != 3 {
		panic("wrong number of arguments in reduce()")
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		panic("first argument to reduce() must be an array")
	}

	prev, ok := args[1].(object.Object)
	if !ok {
		panic("second argument to reduce() must be an object")
	}

	fn, ok := args[2].(object.FunctionObj)
	if !ok {
		panic("third argument to reduce() must be a function")
	} else if len(fn.Parameters) != 2 {
		panic("reduce function must take exactly one argument")
	}

	for _, elem := range arr.Elements {
		funcEnv := NewEnclosedEnvironment(env)
		funcEnv.Create(fn.Parameters[0], prev)
		funcEnv.Create(fn.Parameters[1], elem)
		mapped := EvalStatement(fn.Body, funcEnv)
		prev = mapped
	}

	return prev
}
