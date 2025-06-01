package object

import (
	"hash/fnv"
	"main/ast"
	"strconv"
	"strings"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hashable interface {
	Object
	HashKey() HashKey
}

type IntegerObj struct {
	Value int64
}

func (i IntegerObj) Type() ObjectType { return INT_OBJ }
func (i IntegerObj) Inspect() string  { return strconv.Itoa(int(i.Value)) }
func (i *IntegerObj) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type BooleanObj struct {
	Value bool
}

func (b BooleanObj) Type() ObjectType { return BOOLEAN_OBJ }
func (b BooleanObj) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (b *BooleanObj) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type StringObj struct {
	Value string
}

func (s StringObj) Type() ObjectType { return STRING_OBJ }
func (s StringObj) Inspect() string  { return s.Value }
func (s *StringObj) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type NullObj struct{}

func (n NullObj) Type() ObjectType { return NULL_OBJ }
func (n NullObj) Inspect() string  { return "null" }

type FunctionObj struct {
	Parameters []string
	Body       ast.BlockStatement
}

func (f FunctionObj) Type() ObjectType { return FUNCTION_OBJ }
func (f FunctionObj) Inspect() string {
	return "fn(" + strings.Join(f.Parameters, ", ") + ")"
}

type ArrayObj struct {
	Elements []Object
}

func (a ArrayObj) Type() ObjectType { return ARRAY_OBJ }
func (a ArrayObj) Inspect() string {
	var elements []string
	for _, element := range a.Elements {
		elements = append(elements, element.Inspect())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

type HashObj struct {
	Pairs map[HashKey]HashPair
}

func (h HashObj) Type() ObjectType { return HASH_OBJ }
func (h HashObj) Inspect() string {
	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

type ErrorObj struct {
	Message   string
	SubErrors []ErrorObj
}

func (e ErrorObj) Type() ObjectType { return ERROR_OBJ }
func (e ErrorObj) Inspect() string {
	output := "Error: " + e.Message + "\n"
	for _, subError := range e.SubErrors {
		output += "\t" + subError.Inspect() + "\n"
	}
	return output
}
func (e ErrorObj) Ok() bool {
	return len(e.SubErrors) == 0 && e.Message == ""
}

func NewErrorObj(message string, subErrors ...ErrorObj) ErrorObj {
	return ErrorObj{
		Message:   message,
		SubErrors: subErrors,
	}
}

func EmptyErrorObj() ErrorObj {
	return ErrorObj{
		Message:   "",
		SubErrors: []ErrorObj{},
	}
}
