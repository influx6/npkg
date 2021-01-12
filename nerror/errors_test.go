package nerror_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/influx6/npkg/nerror"
)

func TestErrorCallGraph(t *testing.T) {
	newErr := nerror.New("failed connection: %s", "10.9.1.0")
	assert.NotNil(t, newErr.Frames)
	assert.Equal(t, newErr.Message, "failed connection: 10.9.1.0")
}

func TestErrorWithStack(t *testing.T) {
	newErr := nerror.NewStack("failed connection: %s", "10.9.1.0")
	assert.NotNil(t, newErr.Frames)
	assert.Equal(t, newErr.Message, "failed connection: 10.9.1.0")
}

func TestErrorWithWrapOnly(t *testing.T) {
	newErr := nerror.WrapOnly(doBad())
	assert.NotNil(t, newErr.Frames)
	assert.Equal(t, "Very bad error", nerror.ErrorMessage(newErr).GetMessage())
	fmt.Println(newErr.String())
}

func doBad() error {
	return nerror.WrapOnly(doBad2())
}

func doBad2() error {
	return nerror.WrapOnly(doBad3())
}

func doBad3() error {
	return nerror.New("Very bad error")
}
