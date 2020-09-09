package log

import (
	"regexp"
)

var (
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

	regexDir = regexp.MustCompile(`^(?Uis)(.+)(@.+|)$`)
)
