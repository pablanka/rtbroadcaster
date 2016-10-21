package rtbroadcaster

import "github.com/satori/go.uuid"

// Room maintains the set of active clients and broadcasts messages to the clients.
type Room struct {
	// Room uuid
	uuid uuid.UUID

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Messages that will be sent to new clients.
	stateMesages []message
}

/* PRIVATE FUNCS */

func newRoom() *Room {
	return &Room{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register: // Register a new client to this presentation
			r.clients[client] = true
		case client := <-r.unregister: // Unegister a new client to this presentation
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast: // Broadcast a message to all clients in this presentation
			r.broadcastMessage(message)
		}
	}
}

func (r *Room) broadcastMessage(message []byte) {
	for client := range r.clients {
		if !client.isOwner {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(r.clients, client)
			}
		}
	}
}
