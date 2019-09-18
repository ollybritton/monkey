package object

// Environment is a collection of objects associated with identifiers.
// It holds variables.
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new environment.
func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

// Get gets the object associated with an identifer.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	// If we have an outer environment, use that.
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

// Set sets a value inside the environment.
func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

// NewExtendedEnvironment creates a new extended environment from an exisitng one. This is used for functions.
func NewExtendedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}
