package log

import (
	"fmt"
	"strings"
	"time"
)

type Session struct {
	logger *Logger
	items  []sessionItem
}

func (s *Session) log(msg ...string) {
	s.items = append(s.items, sessionItem{time: time.Now(), msg: strings.Join(msg, ` `)})
}

func (s *Session) flush() {
	sb := strings.Builder{}
	var start time.Time
	var strTime string
	for idx, item := range s.items {
		if idx == 0 {
			start = item.time
			strTime = start.Format(`2006-01-02 15:04:05`)
		} else {
			strTime = item.time.Sub(start).String()
		}
		sb.WriteString(fmt.Sprintf("\t%-10s > %s\n", strTime, item.msg))
	}
	if s.logger != nil {
		s.logger.writer().Write([]byte(sb.String()))
	} else {
		std.Println(sb.String())
	}
}

type sessionItem struct {
	time time.Time
	msg  string
}

func Start(msg ...string) (func(...string), func()) {
	sess := Session{}

	if len(msg) > 0 {
		sess.log(msg...)
	}
	return sess.log, sess.flush
}
