package absolut

import (
	log "github.com/cihub/seelog"
)

var (
	_websocketReadBufferSize  int                 = 1024
	_websocketWriteBufferSize int                 = 1024
	_defaultLogger            log.LoggerInterface = nil
	_defaultLoggerGetter      loggerGetter        = nil
)

func SetWebsocketReadBufferSize(v int) { _websocketReadBufferSize = v }
func GetWebsocketReadBufferSize() int  { return _websocketReadBufferSize }

func SetWebsocketWriteBufferSize(v int) { _websocketWriteBufferSize = v }
func GetWebsocketWriteBufferSize() int  { return _websocketWriteBufferSize }

func SetDefaultLogger(v log.LoggerInterface) { _defaultLogger = v }
func GetDefaultLogger() log.LoggerInterface  { return _defaultLogger }

func SetDefaultLoggerGetter(v loggerGetter) { _defaultLoggerGetter = v }
func GetDefaultLoggerGetter() loggerGetter  { return _defaultLoggerGetter }
