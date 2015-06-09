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
	SetStatus(int) error
	GetStatus() int
	GetStatusText() string
}

type HandlerBase struct {
	status     int
	statusText string
	start      time.Time
	Log        *LogStack
	Child      IHandler

	clientIP string
	method   string
	path     string
}

func (self *HandlerBase) SetStatus(status int) error { self.status = status; return nil }
func (self *HandlerBase) SetStatusAndText(status int, text string) {
	self.status = status
	self.statusText = text
}
func (self *HandlerBase) GetStatus() int { return self.status }
func (self *HandlerBase) GetStatusText() string {
	if self.statusText == "" {
		return http.StatusText(self.GetStatus())
	}
	return self.statusText
}
func (self *HandlerBase) GetLog() *LogStack                 { return self.Log }
func (self *HandlerBase) SetClientIP(v string) *HandlerBase { self.clientIP = v; return self }
func (self *HandlerBase) GetClientIP() string               { return self.clientIP }
func (self *HandlerBase) SetMethod(v string) *HandlerBase   { self.method = v; return self }
func (self *HandlerBase) GetMethod() string                 { return self.method }
func (self *HandlerBase) SetPath(v string) *HandlerBase     { self.path = v; return self }
func (self *HandlerBase) GetPath() string                   { return self.path }

func (self *HandlerBase) GetHandlerName(handler interface{}) string {
	handlerFullName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	parts := strings.Split(handlerFullName, "/")
	return parts[len(parts)-1]
}

func (self *HandlerBase) serveHTTPPreffix(r *http.Request) {
	self.
		SetClientIP(r.RemoteAddr).
		SetMethod(r.Method).
		SetPath(r.URL.Path)
	self.Log.Debugf("Started %s \"%s\" for %s", self.GetMethod(), self.GetPath(), self.GetClientIP())
}

func (self *HandlerBase) serveHTTPSuffix(w http.ResponseWriter) {
	if self.GetStatus() != http.StatusOK && self.GetStatus() != http.StatusSwitchingProtocols {
		w.WriteHeader(self.GetStatus())
		fmt.Fprintln(w, self.GetStatusText())
	}

	latency := time.Now().Sub(self.start)
	msg := fmt.Sprintf("| %3d | %12v | %s | %s %-7s",
		self.GetStatus(),
		latency,
		self.GetClientIP(),
		self.GetMethod(), self.GetPath())

	status := self.GetStatus()
	switch {
	case status >= http.StatusInternalServerError:
		self.Log.Error(msg)
	case status >= http.StatusBadRequest:
		self.Log.Warnf(msg)
	default:
		self.Log.Infof(msg)
	}
}
