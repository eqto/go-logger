package log

import (
	"log"
	"os"
	"runtime"
	"io"
    "time"
    "path"
    "fmt"
    "strings"
)

/**
 * Created by tuxer on 8/30/17.
 */

var (
    defaultPath = `log/`
	logger      *Logger
)

//Logger ...
type Logger struct {
	consoleLogger, errorLogger, warningLogger, debugLogger, infoLogger, fatalLogger *log.Logger
	generalWriter	io.Writer
    errorStyle, warningStyle, debugStyle, infoStyle, fatalStyle *styling

	Path            	string
	IncludeFilename		bool
}

type styling struct {
    prepend string
    color   string
}

//SetDefaultPath ...
func SetDefaultPath(path string)	{
	defaultPath = path
}

//DefaultLogger ...
func DefaultLogger() *Logger {
    if logger == nil    {
        l := Logger{Path: defaultPath, IncludeFilename: true}
        logger = &l
    }
    return logger
}

func (l *Logger) createFileLogger(name string) *log.Logger  {
    if l.Path == ``  {
        l.Path = defaultPath
    }
    os.MkdirAll(l.Path, 0755)
    f, e := os.OpenFile(l.Path + name, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
    if e != nil	{
        log.Fatal(e)
        return nil
    }
    return log.New(f, ``, 0)
}

//W ...
func (l *Logger) W(warnings...interface{})	{
    if l.warningLogger == nil  {
        style := styling{
            prepend: `[WARN ]`, color: fgRed,
        }
        l.warningStyle = &style
        l.warningLogger = l.createFileLogger(`warning.log`)
    }
    l.printLog(l.warningLogger, l.warningStyle, false, warnings...)
}

//W ...
func W(warnings...interface{})	{
	DefaultLogger().W(warnings...)
}

//E ...
func (l *Logger) E(errors...interface{})	{
    if l.errorLogger == nil  {
        style := styling{
            prepend: `[ERROR]`, color: fgRed,
        }
        l.errorStyle = &style
        l.errorLogger = l.createFileLogger(`error.log`)
    }
    l.printLog(l.errorLogger, l.errorStyle, true, errors...)
}

//E ...
func E(errors...interface{})	{
	DefaultLogger().E(errors...)
}

//D ...
func (l *Logger) D(debugs...interface{})	{
    if l.debugLogger == nil  {
        style := styling{
            prepend: `[DEBUG]`, color: fgYellow,
        }
        l.debugStyle = &style
        l.debugLogger = l.createFileLogger(`debug.log`)
    }
    l.printLog(l.debugLogger, l.debugStyle, false, debugs...)
}

//D ...
func D(debugs...interface{})	{
	DefaultLogger().D(debugs...)
}

//I ...
func (l *Logger) I(infos...interface{})	{
    if l.infoLogger == nil  {
        style := styling{
            prepend: `[INFO ]`, color: fgBlue,
        }
        l.infoStyle = &style
        l.infoLogger = l.createFileLogger(`info.log`)
	}
    l.printLog(l.infoLogger, l.infoStyle, false, infos...)
}

//I ...
func I(infos...interface{})	{
	DefaultLogger().I(infos...)
}

//F ...
func (l *Logger) F(fatal...interface{})	{
    if l.fatalLogger == nil  {
        style := styling{
            prepend: `[FATAL]`, color: bgRed + fgWhite,
        }
        l.fatalStyle = &style
        l.fatalLogger = l.createFileLogger(`fatal.log`)
	}
	fatal = append(fatal, bgWhite)
    l.printLog(l.fatalLogger, l.fatalStyle, true, fatal...)
	log.Fatalln()
}

//F ...
func F(fatals...interface{})	{
	DefaultLogger().F(fatals...)
}

func (l *Logger) printLog(fileLogger *log.Logger, style *styling, withStack bool, obj...interface{})	{
    if l.consoleLogger == nil   {
        l.consoleLogger = log.New(os.Stdout, ``, 0)
    }

    date := time.Now().Format(` 2006-01-02 15:04:05 `)

	file := fgBlack + `:`
	if l.IncludeFilename	{
		_, f, line, _ := runtime.Caller(3)
		_, dir := path.Split(path.Dir(f))
		_, f = path.Split(f)
		file = fmt.Sprintf(`(%s/%s:%d)`, dir, f, line)
	}

	console := append([]interface{}{style.color + style.prepend + fgBlack + date + fgCyan + file + fgBlack}, obj...)
    l.consoleLogger.Println(console...)
    console = append([]interface{}{style.prepend + date + file}, obj...)
    fileLogger.Println(console...)

    if withStack    {
        pc := make([]uintptr, 10)
        runtime.Callers(5, pc)
        for _, p := range pc  {
            if p > 0    {
                f := runtime.FuncForPC(p)
                file, line := f.FileLine(p)
                name := f.Name()
                if !strings.HasPrefix(name, `runtime.`) && !strings.HasPrefix(name, `reflect.Value.`)  {
                    _, dir := path.Split(path.Dir(file))
					_, file = path.Split(file)
					
					formatted := fmt.Sprintf(`(%s:%d) %s`, dir + `/` + file, line, f.Name())
                    l.consoleLogger.Println(formatted)
                    fileLogger.Println(formatted)
                }
            }

        }
    }
}

