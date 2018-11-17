package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

/**
 * Created by tuxer on 8/30/17.
 */

var (
	defaultFile = `log/application.log`
	logger      *Logger
)

const (
	//DEBUG ...
	DEBUG = iota
	//WARNING ...
	WARNING
	//INFO ...
	INFO
	//ERROR ...
	ERROR
	//FATAL ...
	FATAL
)

//Logger ...
type Logger struct {
	consoleLogger *log.Logger
	fileLogger    *log.Logger
	generalWriter io.Writer

	errorStyle, warningStyle, debugStyle, infoStyle, fatalStyle *styling

	File            string
	IncludeFilename bool
	Level           int
}

type styling struct {
	prepend string
	color   string
}

//SetDefaultFile ...
func SetDefaultFile(file string) {
	DefaultLogger().File = file
}

//SetLevel ...
func SetLevel(level int) {
	DefaultLogger().Level = level
}

//DefaultLogger ...
func DefaultLogger() *Logger {
	if logger == nil {
		l := Logger{File: defaultFile, IncludeFilename: true, Level: DEBUG}
		logger = &l
	}
	return logger
}

func (l *Logger) createFileLogger(name string) *log.Logger {
	if l.File == `` {
		l.File = defaultFile
	}

	os.MkdirAll(l.File[0:strings.LastIndex(l.File, `/`)], 0755)
	f, e := os.OpenFile(l.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if e != nil {
		log.Fatal(e)
		return nil
	}
	return log.New(f, ``, 0)
}

//W ...
func (l *Logger) W(warnings ...interface{}) {
	l.initLogger()
	if l.Level <= WARNING {
		l.printLog(l.warningStyle, false, warnings...)
	}
}

//W ...
func W(warnings ...interface{}) {
	DefaultLogger().W(warnings...)
}

//E ...
func (l *Logger) E(errors ...interface{}) {
	l.initLogger()
	if l.Level <= ERROR {
		l.printLog(l.errorStyle, true, errors...)
	}
}

//E ...
func E(errors ...interface{}) {
	DefaultLogger().E(errors...)
}

//D ...
func (l *Logger) D(debugs ...interface{}) {
	l.initLogger()
	if l.Level <= DEBUG {
		l.printLog(l.debugStyle, false, debugs...)
	}
}

//D ...
func D(debugs ...interface{}) {
	DefaultLogger().D(debugs...)
}

//I ...
func (l *Logger) I(infos ...interface{}) {
	l.initLogger()
	if l.Level <= INFO {
		l.printLog(l.infoStyle, false, infos...)
	}
}

//I ...
func I(infos ...interface{}) {
	DefaultLogger().I(infos...)
}

//F ...
func (l *Logger) F(fatal ...interface{}) {
	l.initLogger()

	if l.Level <= FATAL {
		fatal = append(fatal, bgWhite)
		l.printLog(l.fatalStyle, true, fatal...)
		log.Fatalln()
	}
}

//F ...
func F(fatals ...interface{}) {
	DefaultLogger().F(fatals...)
}

func (l *Logger) initLogger() {
	if l.fileLogger == nil {
		l.fatalStyle = &styling{prepend: `[FATAL]`, color: bgRed + fgWhite}
		l.errorStyle = &styling{prepend: `[ERROR]`, color: fgRed}
		l.infoStyle = &styling{prepend: `[INFO ]`, color: fgBlue}
		l.warningStyle = &styling{prepend: `[WARN ]`, color: fgRed}
		l.debugStyle = &styling{prepend: `[DEBUG]`, color: fgYellow}
		l.fileLogger = l.createFileLogger(`application.log`)
	}
}

func (l *Logger) printLog(style *styling, withStack bool, obj ...interface{}) {
	if l.consoleLogger == nil {
		l.consoleLogger = log.New(os.Stdout, ``, 0)
	}

	date := time.Now().Format(` 2006-01-02 15:04:05 `)

	file := fgBlack + `:`
	if l.IncludeFilename {
		_, f, line, _ := runtime.Caller(3)
		_, dir := path.Split(path.Dir(f))
		_, f = path.Split(f)
		file = fmt.Sprintf(`(%s/%s:%d)`, dir, f, line)
	}

	console := append([]interface{}{style.color + style.prepend + fgBlack + date + fgCyan + file + fgBlack}, obj...)
	l.consoleLogger.Println(console...)
	console = append([]interface{}{style.prepend + date + file}, obj...)
	l.fileLogger.Println(console...)

	if withStack {
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

					formatted := fmt.Sprintf(`(%s:%d) %s`, dir+`/`+file, line, f.Name())
					l.consoleLogger.Println(formatted)
					l.fileLogger.Println(formatted)
				}
			}

		}
	}
}
