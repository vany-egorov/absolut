package absolut

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketCallbacks interface {
	BeforeUpgrade(http.ResponseWriter, *http.Request, *HandlerWebsocket) error
	AfterConnect(*websocket.Conn)
	OnMessage(int, io.Reader)
	OnClose(error)
}

type HandlerWebsocket struct {
	HandlerBase
	Callbacks WebsocketCallbacks
}

func NewHandlerWebsocket(callbacks WebsocketCallbacks) *HandlerWebsocket {
	self := &HandlerWebsocket{
		HandlerBase: HandlerBase{
			status: http.StatusOK,
			Log:    &logStack{},
		},
		Callbacks: callbacks,
	}
	self.Child = self
	return self
}

func (self *HandlerWebsocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		self.Log.Flush()
	}()

	start := time.Now()
	self.Log.Infof("Started %s \"%s\" for %s", r.Method, r.URL.Path, r.RemoteAddr)

	if e := self.Callbacks.BeforeUpgrade(w, r, self); e != nil {
		self.SetStatus(http.StatusBadRequest)
		self.Log.Errorf("\tBeforeUpgrade failed: %s", e.Error())
	} else {
		ws, e := websocket.Upgrade(w, r, nil, 1024, 1024)
		if e != nil {
			self.SetStatus(http.StatusInternalServerError)
			self.Log.Errorf("\twebsocket.Upgrade failed: %s", e.Error())
		} else {
			self.SetStatus(http.StatusSwitchingProtocols)
			go self.Callbacks.AfterConnect(ws)
			go func() {
				for {
					messageType, r, e := ws.NextReader()
					if e != nil {
						self.Callbacks.OnClose(e)
						return
					}

					self.Callbacks.OnMessage(messageType, r)
				}
			}()
		}
	}

	if self.GetStatus() != http.StatusOK && self.GetStatus() != http.StatusSwitchingProtocols {
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
