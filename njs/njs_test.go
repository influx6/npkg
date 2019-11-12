package njs_test

import (
	"strings"
	"testing"

	"github.com/influx6/npkg/njs/producers"

	"github.com/influx6/npkg/njs"
	"github.com/stretchr/testify/require"
)

func TestClassBlock(t *testing.T) {
	var funcDecl = njs.ClassBlock(njs.Class{
		Name:    "Mix",
		Extends: "",
		Body: njs.Join(
			njs.Function(njs.Func{
				Name:   "Mixer",
				Args:   "name",
				Body:   `console.log("Mixing: ", name)`,
				Async:  true,
				Static: false,
			}),
			njs.Function(njs.Func{
				Name:   "Mix",
				Args:   "name",
				Body:   `console.log("Mixing: ", name)`,
				Async:  false,
				Static: true,
			}),
		),
	})
	require.Equal(t, funcDecl, "\n\n\n\n\n\n\n\n\n\n\n\n\n\nclass Mix {\n    \n    \n  "+
		"  \n    async function Mixer(name) {\n        console.log(\"Mixing: \", name)\n    };\n "+
		"   \n    \n    \n    static Mix(name) {\n        console.log(\"Mixing: \", name)\n    };\n};")
}

func TestAsyncFunctionCall(t *testing.T) {
	var funcDecl = njs.FunctionCall(njs.Func{
		Name:   "Mix",
		Args:   "name",
		Body:   "",
		Async:  true,
		Static: false,
	})
	require.Equal(t, funcDecl, `await Mix(name);`)
}

func TestFunctionCall(t *testing.T) {
	var funcDecl = njs.FunctionCall(njs.Func{
		Name:   "Mix",
		Args:   "name",
		Body:   "",
		Async:  false,
		Static: false,
	})
	require.Equal(t, funcDecl, `Mix(name);`)
}

func TestStaticFunction(t *testing.T) {
	var funcDecl = njs.Function(njs.Func{
		Name:   "Mix",
		Args:   "name",
		Body:   `console.log("Mixing: ", name)`,
		Async:  false,
		Static: true,
	})
	require.Equal(t, funcDecl, "\n\n\nstatic Mix(name) {\n    console.log(\"Mixing: \", name)\n};")
}

func TestAsyncStaticFunction(t *testing.T) {
	var funcDecl = njs.Function(njs.Func{
		Name:   "Mix",
		Args:   "name",
		Body:   `console.log("Mixing: ", name)`,
		Async:  true,
		Static: true,
	})

	require.Equal(t, funcDecl, "\n\n\nstatic async Mix(name) {\n    console.log(\"Mixing: \", name)\n};")
}

func TestAsyncFunction(t *testing.T) {
	var funcDecl = njs.Function(njs.Func{
		Name:   "Mix",
		Args:   "name",
		Body:   `console.log("Mixing: ", name)`,
		Async:  true,
		Static: false,
	})

	require.Equal(t, funcDecl, "\n\n\nasync function Mix(name) {\n    console.log(\"Mixing: \", name)\n};")
}

func TestFunction(t *testing.T) {
	var funcDecl = njs.Function(njs.Func{
		Name:   "Mix",
		Args:   "name",
		Body:   `console.log("Mixing: ", name)`,
		Async:  false,
		Static: false,
	})

	require.Equal(t, funcDecl, "\n\n\nfunction Mix(name) {\n    console.log(\"Mixing: \", name)\n};")
}

func TestVariable(t *testing.T) {
	var variableDecl = njs.Variable(njs.Var{Name: "name", Value: `"Alex"`})
	require.Equal(t, variableDecl, `var name = "Alex";`)
}

func TestComposeProducers(t *testing.T) {
	var composed = njs.ComposeProducers(
		true,
		producers.ComposeVariable(njs.Var{Name: "name", Value: `"Alex"`}),
		producers.ComposeVariable(njs.Var{Name: "age", Value: `20`}),
		producers.ComposeVariable(njs.Var{Name: "range", Value: `1`}),
	)

	var builder strings.Builder
	require.NoError(t, composed(&builder))
	require.Equal(t, builder.String(), "var name = \"Alex\";\nvar age = 20;\nvar range = 1;")
}
