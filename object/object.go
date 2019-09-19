package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ollybritton/monkey/ast"
)

// ObjectType is a representation of a type of object, such as BOOLEAN or INT.
type ObjectType string

// Definition of object types.
const (
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"

	NULL_OBJ    = "NULL"
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"
	ARRAY_OBJ   = "ARRAY"

	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
)

// Object is an interface that represents an object inside the program. The reason this is an interface and not a struct
// is because every value needs a different internal representation.
type Object interface {
	Type() ObjectType // Type is the ObjectType value of that object.
	Inspect() string  // Inspect gets a string representation of the object.
}

// BuiltinFunction is a function that is built-in to the interpreter, such as len()
type BuiltinFunction func(args ...Object) Object

// Integer represents an integer, such as "5" or "1232".
type Integer struct {
	Value int64
}

// Inspect gets the literal value of the integer, as a string.
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type gets the INTEGER_OBJ value.
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Boolean represents a bool, either "true" or "false".
type Boolean struct {
	Value bool
}

// Inspect gets the value of the boolean as a string.
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Type gets the BOOLEAN_OBJ value.
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Null represents null/nil, the lack of a value.
type Null struct{}

// Inspect gets the string "null".
func (n *Null) Inspect() string { return "null" }

// Type gets the NULL_OBJ type.
func (n *Null) Type() ObjectType { return NULL_OBJ }

// ReturnValue represents a value that is being returned.
type ReturnValue struct {
	Value Object
}

// Inspect gets the string of the return value.
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Type gets the RETURN_VALUE_OBJ type.
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Error represents an error that occurs.
type Error struct {
	Message string
}

// Inspect gets the error message.
func (e *Error) Inspect() string { return "error: " + e.Message }

// Type gets the ERROR_OBJ type.
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Function represents a function that is being evaluated.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type gets the FUNCTION_OBJ type.
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// Inspect gets the definition of the function as a string.
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// String represents a string in the program.
type String struct {
	Value string
}

// Type gets the STRING_OBJ type.
func (s *String) Type() ObjectType { return STRING_OBJ }

// Inspect gets the literal value of the string.
func (s *String) Inspect() string { return s.Value }

// Array represents the array data structure.
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Builtin wraps a built-in function so that it is usable inside the program.
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns the BUILTIN_OBJ type.
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Inspect tells us that this is a builtin function.
func (b *Builtin) Inspect() string { return "builtin function" }
