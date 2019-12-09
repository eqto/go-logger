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

/**
 * Created by tuxer on 8/30/17.
 */

//Logger ...
type Logger struct {
	Level int

	File      string
	f         *os.File
	CallDepth int

	format struct {
		value                   *string
		level, date, time, file string
	}

	prefix string

	out io.Writer
}

//Print ...
func (l *Logger) Print(v ...interface{}) {
	l.print(DEBUG, false, ``, v...)
}

//Printf ...
func (l *Logger) Printf(format string, v ...interface{}) {
	l.print(DEBUG, false, format, v...)
}

//Println ...
func (l *Logger) Println(v ...interface{}) {
	l.print(DEBUG, true, ``, v...)
}

//SetFormat ...
func (l *Logger) SetFormat(format string) {
	l.format.value = &format
	l.format.level = regexLevel.FindString(format)
	l.format.date = regexDate.FindString(format)
	l.format.time = regexTime.FindString(format)
	l.format.file = regexFile.FindString(format)
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
	if l.Level <= DEBUG {
		l.print(DEBUG, true, ``, v...)
	}
}

//I ...
func (l *Logger) I(v ...interface{}) {
	if l.Level <= INFO {
		l.print(INFO, true, ``, v...)
	}
}

//W ...
func (l *Logger) W(v ...interface{}) {
	if l.Level <= WARNING {
		l.print(WARNING, true, ``, v...)
	}
}

//E ...
func (l *Logger) E(v ...interface{}) {
	if l.Level <= ERROR {
		l.print(ERROR, true, ``, v...)
	}
}

// F equivalent to Print() followed by a call to os.Exit(1).
func (l *Logger) F(v ...interface{}) {
	l.print(FATAL, false, ``, v...)
	os.Exit(1)
}

func (l *Logger) println(level int, format string, v ...interface{}) {
	l.print(level, true, format, v...)
}

func (l *Logger) print(level int, newline bool, format string, v ...interface{}) {
	buffer := *l.format.value
	now := time.Now()

	if l.format.level != `` {
		buffer = strings.Replace(
			buffer, l.format.level,
			levelColor[level]+strings.Replace(l.format.level, `%level%`, levelName[level], 1)+bgWhite+fgBlack, 1)
	}
	if l.format.date != `` {
		buffer = strings.Replace(
			buffer, l.format.date,
			bgWhite+fgBlack+strings.Replace(l.format.date, `%date%`, now.Format(`2006-01-02`), 1), 1)
	}
	if l.format.time != `` {
		buffer = strings.Replace(
			buffer, l.format.time,
			bgWhite+fgBlack+strings.Replace(l.format.time, `%time%`, now.Format(`15:04:05`), 1), 1)
	}
	if l.format.file != `` {
		_, f, line, _ := runtime.Caller(l.CallDepth + 2)
		_, dir := path.Split(path.Dir(f))
		_, f = path.Split(f)
		buffer = strings.Replace(
			buffer, l.format.file,
			bgWhite+fgCyan+strings.Replace(l.format.file, `%file%`, fmt.Sprintf(`%s/%s:%d`, dir, f, line), 1)+bgWhite+fgBlack, 1)
	}
	if l.prefix != `` {
		buffer = buffer + l.prefix + ` `
	}
	if newline {
		buffer = buffer + fmt.Sprintln(v...)
	} else {
		buffer = buffer + fmt.Sprint(v...)
	}
	if level >= ERROR {
		if !newline {
			buffer = buffer + "\n"
		}
		pc := make([]uintptr, 10)
		runtime.Callers(5, pc)
		for _, p := range pc {
			if p > 0 {
				f := runtime.FuncForPC(p)
				file, line := f.FileLine(p)
				name := f.Name()
				if !strings.HasPrefix(name, `runtime.`) && !strings.HasPrefix(name, `reflect.Value.`) {
					_, dir := path.Split(path.Dir(file))
					_, file = path.Split(file)
					buffer = buffer + `  ` + fmt.Sprintf(`(%s:%d) %s`, dir+`/`+file, line, f.Name()) + "\n"
				}
			}

		}
	}
	if l.out == nil {
		l.out = os.Stdout
	}
	l.out.Write([]byte(buffer))
	if l.f != nil {
		l.f.WriteString(regexStrip.ReplaceAllString(buffer, ``))
	}
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
	std.File = file
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
	std.Printf(format, v...)
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
	std.Print(fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
// Compatibility for built-in go logging library
func Printf(format string, v ...interface{}) {
	std.Printf(fmt.Sprintf(format, v...))
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
	std.format.value = &format
	std.format.level = regexLevel.FindString(format)
	std.format.date = regexDate.FindString(format)
	std.format.time = regexTime.FindString(format)
	std.format.file = regexFile.FindString(format)
}

//New ...
func New(format string, file string) *Logger {
	logger := &Logger{
		CallDepth: 1,
		File:      file,
	}
	if file != `` {
		os.MkdirAll(file[0:strings.LastIndex(file, `/`)], 0755)
		f, e := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if e != nil {
			logger.W(e)
		} else {
			logger.f = f
		}
	}
	logger.SetFormat(format)
	return logger
}

//Default ...
func Default() *Logger {
	return std
}
