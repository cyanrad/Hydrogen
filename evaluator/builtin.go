package evaluator

import (
	"fmt"
	"main/object"
	"os"
)

type BuiltinFunction func(env Environment, args ...object.Object) (object.Object, object.ErrorObj)

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
		"exit":   {Fn: builtin_exit},
		"rest":   {Fn: builtin_rest},
		"filter": {Fn: builtin_filter},
		"map":    {Fn: builtin_map},
		"reduce": {Fn: builtin_reduce},
	}
}

func builtin_len(_ Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) != 1 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("len() requires exactly one argument, got %d", len(args)),
		)
	}

	switch obj := args[0].(type) {
	case *object.StringObj:
		return &object.IntegerObj{Value: int64(len(obj.Value))}, object.EmptyErrorObj()
	case *object.ArrayObj:
		return &object.IntegerObj{Value: int64(len(obj.Elements))}, object.EmptyErrorObj()
	case *object.HashObj:
		return &object.IntegerObj{Value: int64(len(obj.Pairs))}, object.EmptyErrorObj()
	}

	return &object.NullObj{}, object.NewErrorObj(
		"argument type to len() not supported, got " + string(args[0].Type()),
	)
}

func builtin_push(_ Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) != 3 && len(args) != 2 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("push() requires 2 or 3 arguments, got %d", len(args)),
		)
	}

	switch obj := args[0].(type) {
	case *object.ArrayObj:
		if len(args) != 2 {
			return &object.NullObj{}, object.NewErrorObj(
				fmt.Sprintf("push() to an array requires exactly 2 arguments, got %d", len(args)),
			)
		}
		obj.Elements = append(obj.Elements, args[1])
		return obj, object.EmptyErrorObj()
	case *object.HashObj:
		if len(args) != 3 {
			return &object.NullObj{}, object.NewErrorObj(
				fmt.Sprintf("push() to a hash requires exactly 3 arguments, got %d", len(args)),
			)
		}
		key, ok := args[1].(object.Hashable)
		if !ok {
			return &object.NullObj{}, object.NewErrorObj(
				"key to push() to a hash must be hashable, got " + string(args[1].Type()),
			)
		}
		obj.Pairs[key.HashKey()] = object.HashPair{Key: args[1], Value: args[2]}
		return obj, object.EmptyErrorObj()
	}

	return &object.NullObj{}, object.NewErrorObj(
		"argument type to push() not supported, got " + string(args[0].Type()),
	)
}

func builtin_print(_ Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	for _, arg := range args {
		fmt.Print(arg.Inspect() + " ")
	}
	fmt.Print("\n")
	return &object.NullObj{}, object.EmptyErrorObj()
}

func builtin_exit(_ Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) > 1 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("exit() takes at most one argument, got %d", len(args)),
		)
	}

	if len(args) == 0 {
		fmt.Println("Exiting with code 0")
		os.Exit(0)
	} else if intObj, ok := args[0].(*object.IntegerObj); ok {
		fmt.Printf("Exiting with code %d\n", intObj.Value)
		os.Exit(int(intObj.Value))
	} else {
		panic("argument to exit() must be an integer")
	}

	return &object.NullObj{}, object.EmptyErrorObj()
}

func builtin_rest(env Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) != 1 && len(args) != 2 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("rest() requires 1 or 2 arguments, got %d", len(args)),
		)
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"first argument to rest() must be an array, got " + string(args[0].Type()),
		)
	}

	start := 1
	if len(args) == 2 {
		if intObj, ok := args[1].(*object.IntegerObj); ok {
			start = int(intObj.Value)
		} else {
			return &object.NullObj{}, object.NewErrorObj(
				"second argument to rest() must be an integer, got " + string(args[1].Type()),
			)
		}
	}

	if len(arr.Elements) == 0 {
		return &object.ArrayObj{Elements: []object.Object{}}, object.EmptyErrorObj()
	}

	result := make([]object.Object, len(arr.Elements)-start)
	copy(result, arr.Elements[start:])
	return &object.ArrayObj{Elements: result}, object.EmptyErrorObj()
}

func builtin_filter(env Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) != 2 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("filter() requires exactly 2 arguments, got %d", len(args)),
		)
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"first argument to filter() must be an array, got " + string(args[0].Type()),
		)
	}

	fn, ok := args[1].(object.FunctionObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"second argument to filter() must be a function, got " + string(args[1].Type()),
		)
	} else if len(fn.Parameters) != 1 {
		return &object.NullObj{}, object.NewErrorObj(
			"filter function must take exactly one argument, got " + fmt.Sprintf("%d", len(fn.Parameters)),
		)
	}

	result := []object.Object{}
	for _, elem := range arr.Elements {
		funcEnv := NewEnclosedEnvironment(env)
		funcEnv.Create(fn.Parameters[0], elem)
		bool, err := EvalStatement(fn.Body, funcEnv)
		if !err.Ok() {
			return &object.NullObj{}, object.NewErrorObj("error evaluating filter function", err)
		}
		if boolObj, ok := bool.(*object.BooleanObj); ok && boolObj.Value {
			result = append(result, elem)
		}
	}

	return &object.ArrayObj{Elements: result}, object.EmptyErrorObj()
}

func builtin_map(env Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) != 2 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("map() requires exactly 2 arguments, got %d", len(args)),
		)
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"first argument to map() must be an array, got " + string(args[0].Type()),
		)
	}

	fn, ok := args[1].(object.FunctionObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"second argument to map() must be a function, got " + string(args[1].Type()),
		)
	} else if len(fn.Parameters) != 1 {
		return &object.NullObj{}, object.NewErrorObj(
			"map function must take exactly one argument, got " + fmt.Sprintf("%d", len(fn.Parameters)),
		)
	}

	result := []object.Object{}
	for _, elem := range arr.Elements {
		funcEnv := NewEnclosedEnvironment(env)
		funcEnv.Create(fn.Parameters[0], elem)
		mapped, err := EvalStatement(fn.Body, funcEnv)
		if !err.Ok() {
			return &object.NullObj{}, object.NewErrorObj("error evaluating map function", err)
		}

		result = append(result, mapped)
	}

	return &object.ArrayObj{Elements: result}, object.EmptyErrorObj()
}

func builtin_reduce(env Environment, args ...object.Object) (object.Object, object.ErrorObj) {
	if len(args) != 3 {
		return &object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("reduce() requires exactly 3 arguments, got %d", len(args)),
		)
	}

	arr, ok := args[0].(*object.ArrayObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"first argument to reduce() must be an array, got " + string(args[0].Type()),
		)
	}

	prev, ok := args[1].(object.Object)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"second argument to reduce() must be an initial value, got " + string(args[1].Type()),
		)
	}

	fn, ok := args[2].(object.FunctionObj)
	if !ok {
		return &object.NullObj{}, object.NewErrorObj(
			"third argument to reduce() must be a function, got " + string(args[2].Type()),
		)
	} else if len(fn.Parameters) != 2 {
		return &object.NullObj{}, object.NewErrorObj(
			"reduce function must take exactly two arguments, got " + fmt.Sprintf("%d", len(fn.Parameters)),
		)
	}

	for _, elem := range arr.Elements {
		funcEnv := NewEnclosedEnvironment(env)
		funcEnv.Create(fn.Parameters[0], prev)
		funcEnv.Create(fn.Parameters[1], elem)
		mapped, err := EvalStatement(fn.Body, funcEnv)
		if !err.Ok() {
			return &object.NullObj{}, object.NewErrorObj("error evaluating reduce function", err)
		}

		prev = mapped
	}

	return prev, object.EmptyErrorObj()
}
