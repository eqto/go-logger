/**
* Created by Visual Studio Code.
* User: tuxer
* Created At: 2019-11-19 00:09:59
**/

package log

import (
	"regexp"
)

var (
	levelColor = map[int]string{
		LevelDebug: bgWhite + fgGreen,
		LevelInfo:  bgWhite + fgBlue,
		LevelWarn:  bgWhite + fgYellow,
		LevelError: bgWhite + fgRed,
		LevelFatal: bgRed + fgWhite,
	}
	levelName = map[int]string{
		LevelDebug: `DEBUG`,
		LevelInfo:  `INFO `,
		LevelWarn:  `WARN `,
		LevelError: `ERROR`,
		LevelFatal: `FATAL`,
	}

	std        = NewDefault()
	regexLevel = regexp.MustCompile(`\S*%level%\S*`)
	regexDate  = regexp.MustCompile(`%date%`)
	regexTime  = regexp.MustCompile(`%time%`)
	regexFile  = regexp.MustCompile(`\S*%file%\S*`)

	regexStrip = regexp.MustCompile(`\033\[[0-9]+;1m`)
)
