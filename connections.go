package absolut

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Single struct {
	Ws *websocket.Conn
}

func (self *Single) OnConnected(ws *websocket.Conn) {
	self.Ws = ws
}

func (self *Single) OnDisconnected() {
	self.Ws = nil
}

func (self *Single) GetWs() *websocket.Conn {
	return self.Ws
}

func (self *Single) IsConnected() bool {
	return self.Ws != nil
}

type Multiple struct {
	m map[string]*websocket.Conn
}

func newMultiple() *Multiple {
	return &Multiple{
		m: make(map[string]*websocket.Conn),
	}
}

func (self *Multiple) Broadcast(m interface{}) error {
	b, e := json.Marshal(m)

	if e != nil {
		return e
	}

	for _, ws := range self.m {
		if e := ws.WriteMessage(websocket.TextMessage, b); e != nil {
			return e
		}
	}

	return nil
}

func (self *Multiple) OnConnected(id string, ws *websocket.Conn) {
	self.m[id] = ws
}

func (self *Multiple) OnMessage(id string, ws *websocket.Conn) {
	if _, ok := self.m[id]; ok {
		return
	}
	self.OnConnected(id, ws)
}

func (self *Multiple) OnDisconnected(id string) {
	delete(self.m, id)
}
