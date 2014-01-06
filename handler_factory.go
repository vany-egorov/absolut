package absolut

import (
	"net/http"
)

type FactoryHTTP struct {
	HandlerFunc HandlerHTTPFuncType
}

type FactoryWebsocket struct {
	Initializer WebsocketServerInitializer
}

func Δ(params ...interface{}) http.Handler {
	return NewHandlerFactory(params...)
}

func NewHandlerFactory(params ...interface{}) http.Handler {
	if handlerFunc, ok := params[0].(func(http.ResponseWriter, *http.Request, *HandlerHTTP) error); ok {
		return &FactoryHTTP{
			HandlerFunc: handlerFunc,
		}
	} else if initializer, ok := params[0].(WebsocketServerInitializer); ok {
		return &FactoryWebsocket{
			Initializer: initializer,
		}
	}
	return nil
}

func (self *FactoryHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewHandlerHTTP(self.HandlerFunc)
	handler.ServeHTTP(w, r)
}

func (self *FactoryWebsocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewHandlerWebsocket(self.WebsocketServerInitializer)
	handler.ServeHTTP(w, r)
}
