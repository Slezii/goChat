package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan interface{}
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var m chatMessageDto
		err := c.socket.ReadJSON(&m)
		if err != nil {
			log.Print(err)
			return
		}
		m.Author = c.userData["name"].(string)
		c.room.forward <- m
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
