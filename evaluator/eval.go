package evaluator

import (
	"main/ast"
	"main/object"
)

func Eval(p ast.Program, env Environment) object.Object {
	var lastStatement object.Object = object.NullObj{} // we return the value of the last statement in the program

	for _, statement := range p.Statements {
		lastStatement = EvalStatement(statement, env)
	}

	return lastStatement
}
