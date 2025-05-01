package evaluator

import "main/object"

type Environment struct {
	Store map[string]object.Object
	Outer *Environment
}

func NewEnvironment() Environment {
	return Environment{Store: make(map[string]object.Object), Outer: nil}
}

func NewEnclosedEnvironment(env Environment) Environment {
	return Environment{Store: make(map[string]object.Object), Outer: &env}
}

func (e *Environment) Create(name string, value object.Object) {
	if val := e.Get(name); val != nil {
		panic("variable already exists")
	}
	e.Store[name] = value
}

func (e *Environment) Set(name string, value object.Object) {
	if val := e.Get(name); val != nil {
		e.Store[name] = value
	} else if e.Outer != nil {
		e.Outer.Set(name, value)
	} else {
		panic("variable not found")
	}
}

func (e *Environment) Get(name string) object.Object {
	if value, ok := e.Store[name]; ok {
		return value
	} else if e.Outer != nil {
		return e.Outer.Get(name)
	}
	return nil
}
