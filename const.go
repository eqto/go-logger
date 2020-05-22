/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2019-11-19 00:11:07
**/

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
