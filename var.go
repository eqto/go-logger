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
		DEBUG:   bgWhite + fgGreen,
		INFO:    bgWhite + fgBlue,
		WARNING: bgWhite + fgYellow,
		ERROR:   bgWhite + fgRed,
		FATAL:   bgRed + fgWhite,
	}
	levelName = map[int]string{
		DEBUG:   `DEBUG`,
		INFO:    `INFO `,
		WARNING: `WARN `,
		ERROR:   `ERROR`,
		FATAL:   `FATAL`,
	}

	std        = NewDefault()
	regexLevel = regexp.MustCompile(`\S*%level%\S*`)
	regexDate  = regexp.MustCompile(`%date%`)
	regexTime  = regexp.MustCompile(`%time%`)
	regexFile  = regexp.MustCompile(`\S*%file%\S*`)

	regexStrip = regexp.MustCompile(`\033\[[0-9]+;1m`)
)
