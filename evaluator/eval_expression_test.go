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
		{"5 + 5", 10},
		{"5 - 5", 0},
		{"5 * 5", 25},
		{"12 / 6", 2},
		{"5 % 2", 1},
	}

	for _, tt := range intTests {
		evaluated, errors := testEval(tt.input)
		if errors != nil {
			t.Fatalf("unexpected errors: %v", errors)
		}

		testIntegerObject(t, evaluated, tt.expected)
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

func testEval(input string) (object.Object, []error) {
	l := lexer.CreateLexer(input)
	p := parser.CreateParser(l)
	program, errors := p.ParseProgram()

	if len(errors) > 0 {
		return nil, errors
	}

	return Eval(program), nil
}
