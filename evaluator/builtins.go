package evaluator

import (
	"fmt"

	"github.com/ollybritton/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.String{Value: string(arg.Value[0])}
			case *object.Array:
				return arg.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) == 0 {
					return NULL
				}

				return &object.String{Value: string(arg.Value[len(arg.Value)-1])}
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				if len(arg.Value) == 0 {
					return NULL
				}

				return &object.String{Value: string(arg.Value[1:len(arg.Value)])}
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}

				return &object.Array{
					Elements: arg.Elements[1:len(arg.Elements)],
				}
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("wrong number of arguments in call to push: need push(array, elements...). got=%d, want=>1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Array{
					Elements: append(arg.Elements, args[1:]...),
				}
			}

			return NULL
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
