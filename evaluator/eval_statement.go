package evaluator

import (
	"main/ast"
	"main/object"
)

func EvalStatement(s ast.Statement, env Environment) object.Object {
	switch stmt := s.(type) {
	case ast.BlockStatement:
		return evalBlockStatement(stmt, env)
	case ast.ExpressionStatement:
		return EvalExpression(stmt.Expression, env)
	case ast.LetStatement:
		return evalLetStatement(stmt, env)
	case ast.ReturnStatement:
		return evalReturnStatement(stmt)
	default:
		panic("unknown statement type")
	}
}

func evalBlockStatement(block ast.BlockStatement, env Environment) object.Object {
	prog := ast.Program{Statements: block.Statements}
	return Eval(prog, env)
}

func evalLetStatement(stmt ast.LetStatement, env Environment) object.Object {
	// Check if the variable already exists in the environment
	ident := stmt.Identifier.TokenLiteral()
	existingVar := env.GetVariable(ident)
	if existingVar != nil {
		panic("variable already exists")
	}

	env.CreateVariable(ident, EvalExpression(stmt.Expression, env))
	return object.NullObj{}
}

func evalReturnStatement(stmt ast.ReturnStatement) object.Object {
	return object.NullObj{}
}
