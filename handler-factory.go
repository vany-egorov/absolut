package absolut

import (
	"net/http"
	"time"
)

type FactoryHTTP struct {
	HandlerFunc  HandlerHTTPFuncType
	LoggerGetter LoggerGetter
}

type FactoryWebsocket struct {
	Initializer WebsocketServerInitializer
	ReadWait    time.Duration
}

func Δ(params ...interface{}) http.Handler {
	return NewHandlerFactory(params...)
}

func NewHandlerFactory(params ...interface{}) http.Handler {
	if handlerFunc, ok := params[0].(func(http.ResponseWriter, *http.Request, *HandlerHTTP) error); ok {
		loggerGetter := _defaultLoggerGetter
		if len(params) > 1 {
			loggerGetter = params[1].(LoggerGetter)

		}
		return &FactoryHTTP{
			HandlerFunc:  handlerFunc,
			LoggerGetter: loggerGetter,
		}
	} else if initializer, ok := params[0].(WebsocketServerInitializer); ok {
		var readWait time.Duration
		if len(params) > 1 {
			readWait = params[1].(time.Duration)
			readWait = readWait
		} else {
			readWait = 1 * time.Second
		}
		return &FactoryWebsocket{
			Initializer: initializer,
			ReadWait:    readWait,
		}
	}
	return nil
}

func (self *FactoryHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewHandlerHTTP(self.HandlerFunc, self.LoggerGetter)
	handler.contentLengthIN = r.ContentLength
	handler.ServeHTTP(w, r)
}

func (self *FactoryWebsocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewHandlerWebsocket(self.Initializer, self.ReadWait)
	handler.contentLengthIN = r.ContentLength
	handler.ServeHTTP(w, r)
}
