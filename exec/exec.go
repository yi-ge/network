package exec

import (
	"io"
	"os"
	osexec "os/exec"
	"syscall"
)

// ErrExecutableNotFound is returned if the executable is not found.
var ErrExecutableNotFound = osexec.ErrNotFound

// Interface is an interface that presents a subset of the os/exec API.  Use this
// when you want to inject fakeable/mockable exec behavior.
type Interface interface {
	// Command returns a Cmd instance which can be used to run a single command.
	// This follows the pattern of package os/exec.
	Command(cmd string, args ...string) Cmd

	// LookPath wraps os/exec.LookPath
	LookPath(file string) (string, error)
}

// Cmd is an interface that presents an API that is very similar to Cmd from os/exec.
// As more functionality is needed, this can grow.  Since Cmd is a struct, we will have
// to replace fields with get/set method pairs.
type Cmd interface {
	// CombinedOutput runs the command and returns its combined standard output
	// and standard error.  This follows the pattern of package os/exec.
	CombinedOutput() ([]byte, error)
	// Output runs the command and returns standard output, but not standard err
	Output() ([]byte, error)
	SetDir(dir string)
	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(out io.Writer)
	Start() error
	Wait() error
	GetProcess() *os.Process
}

// ExitError is an interface that presents an API similar to os.ProcessState, which is
// what ExitError from os/exec is.  This is designed to make testing a bit easier and
// probably loses some of the cross-platform properties of the underlying library.
type ExitError interface {
	String() string
	Error() string
	Exited() bool
	ExitStatus() int
}

// Implements Interface in terms of really exec()ing.
type executor struct{}

// New returns a new Interface which will os/exec to run commands.
func New() Interface {
	return &executor{}
}

// LookPath is part of the Interface interface
func (executor *executor) LookPath(file string) (string, error) {
	return osexec.LookPath(file)
}

// Wraps exec.Cmd so we can capture errors.
type cmdWrapper osexec.Cmd

func (cmd *cmdWrapper) SetDir(dir string) {
	cmd.Dir = dir
}

func (cmd *cmdWrapper) SetStdin(in io.Reader) {
	cmd.Stdin = in
}

func (cmd *cmdWrapper) SetStdout(out io.Writer) {
	cmd.Stdout = out
}

func (cmd *cmdWrapper) SetStderr(out io.Writer) {
	cmd.Stderr = out
}

func (cmd *cmdWrapper) SetSystemOptions() {
	SetSystemOptions(cmd)
}

// CombinedOutput is part of the Cmd interface.
func (cmd *cmdWrapper) CombinedOutput() ([]byte, error) {
	cmd.Env = os.Environ()
	cmd.SetSystemOptions()
	out, err := (*osexec.Cmd)(cmd).CombinedOutput()
	if err != nil {
		return out, handleError(err)
	}
	return out, nil
}

// Command is part of the Interface interface.
func (executor *executor) Command(cmd string, args ...string) Cmd {
	c, a := parseCommand(cmd, args...)
	return (*cmdWrapper)(osexec.Command(c, a...))
}

func (cmd *cmdWrapper) Output() ([]byte, error) {
	cmd.Env = os.Environ()
	out, err := (*osexec.Cmd)(cmd).Output()
	if err != nil {
		return out, handleError(err)
	}
	return out, nil
}

func (cmd *cmdWrapper) Start() error {
	cmd.SetSystemOptions()
	err := (*osexec.Cmd)(cmd).Start()
	return err
}

func (cmd *cmdWrapper) Wait() error {
	err := (*osexec.Cmd)(cmd).Wait()
	return err
}

func (cmd *cmdWrapper) GetProcess() *os.Process {
	return (*osexec.Cmd)(cmd).Process
}

func handleError(err error) error {
	if ee, ok := err.(*osexec.ExitError); ok {
		// Force a compile fail if exitErrorWrapper can't convert to ExitError.
		var x ExitError = &ExitErrorWrapper{ee}
		return x
	}
	if ee, ok := err.(*osexec.Error); ok {
		if ee.Err == osexec.ErrNotFound {
			return ErrExecutableNotFound
		}
	}
	return err
}

// ExitErrorWrapper is an implementation of ExitError in terms of os/exec ExitError.
// Note: standard exec.ExitError is type *os.ProcessState, which already implements Exited().
type ExitErrorWrapper struct {
	*osexec.ExitError
}

var _ ExitError = ExitErrorWrapper{}

// ExitStatus is part of the ExitError interface.
func (eew ExitErrorWrapper) ExitStatus() int {
	ws, ok := eew.Sys().(syscall.WaitStatus)
	if !ok {
		panic("can't call ExitStatus() on a non-WaitStatus exitErrorWrapper")
	}
	return ws.ExitStatus()
}

// CodeExitError is an implementation of ExitError consisting of an error object
// and an exit code (the upper bits of os.exec.ExitStatus).
type CodeExitError struct {
	Err  error
	Code int
}

var _ ExitError = CodeExitError{}

func (e CodeExitError) Error() string {
	return e.Err.Error()
}

func (e CodeExitError) String() string {
	return e.Err.Error()
}

// Exited .
func (e CodeExitError) Exited() bool {
	return true
}

// ExitStatus .
func (e CodeExitError) ExitStatus() int {
	return e.Code
}
