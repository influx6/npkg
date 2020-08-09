package nerror_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/influx6/npkg/nerror"
)

func TestErrorCallGraph(t *testing.T) {
	newErr, ok := (nerror.New("failed connection: %s", "10.9.1.0")).(*nerror.PointingError)
	assert.True(t, ok)
	assert.NotNil(t, newErr.Frames)
	assert.Equal(t, newErr.Message, "failed connection: 10.9.1.0")
}

func TestErrorWithStack(t *testing.T) {
	newErr, ok := (nerror.NewStack("failed connection: %s", "10.9.1.0")).(*nerror.PointingError)
	assert.True(t, ok)
	assert.NotNil(t, newErr.Frames)
	assert.Equal(t, newErr.Message, "failed connection: 10.9.1.0")
}

func TestErrorWithWrapOnly(t *testing.T) {
	newErr, ok := (nerror.WrapOnly(doBad())).(*nerror.PointingError)
	assert.True(t, ok)
	assert.NotNil(t, newErr.Frames)
	fmt.Printf("Stacks: \n%s\n", newErr.Error())
}

func doBad() error {
	return nerror.New("Very bad error")
}
