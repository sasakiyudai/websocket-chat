package room

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"

	"github.com/sasakiyudai/websocket-chat/trace"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)
var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize}

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	Tracer	trace.Tracer
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
			r.Tracer.Trace("client send message: ", string(msg))
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
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		Tracer: trace.Off(),
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
