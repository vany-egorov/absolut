package absolut

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HandlerHTTP struct {
	HandlerBase
	HandlerFunc HandlerHTTPFuncType
	Extension   Extension
}

type HandlerHTTPFuncType func(http.ResponseWriter, *http.Request, *HandlerHTTP) error

func NewHandlerHTTP(handler HandlerHTTPFuncType) *HandlerHTTP {
	self := &HandlerHTTP{
		HandlerBase: HandlerBase{
			status: http.StatusOK,
			Log:    LogStackNew(),
			start:  time.Now(),
		},
		HandlerFunc: handler,
		Extension:   JSON,
	}
	self.Child = self

	return self
}

func (self *HandlerHTTP) GetExtension() Extension { return self.GetExt() }
func (self *HandlerHTTP) GetExt() Extension       { return self.Extension }

func (self *HandlerHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer self.Log.Flush()

	if extension, ok := mux.Vars(r)["extension"]; ok {
		self.Extension = GuessExtension(extension)
	}

	w.Header().Set("Content-Type", self.GetExt().GetContentType())

	self.serveHTTPPreffix(r)

	self.Log.Debugf("\tProcessing by %s as %s", self.GetHandlerName(self.HandlerFunc), ExtensionText(self.Extension))

	if e := self.HandlerFunc(w, r, self); e != nil {
		self.SetStatus(http.StatusInternalServerError)
		self.Log.Errorf("\tHandlerFunc failed: %s", e.Error())
	}

	self.serveHTTPSuffix(w)
}
