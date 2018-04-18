package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	// Typ client reprezentuje pojedynczego użytkownika
	// prowadzącego konwersację z użyciem komunikatora.
	// socket to gniazdo internetowe do obsługi danego klienta.
	socket *websocket.Conn
	// send to kanał, którym są przesyłane komunikaty.
	send        chan []byte
	onlineCount chan []byte
	// room to pokój rozmów używany przez klienta.
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
	for oCount := range c.onlineCount {
		err := c.socket.WriteMessage(websocket.PingMessage, oCount)
		if err != nil {
			return
		}
	}
}
