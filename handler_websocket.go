package absolut

import (
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
	callbacks WebsocketCallbacks
}

func NewHandlerWebsocket(callbacks WebsocketCallbacks) *HandlerWebsocket {
	self := &HandlerWebsocket{
		HandlerBase: HandlerBase{
			Log:    &logStack{},
			status: http.StatusOK,
			start:  time.Now(),
		},
		callbacks: callbacks,
	}
	self.Child = self
	return self
}

func (self *HandlerWebsocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		self.Log.Flush()
	}()

	self.serveHTTPPreffix(r)

	if e := self.callbacks.BeforeUpgrade(w, r, self); e != nil {
		self.SetStatus(http.StatusBadRequest)
		self.Log.Errorf("\tBeforeUpgrade failed: %s", e.Error())
	} else {
		ws, e := websocket.Upgrade(w, r, nil, 1024, 1024)
		if e != nil {
			self.SetStatus(http.StatusInternalServerError)
			self.Log.Errorf("\twebsocket.Upgrade failed: %s", e.Error())
		} else {
			self.SetStatus(http.StatusSwitchingProtocols)
			go self.callbacks.AfterConnect(ws)
			go func() {
				for {
					messageType, r, e := ws.NextReader()
					if e != nil {
						self.callbacks.OnClose(e)
						return
					}

					self.callbacks.OnMessage(messageType, r)
				}
			}()
		}
	}

	self.serveHTTPSuffix(w)
}
