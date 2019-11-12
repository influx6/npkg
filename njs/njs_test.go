package njs_test

import (
	"strings"
	"testing"

	"github.com/influx6/npkg/njs/producers"

	"github.com/influx6/npkg/njs"
	"github.com/stretchr/testify/require"
)

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
