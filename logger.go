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

	File string
	f    *os.File

	prefix struct {
		value                   *string
		level, date, time, file string
	}

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

func (l *Logger) check() {
	if l.prefix.value == nil {
		l.SetFormat(defaultPrefix)
	}
	if l.File != `` && l.f == nil {
		os.MkdirAll(l.File[0:strings.LastIndex(l.File, `/`)], 0755)
		f, e := os.OpenFile(l.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if e != nil {
			W(e)
		} else {
			l.f = f
		}
	}
}

//SetFormat ...
func (l *Logger) SetFormat(format string) {
	l.prefix.value = &format
	l.prefix.level = regexLevel.FindString(format)
	l.prefix.date = regexDate.FindString(format)
	l.prefix.time = regexTime.FindString(format)
	l.prefix.file = regexFile.FindString(format)
}

//D ...
func (l *Logger) D(d ...interface{}) {
	if l.Level <= DEBUG {
		l.print(DEBUG, true, ``, d...)
	}
}

//I ...
func (l *Logger) I(i ...interface{}) {
	if l.Level <= INFO {
		l.print(INFO, true, ``, i...)
	}
}

//W ...
func (l *Logger) W(w ...interface{}) {
	if l.Level <= WARNING {
		l.print(WARNING, true, ``, w...)
	}
}

//E ...
func (l *Logger) E(e ...interface{}) {
	if l.Level <= ERROR {
		l.print(ERROR, true, ``, e...)
	}
}

// F equivalent to Print() followed by a call to os.Exit(1).
func (l *Logger) F(f ...interface{}) {
	l.print(FATAL, false, ``, f...)
	os.Exit(1)
}

func (l *Logger) print(level int, newLine bool, format string, v ...interface{}) {
	l.check()
	buffer := *l.prefix.value
	now := time.Now()

	if l.prefix.level != `` {
		buffer = strings.Replace(
			buffer, l.prefix.level,
			levelColor[level]+strings.Replace(l.prefix.level, `%level%`, levelName[level], 1)+bgWhite+fgBlack, 1)
	}
	if l.prefix.date != `` {
		buffer = strings.Replace(
			buffer, l.prefix.date,
			bgWhite+fgBlack+strings.Replace(l.prefix.date, `%date%`, now.Format(`2006-01-02`), 1), 1)
	}
	if l.prefix.time != `` {
		buffer = strings.Replace(
			buffer, l.prefix.time,
			bgWhite+fgBlack+strings.Replace(l.prefix.time, `%time%`, now.Format(`15:04:05`), 1), 1)
	}
	if l.prefix.file != `` {
		_, f, line, _ := runtime.Caller(3)
		_, dir := path.Split(path.Dir(f))
		_, f = path.Split(f)
		buffer = strings.Replace(
			buffer, l.prefix.file,
			bgWhite+fgCyan+strings.Replace(l.prefix.file, `%file%`, fmt.Sprintf(`%s/%s:%d`, dir, f, line), 1)+bgWhite+fgBlack, 1)
	}
	buffer = buffer + fmt.Sprint(v...)
	if level >= ERROR {
		buffer = buffer + "\n"
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
	if newLine && level < ERROR {
		l.out.Write([]byte("\n"))
	}
	if l.f != nil {
		l.f.WriteString(regexStrip.ReplaceAllString(buffer, ``))
	}
}

//D ...
func D(d ...interface{}) {
	std.D(d...)
}

//I ...
func I(i ...interface{}) {
	std.I(i...)
}

//W ...
func W(w ...interface{}) {
	std.W(w...)
}

//E ...
func E(e ...interface{}) {
	std.E(e...)
}

// F equivalent to Print() followed by a call to os.Exit(1).
func F(f ...interface{}) {
	std.F(f...)
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
	std.Println(fmt.Sprintln(v...))
}
