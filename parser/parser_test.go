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

func TestBasicExpressionStatements(t *testing.T) {
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

func TestPrefixExpressionStatements(t *testing.T) {
	input := `!5;
-15;
++foobar;
--x;`
	l := lexer.CreateLexer(input)
	p := CreateParser(l)

	prog, err := p.ParseProgram()
	if len(err) != 0 {
		t.Fatal(err)
	}

	statementCount := 4
	if len(prog.Statements) != statementCount {
		t.Fatalf("error - expected: %d statements - got: %d", statementCount, len(prog.Statements))
	}

	expectedProg := ast.Program{
		Statements: []ast.Statement{
			ast.ExpressionStatement{
				Token: token.Token{Type: token.BANG, Literal: "!"},
				Expression: ast.PrefixExpression{
					Token: token.Token{Type: token.BANG, Literal: "!"},
					Expression: ast.IntExpression{
						Token: token.Token{Type: token.INT, Literal: "5"},
					},
				},
			},
			ast.ExpressionStatement{
				Token: token.Token{Type: token.MINUS, Literal: "-"},
				Expression: ast.PrefixExpression{
					Token: token.Token{Type: token.MINUS, Literal: "-"},
					Expression: ast.IntExpression{
						Token: token.Token{Type: token.INT, Literal: "15"},
					},
				},
			},
			ast.ExpressionStatement{
				Token: token.Token{Type: token.INCREMENT, Literal: "++"},
				Expression: ast.PrefixExpression{
					Token: token.Token{Type: token.INCREMENT, Literal: "++"},
					Expression: ast.IdentifierExpression{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "foobar"},
					},
				},
			},
			ast.ExpressionStatement{
				Token: token.Token{Type: token.DECREMENT, Literal: "--"},
				Expression: ast.PrefixExpression{
					Token: token.Token{Type: token.DECREMENT, Literal: "--"},
					Expression: ast.IdentifierExpression{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "x"},
					},
				},
			},
		},
	}

	if ok := reflect.DeepEqual(prog, expectedProg); !ok {
		t.Fatalf("expected: %v - got: %v", expectedProg, prog)
	}
}

func TestInfixExpressionStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b;",
			"((-a) * b);",
		},
		{
			"!-1;",
			"(!(-1));",
		},
		{
			"a + b + c;",
			"((a + b) + c);",
		},
		{
			"a + b - c;",
			"((a + b) - c);",
		},
		{
			"a * b * c;",
			"((a * b) * c);",
		},
		{
			"a * b / c;",
			"((a * b) / c);",
		},
		{
			"a + b / c;",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4; -5 * 5;",
			"(3 + 4);\n((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		// {
		// 	"true",
		// 	"true",
		// },
		// {
		// 	"false",
		// 	"false",
		// },
		// {
		// 	"3 > 5 == false",
		// 	"((3 > 5) == false)",
		// },
		// {
		// 	"3 < 5 == true",
		// 	"((3 < 5) == true)",
		// },
		// {
		// 	"1 + (2 + 3) + 4",
		// 	"((1 + (2 + 3)) + 4)",
		// },
		// {
		// 	"(5 + 5) + 2",
		// 	"((5 + 5) + 2)",
		// },
		// {
		// 	"2 / (5 + 5)",
		// 	"(2 / (5 + 5))",
		// },
		// {
		// 	"-(5 + 5)",
		// 	"(-(5 + 5))",
		// },
		// {
		// 	"!(true == true)",
		// 	"(!(true == true))",
		// },
		// {
		// 	"a + add(b + c) + d",
		// 	"((a + add((b + c))) + d)",
		// },
		// {
		// 	"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		// 	"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		// },
		// {
		// 	"add(a + b + c * d / f + g)",
		// 	"add((((a + b) + ((c * d) / f)) + g))",
		// },
		// {
		// 	"a * [1, 2, 3, 4][b * c] * d",
		// 	"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		// },
		// {
		// 	"add(a * b[2], b[1], 2 * [1,2][1])",
		// 	"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		// },
	}

	for _, tt := range tests {
		l := lexer.CreateLexer(tt.input)
		p := CreateParser(l)
		prog, errs := p.ParseProgram()
		if len(errs) != 0 {
			t.Fatal(errs)
		}

		actual := prog.String()
		if actual != tt.expected {
			t.Errorf("error - expected: %q, got: %q", tt.expected, actual)
		}
	}
}
