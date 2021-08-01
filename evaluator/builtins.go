package evaluator

import (
	"strings"

	"github.com/skanehira/monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{
					Value: int64(len(arg.Value)),
				}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return newError("wrong number of arguments. got=%d, want=one or more", len(args))
			}

			var result []string

			for _, a := range args {
				result = append(result, a.Inspect())
			}

			return &object.String{
				Value: strings.Join(result, " "),
			}
		},
	},
}
