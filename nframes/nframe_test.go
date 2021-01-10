package nframes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getName() string {
	return GetCallerName()
}

func getNameWith2Skip() string {
	return GetCallerNameWith(2)
}

func TestGetCallerName(t *testing.T) {
	require.Equal(t, "github.com/influx6/npkg/nframes.getName", getName())
	require.NotEqual(t, "github.com/influx6/npkg/nframes.getNameWith2Skip", getNameWith2Skip())
	require.Equal(t, "github.com/influx6/npkg/nframes.TestGetCallerName", getNameWith2Skip())
	require.NotEqual(t, getName(), getNameWith2Skip())
}
