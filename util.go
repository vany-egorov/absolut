package absolut

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
)

func JsonToWs(self interface{}, ws *websocket.Conn) error {
	b, e := json.Marshal(self)
	if e != nil {
		return fmt.Errorf("json.Marshal failed: %s", e.Error())
	}
	e = ws.WriteMessage(websocket.TextMessage, b)
	if e != nil {
		return fmt.Errorf("ws.WriteMessage failed: %s", e.Error())
	}
	return nil
}

func WsToJson(self interface{}, r io.Reader) error {
	buffer := new(bytes.Buffer)
	_, e := buffer.ReadFrom(r)
	if e != nil {
		return fmt.Errorf("buffer.ReadFrom failed: %s", e.Error())
	}

	e = json.Unmarshal(buffer.Bytes(), self)
	if e != nil {
		return fmt.Errorf("json.Unmarshal failed: %s", e.Error())
	}

	return nil
}

func JsonToPrettyString(self interface{}) string {
	b, e := json.MarshalIndent(self, " ", " ")
	if e != nil {
		return e.Error()
	}
	return string(b)
}
