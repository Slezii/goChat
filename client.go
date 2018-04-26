package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	// Typ client reprezentuje pojedynczego użytkownika
	// prowadzącego konwersację z użyciem komunikatora.
	// socket to gniazdo internetowe do obsługi danego klienta.
	socket *websocket.Conn
	send   chan interface{}
	room   *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		m := chatMessageDto{}
		err := c.socket.ReadJSON(&m)
		if err != nil {
			return
		}
		c.room.forward <- "dd"
	}
}
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
