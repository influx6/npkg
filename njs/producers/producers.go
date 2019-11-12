package producers

import (
	. "github.com/influx6/npkg/njs"
)

func ComposeComment(data TextBlock) Producer {
	return func() string {
		return Comment(data)
	}
}

func ComposeConsoleLog(args ...string) Producer {
	return func() string {
		return ConsoleLog(args...)
	}
}

func ComposeReturnStatement(value string) Producer {
	return func() string {
		return ReturnStatement(value)
	}
}

func ComposeObject(body string) Producer {
	return func() string {
		return Object(body)
	}
}

func ComposeObjectPropertyAssignment(data ObjWithValue) Producer {
	return func() string {
		return ObjectPropertyAssignment(data)
	}
}

func ComposeVariable(data Var) Producer {
	return func() string {
		return Variable(data)
	}
}

func ComposeReAssignVariable(data Var) Producer {
	return func() string {
		return ReAssignVariable(data)
	}
}

func ComposeIncrementVariable(data Var) Producer {
	return func() string {
		return IncrementVariable(data)
	}
}

func ComposeDecrementVariable(data Var) Producer {
	return func() string {
		return DecrementVariable(data)
	}
}

func ComposeFunction(data Func) Producer {
	return func() string {
		return Function(data)
	}
}

func ComposeFunctionCall(data Func) Producer {
	return func() string {
		return FunctionCall(data)
	}
}

func ComposeObjectPropertyName(data Obj) Producer {
	return func() string {
		return ObjectPropertyName(data)
	}
}

func ComposeChainFunction(name, args string) Producer {
	return func() string {
		return ChainFunction(name, args)
	}
}

func ComposeNewInstance(data struct {
	Type string
	Args string
}) Producer {
	return func() string {
		return NewInstance(data)
	}
}

func ComposeObjectFunction(data ObjFunc) Producer {
	return func() string {
		return ObjectFunction(data)
	}
}

func ComposeObjectFunctionCall(data ObjFunc) Producer {
	return func() string {
		return ObjectFunctionCall(data)
	}
}

func ComposeIfStatement(
	condition string,
	body string,
) Producer {
	return func() string {
		return IfStatement(condition, body)
	}
}

func ComposeElseIfStatement(
	condition string,
	body string,
) Producer {
	return func() string {
		return ElseIfStatement(condition, body)
	}
}

func ComposeElseStatement(body string) Producer {
	return func() string {
		return ElseStatement(body)
	}
}

func ComposeForLoop(data Loop) Producer {
	return func() string {
		return ForLoop(data)
	}
}

func ComposeTryBlock(
	body string,
) Producer {
	return func() string {
		return TryBlock(body)
	}
}

func ComposeCatchBlock(
	exception string,
	body string,
) Producer {
	return func() string {
		return CatchBlock(exception, body)
	}
}

func ComposeJSONStringify(
	args string,
) Producer {
	return func() string {
		return JSONStringify(args)
	}
}

func ComposeJSONParse(
	args string,
) Producer {
	return func() string {
		return JSONParse(args)
	}
}

func ComposeChainObject(
	objName string,
) Producer {
	return func() string {
		return ChainObject(objName)
	}
}

func ComposeWindow() Producer {
	return func() string {
		return Window()
	}
}

func ComposeMath() Producer {
	return func() string {
		return Math()
	}
}

func ComposeSelf() Producer {
	return func() string {
		return Self()
	}
}

func ComposeParseFloat(args string) Producer {
	return func() string {
		return ParseFloat(args)
	}
}

func ComposeParseInt(args string) Producer {
	return func() string {
		return ParseInt(args)
	}
}

func ComposeThis() Producer {
	return func() string {
		return This()
	}
}

func ComposePromise(
	args string,
	body string,
) Producer {
	return func() string {
		return Promise(args, body)
	}
}

func ComposeManyPromise(
	args string,
) Producer {
	return func() string {
		return ManyPromise(args)
	}
}

func ComposePromiseThen(
	args string,
	body string,
) Producer {
	return func() string {
		return PromiseThen(args, body)
	}
}

func ComposePromiseCatch(
	err string,
	body string,
) Producer {
	return func() string {
		return PromiseCatch(err, body)
	}
}

// ComposeIndent returns function that returns
// a new line feed character.
func ComposeIndent() Producer {
	return func() string {
		return AddIndent()
	}
}
