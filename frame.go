package log

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type Frame struct {
	File     string
	Dir      string
	Version  string
	Function string
	Line     int
}

func (f *Frame) String() string {
	return fmt.Sprintf("(%s/%s:%d) %s", f.Dir, f.File, f.Line, f.Function)
}

func newFrame(frame runtime.Frame) Frame {
	dir, file := path.Split(frame.File)
	dirs := strings.SplitN(path.Base(dir), `@`, 2)
	f := Frame{
		File:     file,
		Dir:      dirs[0],
		Function: frame.Function,
		Line:     frame.Line,
	}
	if len(dirs) == 2 {
		f.Version = dirs[1]
	}
	return f
}
