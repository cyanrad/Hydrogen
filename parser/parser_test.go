package parser

import (
	"errors"
	"main/ast"
	"main/lexer"
	"main/token"
	"reflect"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `let x = 10; 
let y = 5;`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()
	if len(err) != 0 {
		t.Fatal(err)
	}

	statementCount := 2
	if len(prog.Statements) != statementCount {
		t.Fatalf("error - expected: %d statements - got: %d", statementCount, len(prog.Statements))
	}

	expectedProg := ast.Program{
		Statements: []ast.Statement{
			ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Identifier: ast.IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "x"},
				},
				Expression: ast.IntExpression{
					Token: token.Token{Type: token.INT, Literal: "10"},
				},
			},
			ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Identifier: ast.IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "y"},
				},
				Expression: ast.IntExpression{
					Token: token.Token{Type: token.INT, Literal: "5"},
				},
			},
		},
	}

	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}

func TestLetStatementErrors(t *testing.T) {
	input := `let x 5;
 let = 10;
 let 838383;
 let x = 10a;`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()

	errorCount := 4
	if len(err) != errorCount {
		t.Fatalf("error - expected: %d errors - got: %d", errorCount, len(err))
	}

	expectedErr := []error{
		errors.New("error - expected: = - got: INT"),
		errors.New("error - expected: IDENTIFIER - got: ="),
		errors.New("error - expected: IDENTIFIER - got: INT"),
		errors.New("error - expected: ; - got: IDENTIFIER"),
	}
	if ok := reflect.DeepEqual(err, expectedErr); !ok {
		t.Fatalf("expected: %v - got: %v", expectedErr, err)
	}

	statementCount := 0
	if len(prog.Statements) != statementCount {
		t.Fatalf("error - expected: %d statements - got: %d", statementCount, len(prog.Statements))
	}

	expectedProg := ast.Program{Statements: []ast.Statement{}}
	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}

func TestReturnStatement(t *testing.T) {
	input := `return 10;
return xyz;`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()
	if len(err) != 0 {
		t.Fatal(err)
	}

	statementCount := 2
	if len(prog.Statements) != statementCount {
		t.Fatalf("error - expected: %d statements - got: %d", statementCount, len(prog.Statements))
	}

	expectedProg := ast.Program{
		Statements: []ast.Statement{
			ast.ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				Expression: ast.IntExpression{
					Token: token.Token{Type: token.INT, Literal: "10"},
				},
			},
			ast.ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				Expression: ast.IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "xyz"},
				},
			},
		},
	}

	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}

func TestReturnStatementErrors(t *testing.T) {
	input := `return ;
return 12xyz;
return xyz();
return =`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()

	errorCount := 4
	if len(err) != errorCount {
		t.Fatalf("error - expected: %d errors - got: %d", errorCount, len(err))
	}

	expectedErr := []error{
		errors.New("error - expected: expression - got: ;"),
		errors.New("error - expected: ; - got: IDENTIFIER"),
		errors.New("error - expected: ; - got: ("),
		errors.New("error - expected: expression - got: ="),
	}
	if ok := reflect.DeepEqual(err, expectedErr); !ok {
		t.Fatalf("expected: %v - got: %v", expectedErr, err)
	}

	statementCount := 0
	if len(prog.Statements) != statementCount {
		t.Fatalf("error - expected: %d statements - got: %d", statementCount, len(prog.Statements))
	}

	expectedProg := ast.Program{Statements: []ast.Statement{}}
	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}

func TestExpressionStatements(t *testing.T) {
	input := `foobar;
5;`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()
	if len(err) != 0 {
		t.Fatal(err)
	}

	statementCount := 2
	if len(prog.Statements) != statementCount {
		t.Fatalf("error - expected: %d statements - got: %d", statementCount, len(prog.Statements))
	}

	expectedProg := ast.Program{
		Statements: []ast.Statement{
			ast.ExpressionStatement{
				Token: token.Token{Type: token.IDENTIFIER, Literal: "foobar"},
				Expression: ast.IdentifierExpression{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "foobar"},
				},
			},
			ast.ExpressionStatement{
				Token: token.Token{Type: token.INT, Literal: "5"},
				Expression: ast.IntExpression{
					Token: token.Token{Type: token.INT, Literal: "5"},
				},
			},
		},
	}

	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}
