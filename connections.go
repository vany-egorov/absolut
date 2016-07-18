package absolut

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Single struct{ Ws *websocket.Conn }

func (self *Single) OnConnected(ws *websocket.Conn) { self.Ws = ws }
func (self *Single) OnDisconnected()                { self.Ws = nil }
func (self *Single) GetWs() *websocket.Conn         { return self.Ws }
func (self *Single) IsConnected() bool              { return self.Ws != nil }

func NewSingle() *Single { return new(Single) }

type Multiple struct {
	sync.RWMutex

	m map[string]*websocket.Conn
}

func (self *Multiple) Len() int      { return len(self.m) }
func (self *Multiple) IsEmpty() bool { return len(self.m) == 0 }

func (self *Multiple) Broadcast(m interface{}) error {
	b, e := json.Marshal(m)

	if e != nil {
		return e
	}

	return self.BroadcastByte(b)
}

func (self *Multiple) BroadcastByte(b []byte) (e error) {
	self.Lock()
	defer self.Unlock()

	for _, ws := range self.m {
		e = ws.WriteMessage(websocket.TextMessage, b)
	}

	return
}

func (self *Multiple) OnConnected(id string, ws *websocket.Conn) {
	self.Lock()
	defer self.Unlock()

	self.m[id] = ws
}
func (self *Multiple) OnMessage(id string, ws *websocket.Conn) {
	if _, ok := self.m[id]; ok {
		return
	}
	self.OnConnected(id, ws)
}

func (self *Multiple) OnDisconnected(id string) {
	self.Lock()
	defer self.Unlock()

	delete(self.m, id)
}

func NewMultiple() *Multiple { return &Multiple{m: make(map[string]*websocket.Conn)} }
