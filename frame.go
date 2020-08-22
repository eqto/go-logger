package log

import (
	"fmt"
	"path"
	"runtime"
)

//Frame ...
type Frame struct {
	File     string
	Dir      string
	Function string
	Line     int
}

func (f *Frame) String() string {
	return fmt.Sprintf("(%s/%s:%d) %s", f.Dir, f.File, f.Line, f.Function)
}

func newFrame(frame runtime.Frame) Frame {
	dir, file := path.Split(frame.File)
	dir = path.Base(dir)
	return Frame{
		File:     file,
		Dir:      dir,
		Function: frame.Function,
		Line:     frame.Line,
	}
}
