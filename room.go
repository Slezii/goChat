package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward to kanał przechowujący nadsyłane komunikaty, // które należy przesłać do przeglądarki użytkownika. forward chan []byte
	forward chan []byte
	// join to kanał dla klientów, którzy chcą dołączyć do pokoju.
	join chan *client
	// leave to kanał dla klientów, którzy chcą opuścić pokój.
	leave chan *client

	clients map[*client]bool
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// dołączanie do pokoju
			r.clients[client] = true
			for connectedClient := range r.clients {
				connectedClient.onlineCount <- []byte("ddd")
			}
		case client := <-r.leave:
			// opuszczanie pokoju
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// rozsyłanie wiadomości do wszystkich klientów
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
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
