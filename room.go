package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan interface{}
	join    chan *client
	leave   chan *client

	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			for connectedClient := range r.clients {
				connectedClient.send <- onlineCountChangedDto{len(r.clients)}
			}
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			for connectedClient := range r.clients {
				connectedClient.send <- onlineCountChangedDto{len(r.clients)}
			}
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Cookie error: ", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan interface{}, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan interface{}),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
