package evaluator

import (
	"fmt"

	"github.com/ollybritton/monkey/ast"
	"github.com/ollybritton/monkey/object"
)

// Reused objects with a fixed number of values.
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Eval evaluates an AST node and returns an object.Object representation of the result.
func Eval(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.BlockStatement:
		return evalBlockStatement(node, environment)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, environment)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, environment)
		if isError(val) {
			return val
		}

		environment.Set(node.Name.Value, val)

	// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, environment)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Env:        environment,
			Body:       node.Body,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, environment)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	}

	return nil
}

func evalProgram(stmts []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, environment)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, environment)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, environment *object.Environment) object.Object {
	condition := Eval(ie.Condition, environment)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, environment)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, environment)
	} else {
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendedFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnVal(evaluated)
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewExtendedEnvironment(fn.Env)

	for paramID, param := range fn.Parameters {
		env.Set(param.Value, args[paramID])
	}

	return env

}

func unwrapReturnVal(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool { return obj.Type() == object.ERROR_OBJ }
