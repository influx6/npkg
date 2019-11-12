// Package njs exists to provide a simple javascript text based code generation, exposing
// composable functions that can be used to create generated javascript code.
//
// Code is taken from https://github.com/mistermoe/js-code-generator.
//
package njs

import (
	"regexp"
	"strings"

	"github.com/influx6/npkg/nerror"
)

const (
	numTabs         = 1
	numSpacesPerTab = 4
)

var alphabets = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o",
	"p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

// Gen represents a generated code with attached input data.
type Gen struct {
	Data interface{}
	Code string
}

// WrapWithData will wrap provided code and attached
// giving data object into it returning a Gen object.
func WrapWithData(code string, data interface{}) Gen {
	return Gen{Data: data, Code: code}
}

// Args returns a joined string containing all values as comma seperated
// values within returned string.
func Args(args ...string) string {
	return strings.Join(args, ", ")
}

// Join joins collection of string with newlines.
func Join(cols ...string) string {
	return JoinWith("\n", cols...)
}

// JoinWith joins giving collection of strings with combiner.
func JoinWith(combiner string, cols ...string) string {
	return CleanCode(strings.Join(cols, combiner))
}

type Producer func() string
type Action func(builder *strings.Builder) error
type ActionProducer func(builder *strings.Builder, producers ...Producer)

// ComposeProducer composes series of string producers into a single Action
// to be applied to a giving string.Builder instance.
func ComposeProducers(indent bool, producers ...Producer) Action {
	return func(builder *strings.Builder) error {
		var total = len(producers)
		for index, producer := range producers {
			if _, err := builder.WriteString(producer()); err != nil {
				return nerror.WrapOnly(err)
			}
			if indent && index < total-1 {
				builder.WriteString("\n")
			}
		}
		return nil
	}
}

// ComposeActions composes series of Actions into a single action
// to be applied to a giving string.Builder instance.
func ComposeActions(actions ...Action) Action {
	return func(builder *strings.Builder) error {
		for _, action := range actions {
			if err := action(builder); err != nil {
				return nerror.WrapOnly(err)
			}
		}
		return nil
	}
}

type Var struct {
	Name  string
	Value string
}

/*
  @param Object data {
    {string} name,
    {string} value (optional)
  }
*/
func Variable(data Var) string {
	var code = `var ` + data.Name + ``
	if data.Value != "" {
		return CleanCode(code + ` = ` + data.Value + `;`)
	}
	return CleanCode(code + `;`)
}

/*
  @param Object data {
    {string} name,
    {string} value
  }
*/
func ReAssignVariable(data Var) string {
	return CleanCode(`` + data.Name + ` = ` + data.Value + `;`)
}

/*
  @param Object data {
    {string} name,
    {string} value (Optional)
  }
*/
func IncrementVariable(data Var) string {
	var code string
	if data.Value != "" {
		code = `` + data.Name + ` += ` + data.Value + `;`
	} else {
		code = `` + data.Name + `++`
	}
	return CleanCode(code)
}

/*
  @param Object data {
    {string} name,
    {string} value (Optional)
  }
*/
func DecrementVariable(data Var) string {
	var code string
	if data.Value != "" {
		code = `` + data.Name + ` -= ` + data.Value + `;`
	} else {
		code = `` + data.Name + `--`
	}
	return CleanCode(code)
}

type Func struct {
	Name  string
	Args  string
	Body  string
	Async bool
}

/*
  @param Object data {
    {string} name,
    {array} args (Optional),
	{bool} async
    {func} body
  }
  @returns {
    data: data,
    code: code
  }
*/
func Function(data Func) string {
	var code = `var ` + data.Name + ` = `
	if data.Async {
		code += `async `
	}
	code += `func(` + data.Args + `) {` + "\n" +
		`` + data.Body + "\n" + `};`
	return Indent(code)
}

func FunctionCall(data Func) string {
	var code string
	if data.Async {
		code += `await `
	}
	code += data.Name + `(` + data.Args + `);`
	return CleanCode(code)
}

/*
  @param Object data {
    {string} type
    {array} args
  }

  @returns Object {
    {string} code
  }
*/
func NewInstance(data struct {
	Type string
	Args string
}) string {
	var code = `new ` + data.Type
	code += `(` + data.Args + `);`
	return CleanCode(code)
}

func JSONParse(args string) string {
	var code = `JSON.parse(` + args + `);`
	return Indent(code)
}

func JSONStringify(args string) string {
	var code = `JSON.stringify(` + args + `);`
	return Indent(code)
}

