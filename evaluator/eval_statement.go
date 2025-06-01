package evaluator

import (
	"fmt"
	"main/ast"
	"main/object"
)

func EvalStatement(s ast.Statement, env Environment) (object.Object, object.ErrorObj) {
	switch stmt := s.(type) {
	case ast.BlockStatement:
		return evalBlockStatement(stmt, env)
	case ast.ExpressionStatement:
		return EvalExpression(stmt.Expression, env)
	case ast.LetStatement:
		return evalLetStatement(stmt, env)
	case ast.ReturnStatement:
		return evalReturnStatement(stmt, env)
	default:
		return object.NullObj{}, object.NewErrorObj(fmt.Sprintf("unknown statement type: %T", stmt))

	}
}

func evalBlockStatement(block ast.BlockStatement, env Environment) (object.Object, object.ErrorObj) {
	prog := ast.Program{Statements: block.Statements}
	val, err := Eval(prog, env)
	if !err.Ok() {
		return object.NullObj{}, object.NewErrorObj("error evaluating block statement", err)
	}

	return val, object.EmptyErrorObj()
}

func evalLetStatement(stmt ast.LetStatement, env Environment) (object.Object, object.ErrorObj) {
	// Check if the variable already exists in the environment
	ident := stmt.Identifier.TokenLiteral()
	existingVar := env.Get(ident)
	if existingVar != nil {
		return object.NullObj{}, object.NewErrorObj(fmt.Sprintf("variable '%s' already exists", ident))
	}

	val, err := EvalExpression(stmt.Expression, env)
	if !err.Ok() {
		return object.NullObj{}, object.NewErrorObj(
			fmt.Sprintf("error evaluating expression for variable '%s'", ident), err,
		)
	}
	env.Create(ident, val)
	return object.NullObj{}, object.EmptyErrorObj()
}

func evalReturnStatement(stmt ast.ReturnStatement, env Environment) (object.Object, object.ErrorObj) {
	if stmt.Expression == nil {
		return object.NullObj{}, object.EmptyErrorObj()
	}

	val, err := EvalExpression(stmt.Expression, env)
	if !err.Ok() {
		return object.NullObj{}, object.NewErrorObj("error evaluating return expression", err)
	}
	return val, object.EmptyErrorObj()
}
