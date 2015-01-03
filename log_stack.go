package absolut

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/cihub/seelog"
)

type LogStack struct {
	sync.Mutex
	log      log.LoggerInterface
	s        []logMesage
	isClosed bool
}

type logMesage struct {
	log.LoggerInterface

	level log.LogLevel
	text  string
}

func LogStackNew(args ...interface{}) *LogStack {
	self := new(LogStack)
	self.isClosed = false

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

func (self *LogStack) Tracef(format string, params ...interface{}) {
	self.appendWithLevelf(log.TraceLvl, format, params...)
}

func (self *LogStack) Debugf(format string, params ...interface{}) {
	self.appendWithLevelf(log.DebugLvl, format, params...)
}

func (self *LogStack) Infof(format string, params ...interface{}) {
	self.appendWithLevelf(log.InfoLvl, format, params...)
}

func (self *LogStack) Warnf(format string, params ...interface{}) error {
	return errors.New(self.appendWithLevelf(log.WarnLvl, format, params...).text)
}

func (self *LogStack) Errorf(format string, params ...interface{}) error {
	return errors.New(self.appendWithLevelf(log.ErrorLvl, format, params...).text)
}

func (self *LogStack) Criticalf(format string, params ...interface{}) error {
	return errors.New(self.appendWithLevelf(log.CriticalLvl, format, params...).text)
}

func (self *LogStack) Trace(v ...interface{}) { self.appendWithLevel(log.TraceLvl, v...) }
func (self *LogStack) Debug(v ...interface{}) { self.appendWithLevel(log.DebugLvl, v...) }
func (self *LogStack) Info(v ...interface{})  { self.appendWithLevel(log.InfoLvl, v...) }

func (self *LogStack) Warn(v ...interface{}) error {
	return errors.New(self.appendWithLevel(log.WarnLvl, v...).text)
}
func (self *LogStack) Error(v ...interface{}) error {
	return errors.New(self.appendWithLevel(log.ErrorLvl, v...).text)
}
func (self *LogStack) Critical(v ...interface{}) error {
	return errors.New(self.appendWithLevel(log.CriticalLvl, v...).text)
}

func (self *LogStack) traceWithCallDepth(callDepth int, message fmt.Stringer)    {}
func (self *LogStack) debugWithCallDepth(callDepth int, message fmt.Stringer)    {}
func (self *LogStack) infoWithCallDepth(callDepth int, message fmt.Stringer)     {}
func (self *LogStack) warnWithCallDepth(callDepth int, message fmt.Stringer)     {}
func (self *LogStack) errorWithCallDepth(callDepth int, message fmt.Stringer)    {}
func (self *LogStack) criticalWithCallDepth(callDepth int, message fmt.Stringer) {}

func (self *LogStack) appendWithLevel(level log.LogLevel, params ...interface{}) *logMesage {
	msg := logMesage{level: level, text: fmt.Sprint(params...)}
	self.s = append(self.s, msg)
	return &msg
}

func (self *LogStack) appendWithLevelf(level log.LogLevel, format string, params ...interface{}) *logMesage {
	msg := logMesage{level: level, text: fmt.Sprintf(format, params...)}
	self.s = append(self.s, msg)
	return &msg
}

func (self *LogStack) Close() {
	self.Flush()
	self.isClosed = true
}

func (self *LogStack) Flush() {
	self.Lock()
	defer self.Unlock()
	for len(self.s) > 0 {
		message := self.s[0]
		switch message.level {
		case log.TraceLvl:
			self.log.Trace(message.text)
		case log.DebugLvl:
			self.log.Debug(message.text)
		case log.InfoLvl:
			self.log.Info(message.text)
		case log.WarnLvl:
			self.log.Warn(message.text)
		case log.ErrorLvl:
			self.log.Error(message.text)
		case log.CriticalLvl:
			self.log.Critical(message.text)
		}
		self.s = self.s[1:len(self.s)]
	}
	self.log.Flush()
}

func (self *LogStack) Closed() bool { return self.isClosed }

func (self *LogStack) SetAdditionalStackDepth(depth int) error { return nil }