func ChainObject(objName string) string {
	return objName + `.`
}

func Self() string {
	return ChainObject(`self`)
}

func Window() string {
	return ChainObject(`window`)
}

func ParseInt(args string) string {
	return FunctionCall(Func{
		Name:  "parseInt",
		Args:  args,
		Body:  "",
		Async: false,
	})
}

func ParseFloat(args string) string {
	return FunctionCall(Func{
		Name:  "parseFloat",
		Args:  args,
		Body:  "",
		Async: false,
	})
}

func Math() string {
	return ChainObject(`Math`)
}

func This() string {
	return ChainObject(`this`)
}

func Object(body string) string {
	var code = `{` + "\n" + body + "\n" + `};`
	return Indent(code)
}

type Obj struct {
	DotNotation bool
	PropName    string
	ObjName     string
}

/*
  @param Object data {
    {string} objName (Optional. Defaults to "this")
    {string} propName,
    {string} value,
    {boolean} dotNotation
  }
*/
func ObjectPropertyName(data Obj) string {
	var code string
	if data.DotNotation == false {
		if data.ObjName == "" {
			code = `this["`
		} else {
			code = data.ObjName + `["`
		}

		code += data.PropName + `"] `
		return CleanCode(code)
	}

	if data.ObjName == "" {
		code = `this.` + data.PropName
	} else {
		code = data.ObjName + `.` + data.PropName
	}

	return CleanCode(code)
}

type ObjWithValue struct {
	Obj
	Value string
}

/*
  @param Object data {
    {string} objName (Optional. Defaults to "this")
    {string} propName,
    {string} value,
    {boolean} dotNotation
  }
*/
func ObjectPropertyAssignment(data ObjWithValue) string {
	var code string
	if data.DotNotation == false {
		if data.ObjName == "" {
			code = `this["`
		} else {
			code = data.ObjName + `["`
		}

		code += data.PropName + `"] = ` + data.Value + `;`
		//name = ``+ data.ObjName +`["`+ data.PropName +`"]`;
		return CleanCode(code)
	}

	if data.ObjName == "" {
		code = `this.` + data.PropName + ` = ` + data.Value + `;`
	} else {
		code = data.ObjName + `.` + data.PropName + ` = ` + data.Value + `;`
	}

	return CleanCode(code)
}

type ObjFuncCall struct {
	ObjectName string
	FuncName   string
	Args       string
	Async      bool
}

type ObjFunc struct {
	ObjFuncCall
	Body string
}

/*
  @param Object data {
    {string} objName (Optional. Defaults to "this")
    {string} funcName,
    {array} args (Optional),
	{bool} async
    {func} body
  }
*/
func ObjectFunction(data ObjFunc) string {
	var code string
	if data.ObjectName == "" {
		code = `this.` + data.FuncName
	} else {
		code = data.ObjectName + `.` + data.FuncName
	}

	code += data.FuncName + ` = `
	if data.Async {
		code += `async `
	}

	code += `func(` + data.Args + `) {` + "\n" +
		`` + data.Body + `` + "\n" +
		`};`

	return Indent(code)
}

/*
  @param Object data {
    {string} objName
    {string} funcName,
	{bool} async
    {array} args
  }
*/
func ObjectFunctionCall(data ObjFunc) string {
	var code string
	if data.Async {
		code += `await `
	}

	if data.ObjectName == "" {
		code = `this.` + data.FuncName
	} else {
		code = data.ObjectName + `.` + data.FuncName
	}

	code += `(` + data.Args + `);`
	return CleanCode(code)
}

/*
  @param Object data {
    {string} condition,
    {func} body
  }
*/
func IfStatement(
	condition string,
	body string,
) string {
	var code = `if (` + condition + `) {` + "\n" +
		`` + body + `` + "\n" +
		`}`

	return Indent(code)
}

/*
  @param Object data {
    {string} condition,
    {func} body
  }
*/
func ElseIfStatement(
	condition string,
	body string,
) string {
	var code = `else if (` + condition + `) {` + "\n" +
		`` + body + `` + "\n" +
		`}`

	return Indent(code)
}

/*
  @param Object data {
    {func} body
  }
*/
func ElseStatement(body string) string {
	var code = `else {` + "\n" +
		`` + body + `` + "\n" +
		`}`

	return Indent(code)
}

type Loop struct {
	StartCondition  string
	StopCondition   string
	IncrementAction string
	Body            string
}

