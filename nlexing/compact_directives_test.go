package nlexing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVariantDirectives(t *testing.T) {
	var specs = []struct {
		Directive  string
		Expected   []string
		ShouldFail bool
	}{
		{
			Directive: "xs:(hover:(bg-color-500!))recon",
			Expected: []string{
				"xs:hover:bg-color-500!",
				"xs:recon",
			},
		},
		{
			Directive: "xs:(hover:(bg-color-500!))",
			Expected: []string{
				"xs:hover:bg-color-500!",
			},
		},
		{
			Directive:  "xs:(text-color-400), fg-color-600)",
			ShouldFail: true,
			Expected:   []string{},
		},
		{
			Directive: "xs:(hover:(bg-color-500, flex:text-color-400), fg-color-600)",
			Expected: []string{
				"xs:hover:bg-color-500",
				"xs:hover:flex:text-color-400",
				"xs:fg-color-600",
			},
		},
		{
			Directive: "xs:(hover:(bg-color-500, text-color-400), fg-color-600)",
			Expected: []string{
				"xs:hover:bg-color-500",
				"xs:hover:text-color-400",
				"xs:fg-color-600",
			},
		},
		{
			Directive: "xs:(hover:(bg-color-500), fg-color-600)",
			Expected: []string{
				"xs:hover:bg-color-500",
				"xs:fg-color-600",
			},
		},
		{
			Directive: "xs:(hover:bg-color-500, fg-color-600)",
			Expected: []string{
				"xs:hover:bg-color-500",
				"xs:fg-color-600",
			},
		},
		{
			Directive: "xs:(bg-color-500!, fg-color-600)",
			Expected: []string{
				"xs:bg-color-500!",
				"xs:fg-color-600",
			},
		},
		{
			Directive: "xs:(bg-color-500, fg-color-600)",
			Expected: []string{
				"xs:bg-color-500",
				"xs:fg-color-600",
			},
		},
		{
			Directive: "xs:(bg-color-500)",
			Expected: []string{
				"xs:bg-color-500",
			},
		},
		{
			Directive: "xs:bg-color-500",
			Expected: []string{
				"xs:bg-color-500",
			},
		},
		{
			Directive: "bg-color-500",
			Expected: []string{
				"bg-color-500",
			},
		},
	}

	for index, spec := range specs {
		var parsedList, parseErr = ParseVariantDirectives(spec.Directive)
		if spec.ShouldFail {
			require.Error(t, parseErr, "Failed at index %d", index)
			continue
		}

		require.NoError(t, parseErr, "Failed at index %d", index)
		require.Equal(t, spec.Expected, parsedList, "Failed at index %d", index)
	}
}
