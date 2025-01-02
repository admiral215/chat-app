package websocket

import "github.com/gorilla/websocket"

type Connection struct {
	Conn   *websocket.Conn
	UserId string
}