/*
  @param Object data {
    {string} startCondition,
    {string} stopCondition,
    {string} incrementAction,
    {func} body
  }
*/
func ForLoop(data Loop) string {
	var code = `for (` + data.StartCondition + `; ` + data.StopCondition + `; ` + data.IncrementAction + `) {` + "\n" +
		`` + data.Body + `` + "\n" +
		`}`

	return Indent(code)
}

/*
  @param Object data {
    {func} body,
  }
*/
func TryBlock(body string) string {
	var code = `try {` + "\n" + body + "\n" + `}`
	return Indent(code)
}

/*
  @param Object data {
    {string} arg
    {func} body,
  }
*/
func CatchBlock(
	exception string,
	body string,
) string {
	var code = `catch(`
	if exception != "" {
		code += exception
	}

	code += `) {` + "\n"
	code += `` + body + `` + "\n" + `}`
	return Indent(code)
}

// Promise returns new Promise with resolve and reject function parameter.
func Promise(args string, body string) string {
	var code = `new Promise((resolve, reject) => {` + "\n" + body + "\n" + `})`
	return Indent(code)
}

func ManyPromise(promiseList string) string {
	return Indent(`Promise.all([` + "\n" + promiseList + "\n" + `]);`)
}

func PromiseThen(args string, body string) string {
	var code = `.then((` + args + `) => {` + "\n" + body + "\n" + `})`
	return Indent(code)
}

func PromiseCatch(err string, body string) string {
	var code = `.catch((` + err + `) => {` + "\n" + body + "\n" + `})`
	return Indent(code)
}

/*
  @param Object data {
    {string} name
    {array} args (optional)
  }
*/
func ChainFunction(name string, args string) string {
	return `.` + name + `(` + args + `)`
}

/*
  @param Object data {
    {string} value
  }
*/
func ReturnStatement(value string) string {
	return `return ` + value + `;`
}

/*
  @param { array<string> } args
*/
func ConsoleLog(args ...string) string {
	var code = `console.log(` + Args(args...) + `);`
	return code + `);`
}

type TextBlock struct {
	Text  string
	Block bool
}

/*
  @param Object data {
    {string} text,
    {boolean} block (Optional. Defaults to false.)
  }
*/
func Comment(data TextBlock) string {
	if data.Block {
		var code = `/*` + "\n" +
			`` + data.Text + `` + "\n" +
			`*/`
		return Indent(code)
	}

	return `// ` + data.Text + ``
}

/*
  @param {string} code

  @description:
    Takes a rendered code and indents lines 2 through (n - 2),
    where n is the number of lines

  @returns: {string}
*/
func Indent(code string) string {
	var items = strings.Split(code, "\n")
	var mapped = make([]string, len(items))
	for i := 0; i < len(items); i++ {
		var line = items[i]
		if i == 0 || i == len(items)-1 {
			mapped = append(mapped, GenerateTabs(0)+CleanCode(line))
		}
		mapped = append(mapped, GenerateTabs(numTabs)+CleanCode(line))
	}
	return strings.Join(mapped, "\n")
}

// AddIndent returns a new line feed character.
func AddIndent() string {
	return "\n"
}

/*
  @param: {integer} numTabs
  @description: generates a string containing the number of tabs requested
  @returns: {string}
*/
func GenerateTabs(numTabs int) string {
	var tabs = ""
	var tabCount = numSpacesPerTab * numTabs
	for i := 0; i < tabCount; i += 1 {
		tabs += " "
	}
	return tabs
}

var cleanRegExp = regexp.MustCompile(";{2,}")

/*
  @param: {string} line
  @description: Cleans a line of code.
  @returns: {string}
*/
func CleanCode(line string) string {
	return cleanRegExp.ReplaceAllString(line, ";")
}

var iteratorPos = 8

/*
  UniqueIteratorName
  @description:
    returns a name for an iterator variable that hasn't been used.
    An example of an iterator variable would be `var i` in a for loop.
    This func Is used to prevent undesired behavior in a nested
    loop. Starts at 'i'

  @returns: {string} iterator
*/
func UniqueIteratorName() string {
	var numChars = 1
	var idx = iteratorPos

	var iterator = ""

	if iteratorPos >= 26 {
		if iteratorPos%26 == 0 {
			numChars = (iteratorPos / 26) + 1
		} else {
			numChars = iteratorPos / 26
		}
		idx = iteratorPos % 26
	}

	for i := 0; i < numChars; i++ {
		iterator += alphabets[idx]
	}

	iteratorPos += 1
	return iterator
}

// ResetIteratorPos resets the position of the current iterator to return back to 'i'
func ResetIteratorPos() {
	iteratorPos = 8
}
