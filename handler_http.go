package absolut

import (
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
			start:  time.Now(),
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

	self.serveHTTPPreffix(r)

	self.Log.Infof("\tProcessing by %s as %s", self.getHandlerName(self.HandlerFunc), r.Header.Get("Accept"))

	if e := self.HandlerFunc(w, r, self); e != nil {
		self.SetStatus(http.StatusInternalServerError)
		self.Log.Errorf("\tHandlerFunc failed: %s", e.Error())
	}

	self.serveHTTPSuffix(w)
}
