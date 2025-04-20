package object

import "strconv"

type Object interface {
	Type() ObjectType
	Inspect() string
}

type IntegerObj struct {
	Value int64
}

func (i IntegerObj) Type() ObjectType { return INT_OBJ }
func (i IntegerObj) Inspect() string  { return strconv.Itoa(int(i.Value)) }

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

type NullObj struct{}

func (n NullObj) Type() ObjectType { return NULL_OBJ }
func (n NullObj) Inspect() string  { return "null" }
