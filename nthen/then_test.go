package nthen

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuture_CollectFor(t *testing.T){
	var ft = NewFuture()
	require.False(t, ft.IsResolved())
	require.NoError(t, ft.WithValue(1))

	var ft2 = NewFuture()
	require.False(t, ft2.IsResolved())
	require.NoError(t, ft2.WithError(errors.New("bad op")))

	var ft3 = CollectFor(ft, ft2)
	ft3.Wait()
	require.True(t, ft.IsResolved())
	require.Len(t, ft3.Value(), 2)
}

func TestFuture_WaitFor(t *testing.T){
	var ft = NewFuture()
	require.False(t, ft.IsResolved())
	require.NoError(t, ft.WithValue(1))

	var ft2 = NewFuture()
	require.False(t, ft2.IsResolved())
	require.NoError(t, ft2.WithError(errors.New("bad op")))

	var ft3 = WaitFor(ft, ft2)
	ft3.Wait()
	require.True(t, ft.IsResolved())
	require.Len(t, ft3.Value(), 1)
	require.Error(t, ft3.Err())
}

func TestFuture(t *testing.T){
	var ft = NewFuture()
	require.False(t, ft.IsResolved())

	require.NoError(t, ft.WithValue(1))
	require.True(t, ft.IsResolved())
	require.NoError(t, ft.Err())
	require.NotNil(t, ft.Value())
}

func TestFuture_Error(t *testing.T){
	var ft = NewFuture()
	require.False(t, ft.IsResolved())

	require.NoError(t, ft.WithError(errors.New("yo")))
	require.True(t, ft.IsResolved())
	require.Error(t, ft.Err())
	require.Nil(t, ft.Value())
}