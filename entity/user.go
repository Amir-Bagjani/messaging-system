package entity

type User struct {
	ID   string
	Conn WebSocketConn // WebSocket connection interface
}
