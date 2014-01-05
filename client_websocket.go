package absolut

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketClientCallbacks interface {
	ClientBeforeConnect(WebsocketHandlerClient)

	AfterConnect(*websocket.Conn)
	OnMessage(int, io.Reader)
	OnError(error)
	OnClose(error)
}

type ClientWebsocket struct {
	url       *url.URL
	callbacks WebsocketClientCallbacks
	Log       *LogStack
}

func (self *ClientWebsocket) GetLog() *LogStack {
	return self.Log
}

func Φ(u *url.URL, callbacks WebsocketClientCallbacks) {
	newClientWebsocket(u, callbacks)
}

func newClientWebsocket(u *url.URL, callbacks WebsocketClientCallbacks) {
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		ticker.Stop()
	}()

	self := &ClientWebsocket{
		url:       u,
		callbacks: callbacks,
		Log:       new(LogStack),
	}

	self.callbacks.ClientBeforeConnect(self)

	for {
		select {
		case <-ticker.C:
			if e := self.handle(); e != nil {
				self.callbacks.OnError(e)
				continue
			}
		}
	}
}

func (self *ClientWebsocket) handle() error {
	ws, e := self.connect()

	if e != nil {
		return fmt.Errorf("wsConnect failed: %s", e.Error())
	}

	defer ws.Close()

	self.callbacks.AfterConnect(ws)

	for {
		messageType, r, e := ws.NextReader()
		if e != nil {
			self.callbacks.OnClose(e)
			return fmt.Errorf("ws.NextReader failed: %s", e.Error())
		}

		self.callbacks.OnMessage(messageType, r)
	}
}

func (self *ClientWebsocket) connect() (ws *websocket.Conn, e error) {
	c, e := net.Dial("tcp", self.url.Host)
	if e != nil {
		return nil, fmt.Errorf("net.Dial failed: %s", e.Error())
	}

	ws, _, e = websocket.NewClient(c, self.url, http.Header{}, 1024, 1024)
	if e != nil {
		return nil, fmt.Errorf("websocket.NewClient failed: %s", e)
	}

	return ws, nil
}
