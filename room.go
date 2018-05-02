package main

import (
	"log"
	"net/http"

	. "./dtos"
	. "./repositories"

	"github.com/gorilla/websocket"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/stretchr/objx"
)

var dao = ChatRepository{}

type room struct {
	forward chan ChatMessageDto
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
				connectedClient.send <- OnlineCountChangedDto{len(r.clients)}
			}
			messages, err := dao.GetLast()
			if err != nil {
				return
			}
			for _, message := range messages {
				client.send <- message
			}
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			for connectedClient := range r.clients {
				connectedClient.send <- OnlineCountChangedDto{len(r.clients)}
			}
		case msg := <-r.forward:
			dao.Insert(msg)
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

	conf := config.NewConfig()
	conf.Load(file.NewSource(
		file.WithPath("config.json"),
	))
	conf.Get("hosts", "database").Scan(&dao)
	dao.Connect()

	return &room{
		forward: make(chan ChatMessageDto),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
