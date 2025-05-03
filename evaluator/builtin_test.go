package evaluator

import (
	"testing"
)

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
	}
	for _, tt := range tests {
		evaluated, _ := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		}
	}
}
