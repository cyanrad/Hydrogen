package evaluator

import (
	"fmt"
	"main/object"
)

var builtins = map[string]*object.Builtin{
	"len":   {Fn: builtin_len},
	"push":  {Fn: builtin_push},
	"print": {Fn: builtin_print},
}

func builtin_len(args ...object.Object) object.Object {
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

func builtin_push(args ...object.Object) object.Object {
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

func builtin_print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Print(arg.Inspect() + " ")
	}
	fmt.Print("\n")
	return &object.NullObj{}
}
