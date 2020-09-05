package njobs

import (
	"fmt"
	"io"
	"os"
	"path"
	"reflect"

	"github.com/influx6/npkg/nerror"
)

type Job interface {
	Do(interface{}) (interface{}, error)
}

type JobFunction func(interface{}) (interface{}, error)

func (fn JobFunction) Do(data interface{}) (interface{}, error) {
	return fn(data)
}

// MoveLastForward moves the last result forward compared to the result
// from this doer.
func MoveLastForward(doer JobFunction) JobFunction {
	return func(lastResult interface{}) (interface{}, error) {
		var _, newErr = doer(lastResult)
		if newErr != nil {
			return nil, nerror.WrapOnly(newErr)
		}
		return lastResult, nil
	}
}

// ReadDir creates a new directory with giving mode.
func ReadDir(dir string) JobFunction {
	return func(rootDirData interface{}) (interface{}, error) {
		var rootDir, ok = rootDirData.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		var targetDir = path.Join(rootDir, dir)
		var openedDir, err = os.Open(targetDir)
		if err != nil {
			return nil, nerror.WrapOnly(err)
		}
		var dirListing, dirListingErr = openedDir.Readdir(-1)
		if dirListingErr != nil {
			return nil, nerror.WrapOnly(dirListingErr)
		}
		return dirListing, nil
	}
}

// ForEach expects a type of list of items and a function handle each items.
// If skipResult is false, then the returned values of the doer function is gathered
// into a new array and returned else, a nil array is always returned.
//
// Because the results are gathered into a new array when skipResult is false,
// this means if an error occurred midway, you will receive an array with partial
// results and an error.
func ForEach(doer JobFunction, skipResult bool) JobFunction {
	return func(targetList interface{}) (interface{}, error) {
		var refValue = reflect.ValueOf(targetList)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		var refKind = refValue.Kind()
		if refKind != reflect.Slice && refKind != reflect.Array {
			return nil, nerror.New("argument is neither a slice or array")
		}

		var results []interface{}

		if !skipResult {
			results = make([]interface{}, 0, refValue.Len())
		}

		for i := 0; i < refValue.Len(); i++ {
			var item = refValue.Index(i)
			var result, resultErr = doer(item.Interface())
			if resultErr != nil {
				return results, nerror.Wrap(resultErr, "Failed processing: %#v", item)
			}
			if skipResult {
				continue
			}
			results = append(results, result)
		}

		return results, nil
	}
}

// Mkdir creates a new directory with giving mode.
func Mkdir(dir string, mod os.FileMode) JobFunction {
	return func(rootDirData interface{}) (interface{}, error) {
		var rootDir, ok = rootDirData.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		var targetDir = path.Join(rootDir, dir)
		if err := os.MkdirAll(targetDir, mod); err != nil && err != os.ErrExist {
			return nil, nerror.WrapOnly(err)
		}
		return targetDir, nil
	}
}

// Println deletes incoming path string
func Println(format string, writer io.Writer) JobFunction {
	return Printf(format+"\n", writer)
}

// Printf deletes incoming path string
func Printf(format string, writer io.Writer) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		if _, err := fmt.Fprintf(writer, format, dir); err != nil {
			return dir, err
		}
		return dir, nil
	}
}

// DeletePath deletes incoming path string
func DeletePath() JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var rootDir, ok = dir.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		var err = os.RemoveAll(rootDir)
		if err != nil {
			return rootDir, nerror.WrapOnly(err)
		}
		return rootDir, nil
	}
}

// DeleteDir returns a new function to delete a dir provided.
func DeleteDir(targetFile string) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var err = os.RemoveAll(targetFile)
		if err != nil {
			return targetFile, nerror.WrapOnly(err)
		}
		return targetFile, nil
	}
}

// DeleteFile returns a new function to delete a file provided.
func DeleteFile(targetFile string) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var err = os.Remove(targetFile)
		if err != nil {
			return targetFile, nerror.WrapOnly(err)
		}
		return targetFile, nil
	}
}

// DeleteDirectoryFrom returns a new function to delete a file within directory passed to function.
func DeleteDirectoryFrom(name string) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var rootDir, ok = dir.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		var targetDir = path.Join(rootDir, name)
		var err = os.RemoveAll(targetDir)
		if err != nil {
			return targetDir, nerror.WrapOnly(err)
		}
		return targetDir, nil
	}
}

// DeleteFileFrom returns a new function to delete a file within directory passed to function.
func DeleteFileFrom(name string) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var rootDir, ok = dir.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		var targetFile = path.Join(rootDir, name)
		var err = os.Remove(targetFile)
		if err != nil {
			return targetFile, nerror.WrapOnly(err)
		}
		return targetFile, nil
	}
}

// JoinPath expects to receive a string which is path which it applies a the
// join string to.
func JoinPath(join string) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var rootDir, ok = dir.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		return path.Join(rootDir, join), nil
	}
}

// BackupPath expects to receive a string which is path which it applies a '..' to
// to backup to the root directory.
func BackupPath() JobFunction {
	return JoinPath("..")
}

// NewFile returns a new function to create a file within directory passed to function.
func NewFile(name string, mod os.FileMode, r io.Reader) JobFunction {
	return func(dir interface{}) (interface{}, error) {
		var rootDir, ok = dir.(string)
		if !ok {
			return nil, nerror.New("Expected rootDir path string as input")
		}
		var targetFile = path.Join(rootDir, name)
		var createdFile, err = os.OpenFile(targetFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mod)
		if err != nil {
			return nil, nerror.WrapOnly(err)
		}
		defer createdFile.Close()
		if _, err := io.Copy(createdFile, r); err != nil {
			return nil, nerror.WrapOnly(err)
		}
		return targetFile, nil
	}
}

// Jobs manages a series of sequential jobs to be executed
// one after another.
type Jobs struct {
	jobs []Job
}

func (j *Jobs) Add(jb Job) {
	j.jobs = append(j.jobs, jb)
}

func (j *Jobs) Do(data interface{}) (interface{}, error) {
	var err error
	var d = data
	for _, job := range j.jobs {
		d, err = job.Do(d)
		if err != nil {
			return d, nerror.WrapOnly(err)
		}
	}
	return d, nil
}
