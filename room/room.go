package room

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"

	"github.com/sasakiyudai/websocket-chat/trace"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize}

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	Tracer  trace.Tracer
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.Tracer.Trace("new client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.Tracer.Trace("client left")
		case msg := <-r.forward:
			r.Tracer.Trace("client send message: ", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg:
					// send message
					r.Tracer.Trace("-- message sent successfully")
				default:
					// send failed
					delete(r.clients, client)
					close(client.send)
					r.Tracer.Trace("-- message sending failed")
				}
			}
		}
	}
}

func NewRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		Tracer:  trace.Off(),
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("getting cookie info", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
