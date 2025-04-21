package evaluator

import (
	"main/ast"
	"main/object"
)

func EvalStatement(s ast.Statement) object.Object {
	switch stmt := s.(type) {
	case ast.BlockStatement:
		return evalBlockStatement(stmt)
	case ast.ExpressionStatement:
		return EvalExpression(stmt.Expression)
	case ast.LetStatement:
		return evalLetStatement(stmt)
	case ast.ReturnStatement:
		return evalReturnStatement(stmt)
	default:
		panic("unknown statement type")
	}
}

func evalBlockStatement(block ast.BlockStatement) object.Object {
	prog := ast.Program{Statements: block.Statements}
	return Eval(prog)
}

func evalLetStatement(stmt ast.LetStatement) object.Object {
	return nil
}

func evalReturnStatement(stmt ast.ReturnStatement) object.Object {
	return nil
}
