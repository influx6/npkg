package nreflect_test

import (
	"testing"

	"github.com/gokit/npkg/nreflect"
	"github.com/stretchr/testify/require"
)

func TestDeepCopy(t *testing.T) {
	var specs = []struct {
		Value interface{}
		Fail  bool
	}{
		{
			Value: "one",
		},
		{
			Value: "1",
		},
		{
			Value: 1,
		},
		{
			Value: float64(1),
		},
		{
			Value: float32(1),
		},
		{
			Value: true,
		},
		{
			Value: false,
		},
		{
			Value: []string{"one", "two", "three"},
		},
		{
			Value: []int{1, 2, 3, 4, 5},
		},
	}

	for _, spec := range specs {
		var dest, err = nreflect.DeepCopy(spec.Value)
		require.NoError(t, err)
		require.NotNil(t, dest)
		require.Equal(t, dest, spec.Value)
	}
}
