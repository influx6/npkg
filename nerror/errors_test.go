package nerror_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gokit/npkg/nerror"
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
