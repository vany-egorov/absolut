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

type WebsocketClientInitializer interface {
	ClientBeforeConnect(WebsocketHandlerClient) (WebsocketClientCallbacks, *http.Header)
}

type WebsocketClientCallbacks interface {
	AfterConnect(*websocket.Conn)
	OnMessage(int, io.Reader)
	OnError(error)
	OnClose(error)
}

type ClientWebsocket struct {
	url         *url.URL
	initializer WebsocketClientInitializer
	callbacks   WebsocketClientCallbacks
	httpHeader  *http.Header
	readWait    time.Duration
	pingPeriod  time.Duration
	Log         *LogStack
}

func (self *ClientWebsocket) GetLog() *LogStack {
	return self.Log
}

func Φ(u *url.URL, initializer WebsocketClientInitializer, readWait time.Duration) {
	newClientWebsocket(u, initializer, readWait*time.Second)
}

func newClientWebsocket(u *url.URL, initializer WebsocketClientInitializer, readWait time.Duration) {
	ticker := time.NewTicker(2 * time.Second)

	defer func() {
		ticker.Stop()
	}()

	self := &ClientWebsocket{
		url:         u,
		initializer: initializer,
		Log:         new(LogStack),
		readWait:    readWait,
		pingPeriod:  ((readWait) * 9) / 10,
	}

	callbacks, httpHeader := self.initializer.ClientBeforeConnect(self)
	if httpHeader == nil {
		self.httpHeader = new(http.Header)
	} else {
		self.httpHeader = httpHeader
	}
	self.callbacks = callbacks

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
	ticker := time.NewTicker(self.pingPeriod)

	if e != nil {
		return fmt.Errorf("wsConnect failed: %s", e.Error())
	}

	defer func() {
		ws.Close()
		ticker.Stop()
	}()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if e := ws.WriteMessage(websocket.PingMessage, []byte{}); e != nil {
					return
				}
			}
		}
	}()

	ws.SetReadDeadline(time.Now().Add(self.readWait))
	ws.SetPongHandler(func(s string) error {
		ws.SetReadDeadline(time.Now().Add(self.readWait))
		return nil
	})

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

	ws, _, e = websocket.NewClient(c, self.url, *self.httpHeader, 1024, 1024)
	if e != nil {
		return nil, fmt.Errorf("websocket.NewClient failed: %s", e)
	}

	return ws, nil
}
