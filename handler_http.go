package absolut

import (
	"fmt"
	"net/http"
	"time"
)

type HandlerHTTP struct {
	HandlerBase
	HandlerFunc HandlerHTTPFuncType
}

type HandlerHTTPFuncType func(http.ResponseWriter, *http.Request, *HandlerHTTP) error

func NewHandlerHTTP(handler HandlerHTTPFuncType) *HandlerHTTP {
	self := &HandlerHTTP{
		HandlerBase: HandlerBase{
			status: http.StatusOK,
			Log:    &logStack{},
		},
		HandlerFunc: handler,
	}
	self.Child = self

	return self
}

func (self *HandlerHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		self.Log.Flush()
	}()

	start := time.Now()
	self.Log.Infof("Started %s \"%s\" for %s", r.Method, r.URL.Path, r.RemoteAddr)
	self.Log.Infof("\tProcessing by %s as %s", self.getHandlerName(self.HandlerFunc), r.Header.Get("Accept"))

	if e := self.HandlerFunc(w, r, self); e != nil {
		self.SetStatus(http.StatusInternalServerError)
		self.Log.Errorf("\tHandlerFunc failed: %s", e.Error())
	}

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
