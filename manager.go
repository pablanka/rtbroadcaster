package rtbroadcaster

import (
	"net/http"

	"github.com/satori/go.uuid"
)

// Manager manages rooms and clients
type Manager struct {
	rooms map[uuid.UUID]*Room
}

/* PUBLIC API */

// CreateNewClient creates a new client
func (mgr *Manager) CreateNewClient(w http.ResponseWriter, r *http.Request) {
	newClient(mgr, w, r)
}

// NewManager creates new broadcast rooms manager.
func NewManager() *Manager {
	return &Manager{
		rooms: make(map[uuid.UUID]*Room),
	}
}

/* PRIVATE FUNCS */

func (mgr *Manager) createNewRoom(client *Client) {
	room := newRoom()
	go room.run()
	room.uuid = uuid.NewV4()
	room.register <- client
	mgr.rooms[room.uuid] = room
	client.room = room
	client.isOwner = true
}

func (mgr *Manager) addToRoom(client *Client, _uuid string) {
	parsedUUID, _ := uuid.FromString(_uuid)
	room := mgr.rooms[parsedUUID]
	if room != nil {
		room.register <- client
		client.room = room
		client.isOwner = false
	}
}
