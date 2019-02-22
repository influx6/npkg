package history

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// vars
var (
	stackSize = 1 << 6
	question  = "???"
)

//**************************************************************
// StackTrace
//**************************************************************

// Trace defines a structure which contains the stack, start and endtime
// on a given from a trace call to trace a given call with stack details
// and execution time.
type StackTrace struct {
	Stack      []byte    `json:"stack"`
	Package    string    `json:"Package"`
	File       string    `json:"file"`
	Function   string    `json:"function"`
	LineNumber int       `json:"line_number"`
	Time       time.Time `json:"end_time"`
}

// TraceAt returns a StackTrace object from the giving depth.
func TraceAt(depth int) StackTrace {
	trace := make([]byte, stackSize)
	trace = trace[:runtime.Stack(trace, false)]

	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = question
	}

	pkgFileBase := file

	var pkg, pkgFile, functionName string
	functionName, _, _ = GetMethod(3)

	if file != question {
		pkgPieces := strings.SplitAfter(pkgFileBase, "/src/")
		if len(pkgPieces) > 1 {
			pkgFileBase = pkgPieces[1]
		}

		pkg = filepath.Dir(pkgFileBase)
		pkgFile = filepath.Base(pkgFileBase)
	}

	return StackTrace{
		Package:    pkg,
		LineNumber: line,
		Stack:      trace,
		File:       pkgFile,
		Time:       time.Now(),
		Function:   functionName,
	}
}

// String returns the giving trace timestamp for the execution time.
func (t StackTrace) String() string {
	return fmt.Sprintf("[Package=%q, File=%q, Time=%+q]", t.Package, t.File, t.Time)
}

//**************************************************************
// GetMethod
//**************************************************************

// GetMethod returns the caller of the function that called it :)
func GetMethod(depth int) (string, string, int) {
	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(depth, fpcs)
	if n == 0 {
		return "Unknown()", "???", 0
	}

	funcPtr := fpcs[0]
	funcPtrArea := funcPtr - 1

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(funcPtrArea)
	if fun == nil {
		return "Unknown()", "???", 0
	}

	fileName, line := fun.FileLine(funcPtrArea)

	// return its name
	return fun.Name(), fileName, line
}

// GetMethodGraph returns the caller of the function that called it :)
func GetMethodGraph(depth int) CallGraph {
	var graph CallGraph
	graph.In.File = "???"
	graph.By.File = "???"
	graph.In.Function = "Unknown()"
	graph.By.Function = "Unknown()"

	// we get the callers as uintptrs - but we just need 1
	lower := make([]uintptr, 2)

	// skip 3 levels to get to the caller of whoever called Caller()
	if n := runtime.Callers(depth, lower); n == 0 {
		return graph
	}

	lowerPtr := lower[0] - 1
	higherPtr := lower[1] - 1

	// get the info of the actual function that's in the pointer
	if lowerFun := runtime.FuncForPC(lowerPtr); lowerFun != nil {
		graph.By.File, graph.By.Line = lowerFun.FileLine(lowerPtr)
		graph.By.Function = lowerFun.Name()

	}

	if higherFun := runtime.FuncForPC(higherPtr); higherFun != nil {
		graph.In.File, graph.In.Line = higherFun.FileLine(higherPtr)
		graph.In.Function = higherFun.Name()
	}

	return graph
}
