package absolut

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HandlerHTTP struct {
	HandlerBase
	HandlerFunc HandlerHTTPFuncType
}

type HandlerHTTPFuncType func(http.ResponseWriter, *http.Request, *HandlerHTTP) error

func NewHandlerHTTP(handler HandlerHTTPFuncType, loggerGetter LoggerGetter) *HandlerHTTP {
	var logStack *LogStack
	if loggerGetter == nil {
		logStack = LogStackNew()
	} else {
		logStack = LogStackNew(loggerGetter)
	}

	self := &HandlerHTTP{
		HandlerBase: HandlerBase{
			status: http.StatusOK,
			log:    logStack,
			start:  time.Now(),

			isPoll: false,
			pollStatuses: map[int]bool{
				http.StatusNotFound:    true,
				http.StatusNotModified: true,
			},
		},
		HandlerFunc: handler,
	}
	self.extension = JSON
	self.Child = self

	return self
}

func (self *HandlerHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer self.Log().Flush()

	if extension, ok := mux.Vars(r)["extension"]; ok {
		self.extension = GuessExtension(extension)
	} else {
		self.extension = GuessExtension("json")
	}

	self.serveHTTPPreffix(r)

	if e := self.HandlerFunc(w, r, self); e != nil {
		self.SetStatus(http.StatusInternalServerError)
		self.Log().Errorf("HandlerFunc failed: %s", e.Error())
	}

	w.Header().Set("Content-Type", self.GetExt().GetContentType())

	self.serveHTTPSuffix(w)
}
