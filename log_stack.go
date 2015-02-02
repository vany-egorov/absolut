package absolut

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cihub/seelog"
)

type loggerGetter interface {
	GetLogger() seelog.LoggerInterface
}

type LogStack struct {
	sync.Mutex

	log          seelog.LoggerInterface
	loggerGetter loggerGetter

	s        []logMesage
	isClosed bool
}

type logMesage struct {
	level seelog.LogLevel
	text  string
}

func LogStackNew(args ...interface{}) *LogStack {
	self := new(LogStack)
	self.isClosed = false

	if len(args) > 0 {
		if v, ok := args[0].(seelog.LoggerInterface); ok {
			self.log = v
		}
		if v, ok := args[0].(loggerGetter); ok {
			self.loggerGetter = v
		}
	} else {
		if defautLogger := GetDefaultLogger(); defautLogger != nil {
			self.log = defautLogger
		} else {
			if defaultLoggerGetter := GetDefaultLoggerGetter(); defaultLoggerGetter != nil {
				self.log = defaultLoggerGetter.GetLogger()
			} else {
				self.log = seelog.Current
			}
		}
	}

	return self
}

func (self *LogStack) Tracef(format string, params ...interface{}) {
	self.appendWithLevelf(seelog.TraceLvl, format, params...)
}

func (self *LogStack) Debugf(format string, params ...interface{}) {
	self.appendWithLevelf(seelog.DebugLvl, format, params...)
}

func (self *LogStack) Infof(format string, params ...interface{}) {
	self.appendWithLevelf(seelog.InfoLvl, format, params...)
}

func (self *LogStack) Warnf(format string, params ...interface{}) error {
	return errors.New(self.appendWithLevelf(seelog.WarnLvl, format, params...).text)
}

func (self *LogStack) Errorf(format string, params ...interface{}) error {
	return errors.New(self.appendWithLevelf(seelog.ErrorLvl, format, params...).text)
}

func (self *LogStack) Criticalf(format string, params ...interface{}) error {
	return errors.New(self.appendWithLevelf(seelog.CriticalLvl, format, params...).text)
}

func (self *LogStack) Trace(v ...interface{}) { self.appendWithLevel(seelog.TraceLvl, v...) }
func (self *LogStack) Debug(v ...interface{}) { self.appendWithLevel(seelog.DebugLvl, v...) }
func (self *LogStack) Info(v ...interface{})  { self.appendWithLevel(seelog.InfoLvl, v...) }

func (self *LogStack) Warn(v ...interface{}) error {
	return errors.New(self.appendWithLevel(seelog.WarnLvl, v...).text)
}
func (self *LogStack) Error(v ...interface{}) error {
	return errors.New(self.appendWithLevel(seelog.ErrorLvl, v...).text)
}
func (self *LogStack) Critical(v ...interface{}) error {
	return errors.New(self.appendWithLevel(seelog.CriticalLvl, v...).text)
}

func (self *LogStack) traceWithCallDepth(callDepth int, message fmt.Stringer)    {}
func (self *LogStack) debugWithCallDepth(callDepth int, message fmt.Stringer)    {}
func (self *LogStack) infoWithCallDepth(callDepth int, message fmt.Stringer)     {}
func (self *LogStack) warnWithCallDepth(callDepth int, message fmt.Stringer)     {}
func (self *LogStack) errorWithCallDepth(callDepth int, message fmt.Stringer)    {}
func (self *LogStack) criticalWithCallDepth(callDepth int, message fmt.Stringer) {}

func (self *LogStack) appendWithLevel(level seelog.LogLevel, params ...interface{}) *logMesage {
	msg := logMesage{level: level, text: fmt.Sprint(params...)}
	self.s = append(self.s, msg)
	return &msg
}

func (self *LogStack) appendWithLevelf(level seelog.LogLevel, format string, params ...interface{}) *logMesage {
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

	var log seelog.LoggerInterface
	if self.log != nil {
		log = self.log
	}
	if self.loggerGetter != nil {
		log = self.loggerGetter.GetLogger()
	}

	for len(self.s) > 0 {
		message := self.s[0]
		switch message.level {
		case seelog.TraceLvl:
			log.Trace(message.text)
		case seelog.DebugLvl:
			log.Debug(message.text)
		case seelog.InfoLvl:
			log.Info(message.text)
		case seelog.WarnLvl:
			log.Warn(message.text)
		case seelog.ErrorLvl:
			log.Error(message.text)
		case seelog.CriticalLvl:
			log.Critical(message.text)
		}
		self.s = self.s[1:len(self.s)]
	}
	log.Flush()
}

func (self *LogStack) Closed() bool { return self.isClosed }

func (self *LogStack) SetAdditionalStackDepth(depth int) error { return nil }
