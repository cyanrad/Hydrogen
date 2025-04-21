package evaluator

import (
	"main/object"
	"testing"

	"main/lexer"
	"main/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}
	for _, tt := range tests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalPrefixExpression(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		{"-5", -5},
		{"--5", 4},
		{"++5", 6},
	}
	for _, tt := range intTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}

	boolTests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
	}
	for _, tt := range boolTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalInfixExpression(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		// Arithmetic Operations
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"12 / 6", 2},
		{"5 % 2", 1},
		{"5 + 5 - 5", 5},
		{"5 * 2 + 10", 20},
		{"10 - 5 * 2", 0},

		// Bitwise Operations
		{"5 & 3", 1},
		{"5 | 3", 7},
	}

	for _, tt := range intTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}

	boolTests := []struct {
		input    string
		expected bool
	}{
		// Boolean Operations
		{"true == true", true},
		{"true == false", false},
		{"true != false", true},
		{"false == false", true},

		// Comparison Operations
		{"5 < 10", true},
		{"10 > 5", true},
		{"5 <= 5", true},
		{"5 >= 5", true},
		{"5 < 5", false},
		{"5 > 5", false},
		{"5 <= 4", false},
		{"5 >= 6", false},
		{"5 == 5", true},
		{"5 != 5", false},
		{"5 == 6", false},
		{"5 != 6", true},

		// Logical Operations
		{"true && true", true},
		{"true && false", false},
		{"false || true", true},
		{"false || false", false},

		// Combined Operations
		{"5 + 5 == 10", true},
		{"5 + 5 != 10", false},
		{"5 * 2 < 10", false},

		// With Parentheses
		{"(5 + 5) == 10", true},
		{"(5 + 5) != 10", false},
		{"(5 * 2) < (10 - 1)", false},
		{"(5 * 2 > 10) || (5 < 8)", true},
	}

	for _, tt := range boolTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	intTests := []struct {
		input    string
		expected int64
	}{
		{"if (true) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},

		// nested if
		{"if (2 > 1) { if (true) { 10 } else { 20 } } else { 30 }", 10},

		// else if
		{"if (1 > 2) { 10 } else if (true) { 20 } else { 30 }", 20},
	}

	for _, tt := range intTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testIntegerObject(t, evaluated, tt.expected)
	}

	nullTests := []struct {
		input string
	}{
		{"if (false) { 10 }"},
		{"if (0) { 10 }"},
	}

	for _, tt := range nullTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testNullObject(t, evaluated)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.BooleanObj)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.IntegerObj)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.NullObj)
	if !ok {
		t.Errorf("object is not Null. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string) (object.Object, []error) {
	l := lexer.CreateLexer(input)
	p := parser.CreateParser(l)
	program, errors := p.ParseProgram()

	if len(errors) > 0 {
		return nil, errors
	}

	return Eval(program), nil
}
