package absolut

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type IHandler interface {
	SetStatus(int)
	GetStatus() int
	GetStatusText() string
}

type HandlerBase struct {
	status int
	start  time.Time
	Log    *LogStack
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

func (self *HandlerBase) GetLog() *LogStack {
	return self.Log
}

func (self *HandlerBase) GetHandlerName(handler interface{}) string {
	handlerFullName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	parts := strings.Split(handlerFullName, "/")
	return parts[len(parts)-1]
}

func (self *HandlerBase) serveHTTPPreffix(r *http.Request) {
	self.Log.Infof("Started %s \"%s\" for %s", r.Method, r.URL.Path, r.RemoteAddr)
}

func (self *HandlerBase) serveHTTPSuffix(w http.ResponseWriter) {
	if self.GetStatus() != http.StatusOK && self.GetStatus() != http.StatusSwitchingProtocols {
		http.Error(w, self.GetStatusText(), self.GetStatus())
	}

	msg := fmt.Sprintf(
		"Completed (%d - %s) in %f ms\n",
		self.GetStatus(),
		self.GetStatusText(),
		time.Since(self.start).Seconds()*1000,
	)

	status := self.GetStatus()
	switch {
	case status >= http.StatusInternalServerError:
		self.Log.Errorf("%s", msg)
	case status >= http.StatusBadRequest:
		self.Log.Warnf("%s", msg)
	default:
		self.Log.Infof("%s", msg)
	}
}
