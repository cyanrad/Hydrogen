package evaluator

import "main/object"

type Environment struct {
	Store map[string]object.Object
}

func NewEnvironment() Environment {
	return Environment{Store: make(map[string]object.Object)}
}

func (e *Environment) CreateVariable(name string, value object.Object) {
	if _, ok := e.Store[name]; ok {
		panic("variable already exists")
	}
	e.Store[name] = value
}

func (e *Environment) SetVariable(name string, value object.Object) {
	if _, ok := e.Store[name]; ok {
		e.Store[name] = value
	} else {
		panic("variable not found")
	}
}

func (e *Environment) GetVariable(name string) object.Object {
	if value, ok := e.Store[name]; ok {
		return value
	}
	return nil
}
