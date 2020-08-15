package nbytes_test

import (
	"testing"

	"github.com/influx6/npkg/nbytes"
	"github.com/stretchr/testify/require"
)

var (
	sentences = []string{
		"I went into park stream all alone before the isle lands.",
		"Isle lands of YOR, before the dream verse began we found the diskin.",
		"Break fast in bed, love and eternality for ever",
		"Awaiting the ending seen of waiting for you?",
		"Done be such a waste!",
		"{\"log\":\"token\", \"centry\":\"20\"}",
	}
)

func TestBuilder(t *testing.T) {
	var bb = nbytes.NewBuildReader()
	bb.WriteString("be nice")

	require.Equal(t, "be nice", bb.String())
}

func TestBuilder_Reset(t *testing.T) {
	var bb = nbytes.NewBuildReader()
	bb.WriteString("be nice")
	bb.Reset(false)
	bb.WriteString("be nice")

	require.Equal(t, "be nice", bb.String())
}
