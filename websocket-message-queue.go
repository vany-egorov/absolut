package absolut

import (
	"time"

	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
)

type websocketMessageQueueConnection interface {
	GetWs() *websocket.Conn
}

type websocketMessageQueuer interface {
	Dequeue()
	First() interface{}
	Len() int
}

type WebsocketMessageQueue struct {
	q []interface{}
}

func (self *WebsocketMessageQueue) Enqueue(response interface{}) {
	self.q = append(self.q, response)
}

func (self *WebsocketMessageQueue) Dequeue() {
	self.q = self.q[1:len(self.q)]
}

func (self *WebsocketMessageQueue) First() interface{} {
	return self.q[0]
}

func (self *WebsocketMessageQueue) Len() int {
	return len(self.q)
}

func Ψ(
	q websocketMessageQueuer,
	c websocketMessageQueueConnection,
	period time.Duration,
) {
	responsesDispatcher(q, c, period)
}

func responsesDispatcher(
	q websocketMessageQueuer,
	c websocketMessageQueueConnection,
	period time.Duration,
) {
	ticker := time.NewTicker(period * time.Second)

	defer func() {
		ticker.Stop()
		log.Flush()
	}()

	for {
		select {
		case <-ticker.C:

		whileQueueIsNotEmpty:
			for q.Len() > 0 {
				if q.Len() == 0 {
					break whileQueueIsNotEmpty
				}

				message := q.First()
				ws := c.GetWs()
				if ws == nil {
					break whileQueueIsNotEmpty
				}

				if e := JsonToWs(message, ws); e != nil {
					log.Errorf("[Queue.Len() => %4d] JsonToWs failed: %s \n", q.Len(), e.Error())
				}

				q.Dequeue()
			}
		}
	}
}
