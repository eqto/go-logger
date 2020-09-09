package log

const (
	DefaultFormat = `[%level%] %date% %time% (%file%) `
	DefaultFile   = `logs/application.log`
)

const (
	//LevelAll ...
	LevelAll = iota
	//LevelDebug ...
	LevelDebug
	//LevelInfo ...
	LevelInfo
	//LevelWarn ...
	LevelWarn
	//LevelError ...
	LevelError
	//LevelFatal ...
	LevelFatal
)

const (
	colorDebug = bgWhite + fgGreen
	colorInfo  = bgWhite + fgBlue
	colorWarn  = bgWhite + fgYellow
	colorError = bgWhite + fgRed
	colorFatal = bgRed + fgWhite
)

func levelColor(level int) string {
	switch level {
	case LevelDebug:
		return colorDebug
	case LevelInfo:
		return colorInfo
	case LevelWarn:
		return colorWarn
	case LevelError:
		return colorError
	case LevelFatal:
		return colorFatal
	}
	return bgWhite + fgWhite
}
