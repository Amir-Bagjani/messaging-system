package entity

type WebSocketConn interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}
