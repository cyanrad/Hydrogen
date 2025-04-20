package object

import "strconv"

type Object interface {
	Type() ObjectType
	Insepect() string
}

type IntegerObj struct {
	Value int64
}

func (i IntegerObj) Type() ObjectType { return INT_OBJ }
func (i IntegerObj) Insepect() string { return strconv.Itoa(int(i.Value)) }

type BooleanObj struct {
	Value bool
}

func (b BooleanObj) Type() ObjectType { return BOOLEAN_OBJ }
func (b BooleanObj) Insepect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type NullObj struct{}

func (n NullObj) Type() ObjectType { return NULL_OBJ }
func (n NullObj) Insepect() string { return "null" }
