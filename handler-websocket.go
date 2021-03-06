package absolut

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketServerInitializer interface {
	HandlerBeforeUpgrade(http.ResponseWriter, *http.Request, WebsocketHandlerClient) (WebsocketServerCallbacks, error)
}

type WebsocketServerCallbacks interface {
	AfterConnect(*websocket.Conn)
	OnMessage(int, io.Reader, *websocket.Conn)
	OnClose(error)
}

type HandlerWebsocket struct {
	HandlerBase
	readWait    time.Duration
	pingPeriod  time.Duration
	initializer WebsocketServerInitializer
	callbacks   WebsocketServerCallbacks
}

func NewHandlerWebsocket(initializer WebsocketServerInitializer, readWait time.Duration) *HandlerWebsocket {
	self := &HandlerWebsocket{
		HandlerBase: HandlerBase{
			log:    LogStackNew(),
			status: http.StatusOK,
			start:  time.Now(),

			isPoll:       false,
			pollStatuses: make(map[int]bool),
		},
		readWait:    readWait,
		pingPeriod:  ((readWait) * 9) / 10,
		initializer: initializer,
	}
	self.Child = self
	return self
}

func (self *HandlerWebsocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer self.Log().Flush()

	self.serveHTTPPreffix(r)

	if callbacks, e := self.initializer.HandlerBeforeUpgrade(w, r, self); e != nil {
		self.SetStatus(http.StatusBadRequest)
		self.Log().Errorf("HandlerBeforeUpgrade failed: %s", e.Error())
	} else {
		self.callbacks = callbacks
		upgrader := websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  GetWebsocketReadBufferSize(),
			WriteBufferSize: GetWebsocketWriteBufferSize(),
		}
		ws, e := upgrader.Upgrade(w, r, nil)
		if e != nil {
			self.SetStatus(http.StatusInternalServerError)
			self.Log().Errorf("websocket.Upgrade failed: %s", e.Error())
		} else {
			self.callbacks.AfterConnect(ws)

			go func() {
				ticker := time.NewTicker(self.pingPeriod)

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

				self.SetStatus(http.StatusSwitchingProtocols)
				for {
					messageType, r, e := ws.NextReader()
					if e != nil {
						self.callbacks.OnClose(e)
						return
					}

					self.callbacks.OnMessage(messageType, r, ws)
				}
			}()
		}
	}

	self.serveHTTPSuffix(w)
}
