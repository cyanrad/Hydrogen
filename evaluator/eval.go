package evaluator

import (
	"main/ast"
	"main/object"
)

func Eval(p ast.Program, env Environment) (object.Object, object.ErrorObj) {
	var lastStatement object.Object = object.NullObj{} // we return the value of the last statement in the program
	var err object.ErrorObj

	for _, statement := range p.Statements {
		lastStatement, err = EvalStatement(statement, env)
		if !err.Ok() {
			return object.NullObj{}, err
		}

		// halt if statement is a return statement
		if _, ok := statement.(ast.ReturnStatement); ok {
			return lastStatement, object.EmptyErrorObj()
		}
	}

	return lastStatement, object.EmptyErrorObj()
}
