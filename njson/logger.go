package njson

import (
	"sync"

	"github.com/influx6/npkg"
)

func jsonMaker() npkg.Encoder {
	return JSONB()
}

var (
	logPool = &sync.Pool{
		New: func() interface{} {
			return &LogStack{npkg.NewWriteStack(jsonMaker, nil)}
		},
	}
)

func Log(log Logger) *LogStack {
	var writer = &writeLogger{log}
	var newStack, isStack = logPool.Get().(*LogStack)
	if !isStack {
		newStack = &LogStack{npkg.NewWriteStack(jsonMaker, nil)}
	}
	newStack.SetWriter(writer)
	return newStack
}

func ReleaseLogStack(ll *LogStack) {
	logPool.Put(ll)
}

type LogStack struct {
	*npkg.WriteStack
}

type Logger interface {
	Log(*JSON)
}

type writeLogger struct {
	Logger
}

func (l *writeLogger) Write(v npkg.Encoded) {
	if vjson, ok := v.(*JSON); ok {
		l.Log(vjson)
	}
}
