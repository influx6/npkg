package nenv_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gokit/npkg/nenv"
)

func TestEnv(t *testing.T) {
	var envs, err = nenv.New("", func() ([]string, error) {
		return []string{
			"BADGER_MS=10ms",
			"BADGER_MSI=msa",
			"BADGER_IN_KEY=ricko",
			"BADGER_DAY_OFF=tuesday",
			"BAD_SUNDAY=rabbit",
			"rick=rabbit",
		}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, envs)

	require.False(t, envs.Has("sunday"))
	require.False(t, envs.Has("ms"))
	require.True(t, envs.Has("badger_msi"))
	require.True(t, envs.Has("BADGER_in_key"))
	require.True(t, envs.Has("BADGER_DAY_OFF"))
}

func TestEnvWithPrefix(t *testing.T) {
	var envs, err = nenv.New("badger", func() ([]string, error) {
		return []string{
			"BADGER_MS=10ms",
			"BADGER_MSI=msa",
			"BADGER_IN_KEY=ricko",
			"BADGER_DAY_OFF=tuesday",
			"BAD_SUNDAY=rabbit",
			"rick=rabbit",
		}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, envs)

	require.False(t, envs.Has("sunday"))
	require.True(t, envs.Has("ms"))
	require.True(t, envs.Has("msi"))
	require.True(t, envs.Has("in_key"))
	require.True(t, envs.Has("day_off"))
}
