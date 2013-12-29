package absolut

import (
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
)

type logStack struct {
	sync.Mutex
	s [][2]string
}

func (self *logStack) Tracef(format string, params ...interface{}) {
	self.appendWithLevel("trace", format, params...)
}

func (self *logStack) Debugf(format string, params ...interface{}) {
	self.appendWithLevel("debug", format, params...)
}

func (self *logStack) Infof(format string, params ...interface{}) {
	self.appendWithLevel("info", format, params...)
}

func (self *logStack) Warnf(format string, params ...interface{}) {
	self.appendWithLevel("warn", format, params...)
}

func (self *logStack) Errorf(format string, params ...interface{}) {
	self.appendWithLevel("error", format, params...)
}

func (self *logStack) Criticalf(format string, params ...interface{}) {
	self.appendWithLevel("critical", format, params...)
}

func (self *logStack) appendWithLevel(level string, format string, params ...interface{}) {
	self.s = append(self.s, [2]string{level, fmt.Sprintf(format, params...)})
}

func (self *logStack) Flush() {
	self.Lock()
	for len(self.s) > 0 {
		level, message := self.s[0][0], self.s[0][1]
		switch level {
		case "trace":
			log.Trace(message)
		case "debug":
			log.Debug(message)
		case "info":
			log.Info(message)
		case "warn":
			log.Warn(message)
		case "error":
			log.Error(message)
		case "critical":
			log.Critical(message)
		}
		self.s = self.s[1:len(self.s)]
	}
	self.Unlock()
	log.Flush()
}
