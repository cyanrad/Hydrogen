package evaluator

import "main/object"

var builtins = map[string]*object.Builtin{
	"len": {Fn: builtin_len, ArgCount: 1},
}

func builtin_len(args ...object.Object) object.Object {
	if len(args) != 1 {
		panic("wrong number of arguments in len()")
	}
	if strObj, ok := args[0].(*object.StringObj); ok {
		return &object.IntegerObj{Value: int64(len(strObj.Value))}
	}

	panic("argument type to `len` not supported, got " + args[0].Type())
}
