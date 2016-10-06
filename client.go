package rtbroadcaster

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is an middleman between the websocket connection and the hub.
type Client struct {
	// rooms manager
	manager *Manager

	// room whick this client belows
	room *Room

	// is this client is the broadcast owner
	isOwner bool

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the room.
func (c *Client) readPump() {
	defer func() {
		if c.room != nil {
			c.room.unregister <- c
		}
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		_message = bytes.TrimSpace(bytes.Replace(_message, newline, space, -1))

		if msg := decodeMessageFromJSON(_message); msg.Status.Value == 0 {
			if msg.UUID == "" {
				c.manager.createNewRoom(c)

				byteUUID, err := c.room.uuid.MarshalText()
				if err != nil {
					log.Println(err)
					return
				}
				stringUUID := string(byteUUID)

				connectionMsg := &message{
					UUID: stringUUID,
					Status: messageStatus{
						Value: 1,
						Text:  "Connected",
					},
					FuncKey:    "",
					FuncParams: nil,
				}
				c.send <- encodeJSONFromMessage(connectionMsg)
			} else {
				c.manager.addToRoom(c, msg.UUID)

				connectionMsg := &message{
					UUID: msg.UUID,
					Status: messageStatus{
						Value: 1,
						Text:  "Connected",
					},
					FuncKey:    "",
					FuncParams: nil,
				}
				c.send <- encodeJSONFromMessage(connectionMsg)
			}
		} else {
			if c.room != nil && c.isOwner {
				c.room.broadcast <- _message
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Client) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

// writePump pumps messages from the room to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case _message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(_message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(mgr *Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{manager: mgr, conn: conn, send: make(chan []byte, 256)}
	go client.writePump()
	client.readPump()
}