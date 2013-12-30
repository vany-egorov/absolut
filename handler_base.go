package absolut

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type IHandler interface {
	SetStatus(int)
	GetStatus() int
	GetStatusText() string
}

type HandlerBase struct {
	status int
	Log    *logStack
	Child  IHandler
}

func (self *HandlerBase) SetStatus(status int) {
	self.status = status
}

func (self *HandlerBase) GetStatus() int {
	return self.status
}

func (self *HandlerBase) GetStatusText() string {
	return http.StatusText(self.GetStatus())
}

func (self *HandlerBase) getHandlerName(handler interface{}) string {
	handlerFullName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	parts := strings.Split(handlerFullName, "/")
	return parts[len(parts)-1]
}
