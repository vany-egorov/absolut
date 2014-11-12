package absolut

import (
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
)

type LogStack struct {
	sync.Mutex
	log log.LoggerInterface
	s   [][2]string
}

func LogStackNew(args ...interface{}) *LogStack {
	self := new(LogStack)

	if len(args) > 0 {
		if v, ok := args[0].(log.LoggerInterface); ok {
			self.log = v
		}
	} else {
		if defautLogger := GetDefaultLogger(); defautLogger != nil {
			self.log = defautLogger
		} else {
			self.log = log.Current
		}
	}

	return self
}

func (self *LogStack) Tracef(format string, params ...interface{}) *LogStack {
	self.appendWithLevel(log.TraceStr, format, params...)
	return self
}

func (self *LogStack) Debugf(format string, params ...interface{}) *LogStack {
	self.appendWithLevel(log.DebugStr, format, params...)
	return self
}

func (self *LogStack) Infof(format string, params ...interface{}) *LogStack {
	self.appendWithLevel(log.InfoStr, format, params...)
	return self
}

func (self *LogStack) Warnf(format string, params ...interface{}) *LogStack {
	self.appendWithLevel(log.WarnStr, format, params...)
	return self
}

func (self *LogStack) Errorf(format string, params ...interface{}) *LogStack {
	self.appendWithLevel(log.ErrorStr, format, params...)
	return self
}

func (self *LogStack) Criticalf(format string, params ...interface{}) *LogStack {
	self.appendWithLevel(log.CriticalStr, format, params...)
	return self
}

func (self *LogStack) appendWithLevel(level string, format string, params ...interface{}) {
	self.s = append(self.s, [2]string{level, fmt.Sprintf(format, params...)})
}

func (self *LogStack) Flush() *LogStack {
	self.Lock()
	for len(self.s) > 0 {
		level, message := self.s[0][0], self.s[0][1]
		switch level {
		case log.TraceStr:
			self.log.Trace(message)
		case log.DebugStr:
			self.log.Debug(message)
		case log.InfoStr:
			self.log.Info(message)
		case log.WarnStr:
			self.log.Warn(message)
		case log.ErrorStr:
			self.log.Error(message)
		case log.CriticalStr:
			self.log.Critical(message)
		}
		self.s = self.s[1:len(self.s)]
	}
	self.Unlock()
	self.log.Flush()

	return self
}
