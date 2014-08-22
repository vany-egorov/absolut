package absolut

var (
	_websocketReadBufferSize  int = 1024
	_websocketWriteBufferSize int = 1024
)

func SetWebsocketReadBufferSize(v int) {
	_websocketReadBufferSize = v
}

func GetWebsocketReadBufferSize() int {
	return _websocketReadBufferSize
}

func SetWebsocketWriteBufferSize(v int) {
	_websocketWriteBufferSize = v
}

func GetWebsocketWriteBufferSize() int {
	return _websocketWriteBufferSize
}
