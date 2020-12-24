package nexec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// nerror ...
var (
	ErrCommandFailed = errors.New("Command failed to execute succcesfully")
)

type Log interface {
}

// CommanderOption defines a function type that aguments a commander's field.
type CommanderOption func(*Commander)

// Command sets the command for the Commander.
func Command(format string, m ...interface{}) CommanderOption {
	return func(cm *Commander) {
		cm.Command = fmt.Sprintf(format, m...)
	}
}

// SubCommands sets the subcommands for the Commander exec call.
// If subcommands are set then the Binary, Flag and Command are ignored
// and the values of the subcommand is used.
func SubCommands(p ...string) CommanderOption {
	return func(cm *Commander) {
		cm.SubCommands = p
	}
}

// Dir sets the Directory for the Commander exec call.
func Dir(p string) CommanderOption {
	return func(cm *Commander) {
		cm.Dir = p
	}
}

// Binary sets the binary command for the Commander.
func Binary(bin string, flag string) CommanderOption {
	return func(cm *Commander) {
		cm.Binary = bin
		cm.Flag = flag
	}
}

// Timeout sets the commander to run in synchronouse mode.
func Timeout(d time.Duration) CommanderOption {
	return func(cm *Commander) {
		cm.Timeout = d
	}
}

// Sync sets the commander to run in synchronouse mode.
func Sync() CommanderOption {
	return SetAsync(false)
}

// Async sets the commander to run in asynchronouse mode.
func Async() CommanderOption {
	return SetAsync(true)
}

// SetAsync sets the command for the Commander.
func SetAsync(b bool) CommanderOption {
	return func(cm *Commander) {
		cm.Async = b
	}
}

// Input sets the input reader for the Commander.
func Input(in io.Reader) CommanderOption {
	return func(cm *Commander) {
		cm.In = in
	}
}

// Output sets the output writer for the Commander.
func Output(out io.Writer) CommanderOption {
	return func(cm *Commander) {
		cm.Out = out
	}
}

// Err sets the error writer for the Commander.
func Err(err io.Writer) CommanderOption {
	return func(cm *Commander) {
		cm.Err = err
	}
}

// Envs sets the map of environment for the Commander.
func Envs(envs map[string]string) CommanderOption {
	return func(cm *Commander) {
		cm.Envs = envs
	}
}

// Apply takes the giving series of CommandOption returning a function that always applies them to passed in commanders.
func Apply(ops ...CommanderOption) CommanderOption {
	return func(cm *Commander) {
		for _, op := range ops {
			op(cm)
		}
	}
}

// ApplyImmediate applies the options immediately to the Commander.
func ApplyImmediate(cm *Commander, ops ...CommanderOption) *Commander {
	for _, op := range ops {
		op(cm)
	}

	return cm
}

// Commander runs provided command within a /bin/sh -c "{COMMAND}", returning
// response associatedly. It also attaches if provided stdin, stdout and stderr readers/writers.
// Commander allows you to set the binary to use and flag, where each defaults to /bin/sh for binary
// and -c for flag respectively.
type Commander struct {
	Async       bool
	Command     string
	SubCommands []string
	Timeout     time.Duration
	Dir         string
	Binary      string
	Flag        string
	Envs        map[string]string
	In          io.Reader
	Out         io.Writer
	Err         io.Writer
}

// New returns a new Commander instance.
func New(ops ...CommanderOption) *Commander {
	cm := new(Commander)

	for _, op := range ops {
		op(cm)
	}

	return cm
}

// Exec executes giving command associated within the command with os/exec.
func (c *Commander) Exec(ctx context.Context) (int, error) {
	if c.Binary == "" {
		if runtime.GOOS != "windows" {
			c.Binary = "/bin/sh"
			if c.Flag == "" {
				c.Flag = "-c"
			}
		}

		if runtime.GOOS == "windows" {
			c.Binary = "cmd"
			if c.Flag == "" {
				c.Flag = "/C"
			}
		}
	}

	var cancel func()
	if c.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
	}

	if cancel != nil {
		defer cancel()
	}

	var execCommand []string

	switch {
	case c.Command == "" && len(c.SubCommands) != 0:
		execCommand = c.SubCommands
	case c.Command != "" && len(c.SubCommands) == 0 && c.Binary != "":
		execCommand = append(execCommand, c.Binary, c.Flag, c.Command)
	case c.Command != "" && len(c.SubCommands) != 0 && c.Binary != "":
		execCommand = append(append(execCommand, c.Binary, c.Flag, c.Command), c.SubCommands...)
	case c.Command != "" && len(c.SubCommands) == 0 && c.Binary == "":
		execCommand = append(execCommand, c.Command)
	case c.Command != "" && len(c.SubCommands) != 0:
		execCommand = append(append(execCommand, c.Command), c.SubCommands...)
	default:
		return -1, errors.New("commands with/without subcommands must be specified")
	}

	var errCopy bytes.Buffer
	var multiErr io.Writer

	if c.Err != nil {
		multiErr = io.MultiWriter(&errCopy, c.Err)
	} else {
		multiErr = &errCopy
	}

	cmder := exec.Command(execCommand[0], execCommand[1:]...)
	cmder.Dir = c.Dir
	cmder.Stderr = multiErr
	cmder.Stdin = c.In
	cmder.Stdout = c.Out
	cmder.Env = os.Environ()

	if c.Envs != nil {
		for name, val := range c.Envs {
			cmder.Env = append(cmder.Env, fmt.Sprintf("%s=%s", name, val))
		}
	}

	if !c.Async {
		err := cmder.Run()
		return getExitStatus(err), err
	}

	if err := cmder.Start(); err != nil {
		return getExitStatus(err), err
	}

	go func() {
		<-ctx.Done()
		if cmder.Process == nil {
			return
		}

		cmder.Process.Kill()
	}()

	if err := cmder.Wait(); err != nil {
		return getExitStatus(err), err
	}

	if cmder.ProcessState == nil {
		return 0, nil
	}

	if !cmder.ProcessState.Success() {
		return -1, ErrCommandFailed
	}

	return 0, nil
}

type exitStatus interface {
	ExitStatus() int
}

func getExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if e, ok := err.(exitStatus); ok {
		return e.ExitStatus()
	}
	if e, ok := err.(*exec.ExitError); ok {
		if ex, ok := e.Sys().(exitStatus); ok {
			return ex.ExitStatus()
		}
	}
	return 1
}
