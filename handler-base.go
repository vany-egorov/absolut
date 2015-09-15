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
	log        *LogStack
	Child      IHandler

	clientIP         string
	method           string
	path             string
	extension        Extension
	contentLengthIN  int64
	contentLengthOUT int

	isPoll       bool
	pollStatuses map[int]bool
}

func (self *HandlerBase) SetStatus(status int) error { self.status = status; return nil }
func (self *HandlerBase) SetStatusAndText(status int, text string) error {
	self.status = status
	self.statusText = text
	return nil
}
func (self *HandlerBase) GetStatus() int { return self.status }
func (self *HandlerBase) GetStatusText() string {
	if self.statusText == "" {
		return http.StatusText(self.GetStatus())
	}
	return self.statusText
}
func (self *HandlerBase) GetExtension() Extension           { return self.GetExt() }
func (self *HandlerBase) GetExt() Extension                 { return self.extension }
func (self *HandlerBase) SetExt(v Extension) *HandlerBase   { self.extension = v; return self }
func (self *HandlerBase) Ext() Extension                    { return self.extension }
func (self *HandlerBase) Log() *LogStack                    { return self.log }
func (self *HandlerBase) GetLog() *LogStack                 { return self.log }
func (self *HandlerBase) SetClientIP(v string) *HandlerBase { self.clientIP = v; return self }
func (self *HandlerBase) GetClientIP() string               { return self.clientIP }
func (self *HandlerBase) SetMethod(v string) *HandlerBase   { self.method = v; return self }
func (self *HandlerBase) GetMethod() string                 { return self.method }
func (self *HandlerBase) SetPath(v string) *HandlerBase     { self.path = v; return self }
func (self *HandlerBase) GetPath() string                   { return self.path }
func (self *HandlerBase) SetContentLengthOUT(v int) *HandlerBase {
	self.contentLengthOUT = v
	return self
}
func (self *HandlerBase) SetIsPoll() { self.isPoll = true }
func (self *HandlerBase) SetPollStatuses(vs []int) *HandlerBase {
	self.pollStatuses = make(map[int]bool)
	for _, v := range vs {
		self.pollStatuses[v] = true
	}
	return self
}

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
}

func (self *HandlerBase) serveHTTPSuffix(w http.ResponseWriter) {
	if self.GetStatus() != http.StatusOK && self.GetStatus() != http.StatusSwitchingProtocols {
		w.WriteHeader(self.GetStatus())
		n, _ := fmt.Fprintln(w, self.GetStatusText())
		self.SetContentLengthOUT(n)
	}

	latency := time.Now().Sub(self.start)
	msg := fmt.Sprintf("| %3d | %12v | %s | %s %-7s | %s | <- %d | -> %d",
		self.GetStatus(),
		latency,
		self.GetClientIP(),
		self.GetMethod(), self.GetPath(),
		ExtensionText(self.Ext()),
		self.contentLengthIN,
		self.contentLengthOUT)

	status := self.GetStatus()
	switch {
	case status >= http.StatusInternalServerError:
		self.Log().Error(msg)
	case status >= http.StatusBadRequest:
		if ok := self.pollStatuses[status]; self.isPoll && ok {
			self.Log().Debug(msg)
		} else {
			self.Log().Warn(msg)
		}
	default:
		if ok := self.pollStatuses[status]; self.isPoll && ok {
			self.Log().Debug(msg)
		} else {
			self.Log().Info(msg)
		}
	}
}
