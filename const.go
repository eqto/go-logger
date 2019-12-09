/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2019-11-19 00:11:07
**/

package log

const (
	DefaultFormat = `[%level%] %date% %time% (%file%) `
	DefaultFile   = `log/app.log`
)

const (
	//DEBUG ...
	DEBUG = iota
	//INFO ...
	INFO
	//WARNING ...
	WARNING
	//ERROR ...
	ERROR
	//FATAL ...
	FATAL
)
