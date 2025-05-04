package evaluator

import "main/object"

var builtins = map[string]*object.Builtin{
	"len":  {Fn: builtin_len, ArgCount: 1},
	"push": {Fn: builtin_push, ArgCount: 2},
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
	}

	panic("argument type to `len` not supported, got " + args[0].Type())
}

func builtin_push(args ...object.Object) object.Object {
	if len(args) != 2 {
		panic("wrong number of arguments in push()")
	}

	switch obj := args[0].(type) {
	case *object.ArrayObj:
		obj.Elements = append(obj.Elements, args[1])
		return obj
	}

	panic("argument type to `push` not supported, got " + args[0].Type())
}
