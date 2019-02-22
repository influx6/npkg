package nerror_test

import (
	"testing"

	nerror "github.com/gokit/npkg/nerror"
	"github.com/stretchr/testify/assert"
)

func TestErrorCallGraph(t *testing.T) {
	newErr, ok := (nerror.New("failed connection: %s", "10.9.1.0")).(*nerror.PointingError)
	assert.True(t, ok)
	assert.Nil(t, newErr.Stack)
	assert.Equal(t, newErr.Message, "failed connection: 10.9.1.0")
	assert.Contains(t, newErr.Call.By.File, "errors_test.go")
}

func TestErrorWithStack(t *testing.T) {
	newErr, ok := (nerror.NewStack("failed connection: %s", "10.9.1.0")).(*nerror.PointingError)
	assert.True(t, ok)
	assert.NotNil(t, newErr.Stack)
	assert.Equal(t, newErr.Message, "failed connection: 10.9.1.0")
	assert.Contains(t, newErr.Call.By.File, "errors_test.go")
}
