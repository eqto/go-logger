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
