package absolut

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type BaseHandler struct {
	Handler BaseHandlerFunc
	status  int
	Log     *logStack
}

type BaseHandlerFunc func(http.ResponseWriter, *http.Request, *BaseHandler)

func NewBaseHandler(handler BaseHandlerFunc) *BaseHandler {
	return &BaseHandler{
		Handler: handler,
		status:  http.StatusOK,
		Log:     &logStack{},
	}
}

func (self *BaseHandler) getHandlerName() string {
	return runtime.FuncForPC(reflect.ValueOf(self.Handler).Pointer()).Name()
}

func (self *BaseHandler) SetStatus(status int) {
	self.status = status
}

func (self *BaseHandler) GetStatus() int {
	return self.status
}

func (self *BaseHandler) GetStatusText() string {
	return http.StatusText(self.GetStatus())
}

func (self *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		self.Log.Flush()
	}()

	start := time.Now()
	self.Log.Infof("Started %s \"%s\" for %s", r.Method, r.URL.Path, r.RemoteAddr)
	self.Log.Infof("\tProcessing by %s as HTML", self.getHandlerName())

	self.Handler(w, r, self)

	if self.GetStatus() != http.StatusOK {
		http.Error(w, self.GetStatusText(), self.GetStatus())
	}

	msg := fmt.Sprintf(
		"Completed (%d - %s) in %f ms\n",
		self.GetStatus(),
		self.GetStatusText(),
		time.Since(start).Seconds()*1000,
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
