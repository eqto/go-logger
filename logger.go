package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type Format struct {
	value                   string
	level, date, time, file string
}

//Logger ...
type Logger struct {
	Level int

	File      string
	f         *os.File
	callDepth int

	formats map[int]*Format
	prefix  string
	out     io.Writer
}

//Print ...
func (l *Logger) Print(v ...interface{}) {
	l.print(LevelDebug, false, v...)
}

//Println ...
func (l *Logger) Println(v ...interface{}) {
	l.print(LevelDebug, true, v...)
}

//Format ...
func (l *Logger) Format(level int) *Format {
	if f, ok := l.formats[level]; ok {
		return f
	}
	return l.formats[LevelAll]
}

//SetLevelFormat ...
func (l *Logger) SetLevelFormat(level int, format string) {
	if l.formats == nil {
		l.formats = make(map[int]*Format)
	}
	l.formats[level] = newFormat(level, format)
	if l.formats[LevelAll] == nil {
		l.formats[LevelAll] = l.formats[level]
	}
}

//SetFormat ...
func (l *Logger) SetFormat(format string) {
	l.SetLevelFormat(LevelAll, format)
}

// Prefix returns the output prefix for the logger.
func (l *Logger) Prefix() string {
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

//D ...
func (l *Logger) D(v ...interface{}) {
	if l.Level <= LevelDebug {
		l.print(LevelDebug, true, v...)
	}
}

//I ...
func (l *Logger) I(v ...interface{}) {
	if l.Level <= LevelInfo {
		l.print(LevelInfo, true, v...)
	}
}

//W ...
func (l *Logger) W(v ...interface{}) {
	if l.Level <= LevelWarn {
		l.print(LevelWarn, true, v...)
	}
}

//E ...
func (l *Logger) E(v ...interface{}) {
	if l.Level <= LevelError {
		l.print(LevelError, true, v...)
	}
}

// F equivalent to Print() followed by a call to os.Exit(1).
func (l *Logger) F(v ...interface{}) {
	l.print(LevelFatal, false, v...)
	os.Exit(1)
}

//SetCallDepth ...
func (l *Logger) SetCallDepth(depth int) {
	l.callDepth = depth
}

//SetFile ...
func (l *Logger) SetFile(file string) {
	if file != `` {
		l.File = file
	}
}

func (l *Logger) println(level int, format string, v ...interface{}) {
	l.print(level, true, v...)
}

func (l *Logger) print(level int, newline bool, v ...interface{}) {
	var f *Format
	if format, ok := l.formats[level]; ok {
		f = format
	} else {
		f = l.formats[0]
	}
	out := string(f.value)
	now := time.Now()

	if f.level != `` {
		out = strings.Replace(
			out, f.level,
			levelColor(level)+strings.Replace(f.level, `%level%`, levelName[level], 1)+bgWhite+fgBlack, 1)
	}
	if f.date != `` {
		out = strings.Replace(
			out, f.date,
			bgWhite+fgBlack+strings.Replace(f.date, `%date%`, now.Format(`2006-01-02`), 1), 1)
	}
	if f.time != `` {
		out = strings.Replace(
			out, f.time,
			bgWhite+fgBlack+strings.Replace(f.time, `%time%`, now.Format(`15:04:05`), 1), 1)
	}
	if f.file != `` {
		_, file, line, _ := runtime.Caller(l.callDepth + 2)
		_, dir := path.Split(path.Dir(file))
		if dir == `runtime` {
			_, file, line, _ = runtime.Caller(l.callDepth + 1)
			_, dir = path.Split(path.Dir(file))
		}
		_, file = path.Split(file)
		dirs := regexDir.FindStringSubmatch(dir)
		dir = dirs[1]
		out = strings.Replace(
			out, f.file,
			bgWhite+fgCyan+strings.Replace(f.file, `%file%`, fmt.Sprintf(`%s/%s:%d`, dir, file, line), 1)+bgWhite+fgBlack, 1)
	}
	buf := strings.Builder{}
	buf.WriteString(out)

	if l.prefix != `` {
		buf.WriteString(l.prefix + ` `)
	}
	if newline {
		buf.WriteString(fmt.Sprintln(v...))
	} else {
		buf.WriteString(fmt.Sprint(v...))
	}
	if level >= LevelError {
		if !newline {
			buf.WriteString("\n")
		}
		frames := stacktrace(5)
		for _, frame := range frames {
			buf.WriteString(`    ` + frame.String() + "\n")
		}
	}
	if l.out == nil {
		l.out = os.Stderr
	}
	out = buf.String()
	l.out.Write([]byte(out))
	if l.f == nil && l.File != `` {
		if idx := strings.LastIndex(l.File, `/`); idx >= 0 {
			os.MkdirAll(l.File[0:strings.LastIndex(l.File, `/`)], 0755)
		}
		f, e := os.OpenFile(l.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if e != nil {
			l.W(e)
		} else {
			l.f = f
		}
	}
	if l.f != nil {
		l.f.WriteString(regexStrip.ReplaceAllString(out, ``))
	}
}

func stacktrace(skip int) []Frame {
	pc := make([]uintptr, 10+skip)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return nil
	}
	pc = pc[:n]
	i := runtime.CallersFrames(pc)
	frames := []Frame{}
	for {
		frame, more := i.Next()
		if !strings.HasPrefix(frame.Function, `runtime.`) &&
			!strings.HasPrefix(frame.Function, `reflect.Value.`) {
			frames = append(frames, newFrame(frame))
		}
		if !more {
			break
		}
	}
	return frames
}

//Stacktrace ...
func Stacktrace(skip int) []Frame {
	return stacktrace(3 + skip)
}

//D ...
func D(v ...interface{}) {
	std.D(v...)
}

//I ...
func I(v ...interface{}) {
	std.I(v...)
}

//W ...
func W(v ...interface{}) {
	std.W(v...)
}

//E ...
func E(v ...interface{}) {
	std.E(v...)
}

// F equivalent to Print() followed by a call to os.Exit(1).
func F(v ...interface{}) {
	std.F(v...)
}

//SetLevel ...
func SetLevel(level int) {
	std.Level = level
}

//SetFile ...
func SetFile(file string) {
	std.SetFile(file)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
// Compatibility for built-in go logging library
func Fatal(v ...interface{}) {
	std.Print(v...)
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
// Compatibility for built-in go logging library
func Fatalf(format string, v ...interface{}) {
	std.Print(fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is alias for F()
// Compatibility for built-in go logging library
func Fatalln(v ...interface{}) {
	std.F(v...)
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
// Compatibility for built-in go logging library
func Print(v ...interface{}) {
	std.Print(v...)
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
// Compatibility for built-in go logging library
func Printf(format string, v ...interface{}) {
	std.Print(fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
// Compatibility for built-in go logging library
func Println(v ...interface{}) {
	std.Print(fmt.Sprintln(v...))
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

//SetFormat ...
func SetFormat(format string) {
	std.SetFormat(format)
}

func newFormat(level int, format string) *Format {
	return &Format{
		value: format,
		level: regexLevel.FindString(format),
		date:  regexDate.FindString(format),
		time:  regexTime.FindString(format),
		file:  regexFile.FindString(format),
	}
}

//Default ...
func Default() *Logger {
	return std
}

//New ...
func New() *Logger {
	return NewWithFormat(DefaultFormat)
}

//NewWithFormat ...
func NewWithFormat(format string) *Logger {
	logger := &Logger{callDepth: 1}
	logger.SetFormat(format)
	return logger
}

//NewWithFile ...
func NewWithFile(filename string) *Logger {
	logger := NewWithFormat(DefaultFormat)
	logger.SetFile(filename)
	return logger
}
