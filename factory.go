package absolut

import (
	"net/http"
)

type BaseHandlerFactory struct {
	Handler BaseHandlerFunc
}

func Δ(handler BaseHandlerFunc) *BaseHandlerFactory {
	return NewBaseHandlerFactory(handler)
}

func NewBaseHandlerFactory(handler BaseHandlerFunc) *BaseHandlerFactory {
	return &BaseHandlerFactory{
		Handler: handler,
	}
}

func (self *BaseHandlerFactory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := NewBaseHandler(self.Handler)
	handler.ServeHTTP(w, r)
}
