package nreflect_test

import (
	"testing"

	"github.com/gokit/npkg/nreflect"
	"github.com/stretchr/testify/require"
)

func TestDeepCopy(t *testing.T) {
	type moma struct {
		Number int
	}

	var specs = []struct {
		Value interface{}
		Fail  bool
	}{
		{
			Value: map[string]int{"one": 1, "two": 2, "three": 3},
		},
		{
			Value: []string{"one", "two", "three"},
		},
		{
			Value: [3]string{"one", "two", "three"},
		},
		{
			Value: []int{1, 2, 3, 4, 5},
		},
		{
			Value: moma{Number: 20},
		},
		{
			Value: &moma{Number: 20},
		},
		{
			Value: struct {
				Name string
			}{
				Name: "wonder",
			},
		},
		{
			Value: &(struct {
				Name string
			}{
				Name: "wonder",
			}),
		},
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
			Value: float64(1.6),
		},
		{
			Value: float32(11.1),
		},
		{
			Value: true,
		},
		{
			Value: false,
		},
	}

	for _, spec := range specs {
		var dest, err = nreflect.DeepCopy(spec.Value)
		require.NoError(t, err)
		require.NotNil(t, dest)
		require.Equal(t, spec.Value, dest)
	}
}
