package njson

import (
	"github.com/influx6/npkg"
)

func jsonMaker() npkg.Encoder {
	return JSONB()
}

func Log(log Logger) *LogStack {
	var writer = &writeLogger{log}
	var newStack = &LogStack{npkg.NewWriteStack(jsonMaker, writer)}
	return newStack
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
