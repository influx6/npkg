package njobs_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/influx6/npkg/njobs"
)

func TestJobFunction_OS_Jobs(t *testing.T) {
	var jobs njobs.Jobs
	jobs.Add(njobs.Mkdir("./tmp", 0777))
	jobs.Add(njobs.NewFile("sample.txt", 0777, bytes.NewReader([]byte("yo"))))
	jobs.Add(njobs.JobFunction(func(d interface{}) (interface{}, error) {
		require.Equal(t, "tmp/sample.txt", d)
		return d, nil
	}))
	jobs.Add(njobs.DeleteDir("./tmp"))

	var result, err = jobs.Do(".")
	require.NoError(t, err)
	require.Equal(t, "./tmp", result)

}

func TestJobFunction_Do(t *testing.T) {
	var jobs njobs.Jobs

	var doOne = njobs.JobFunction(func(d interface{}) (interface{}, error) {
		require.Equal(t, 1, d)
		return 2, nil
	})

	var doTwo = njobs.JobFunction(func(d interface{}) (interface{}, error) {
		require.Equal(t, 2, d)
		return 3, nil
	})

	jobs.Add(doOne)
	jobs.Add(doTwo)

	var result, err = jobs.Do(1)
	require.NoError(t, err)
	require.Equal(t, 3, result)
}

func TestJobFunction_Do_Failed(t *testing.T) {
	var jobs njobs.Jobs

	var doOne = njobs.JobFunction(func(d interface{}) (interface{}, error) {
		require.Equal(t, 1, d)
		return 2, errors.New("bad")
	})

	var doTwo = njobs.JobFunction(func(d interface{}) (interface{}, error) {
		require.Equal(t, 2, d)
		return 3, nil
	})

	jobs.Add(doOne)
	jobs.Add(doTwo)

	var result, err = jobs.Do(1)
	require.Error(t, err)
	require.Equal(t, 2, result)
}
