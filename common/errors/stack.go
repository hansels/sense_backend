package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/hansels/sense_backend/common/bufpool"
)

// WithStack adds a stack to an error.
func WithStack(err error) error {
	return withStack(err, 2)
}

func withStack(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &stack{
		cause:   err,
		callers: callers(skip + 1),
	}
}

type stack struct {
	cause   error
	callers []uintptr
}

func (err *stack) StackFrames() *runtime.Frames {
	return runtime.CallersFrames(err.callers)
}

func (err *stack) Message() string {
	buf := bufpool.Get()
	buf.WriteString("stack")
	fs := err.StackFrames()
	for more := true; more; {
		var f runtime.Frame
		f, more = fs.Next()
		_, file := filepath.Split(f.File)
		_, _ = fmt.Fprintf(buf, "\n\t%s %s:%d", f.Function, file, f.Line)
	}
	s := buf.String()
	bufpool.Put(buf)
	return s
}

func (err *stack) Error() string { return Error(err) }
func (err *stack) Cause() error  { return err.cause }

// StackFrames returns the the list of runtime.Frames associated to an error.
func StackFrames(err error) []*runtime.Frames {
	var fss []*runtime.Frames
	for err := err; err != nil; err = Cause(err) {
		fs := getStackFrames(err)
		if fs != nil {
			fss = append(fss, fs)
		}
	}
	return fss
}

func getStackFrames(err error) *runtime.Frames {
	type stackFramer interface {
		StackFrames() *runtime.Frames
	}
	if err, ok := err.(stackFramer); ok {
		return err.StackFrames()
	}
	return nil
}

func ensureStack(err error, skip int) error {
	if !hasStack(err) {
		err = withStack(err, skip+1)
	}
	return err
}

func hasStack(err error) bool {
	type stackFramer interface {
		StackFrames() *runtime.Frames
	}
	for err := err; err != nil; err = Cause(err) {
		if _, ok := err.(stackFramer); ok {
			return true
		}
	}
	return false
}

const callersMaxLength = 1 << 16

var callersPool = sync.Pool{
	New: func() interface{} {
		return make([]uintptr, callersMaxLength)
	},
}

func callers(skip int) []uintptr {
	pcItf := callersPool.Get()
	pc := pcItf.([]uintptr)
	n := runtime.Callers(skip+1, pc)
	pcRes := make([]uintptr, n)
	copy(pcRes, pc)
	callersPool.Put(pcItf)
	return pcRes
}
