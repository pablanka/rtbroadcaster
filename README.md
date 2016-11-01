# rtbroadcaster

It broadcasts messages from one websocket client to all connected clients. Golang based, using Gorilla websocket package.

## How does it works

* **Manager:**
    
    Creates new client connections and rooms. It adds/removes clients to specific room.

* **Room:**

    Maintains the set of active clients and broadcasts messages to the room's clients.

* **Client:**

    Is a middleman between the websocket connection and its room. 
    There is only one room's owner. Only the room's owner can broadcast messages.
    - Pumps messages from the websocket connection to the room. 
    - Pumps messages from the room to the websocket connection.

## How to use it (Server side)

Get the package with "go get" command:

```
go get github.com/pablanka/rtbroadcaster
```

Create a new broadcast manager.

```go
broadcastsMgr := rtbroadcaster.NewManager() // Creates broadcast manager
```

Create a new websocket client connection.

```go
http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
	broadcastsMgr.CreateNewClient(w, r) // create a new socket client and manage it.
})
```

Full example:

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pablanka/rtbroadcaster"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	broadcastsMgr := rtbroadcaster.NewManager() // Creates broadcast manager

	// Handle requests
	http.HandleFunc("/broadcasting", func(w http.ResponseWriter, r *http.Request) {
		broadcastsMgr.CreateNewClient(w, r) // create a new socket client and manage it.
	})
	log.Println("Server running")

	// Serve and listen
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
```

### Message structure

```go
type messageStatus struct {
	// Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
	Value int

	// Status message
	Text string
}
```

```go
type message struct {

	// Room uuid
	UUID string

	// Connection status
	Status messageStatus

	// Function to execute key
	FuncKey string

	// Function to execute parameters
	FuncParams []string

	// If it should be saved as state message
	StateMessage bool
}
```

### Message configuration

Messages are used for:

#### Create new broacast room:

To create a new broadcast room, client must to send a message with **status.value = 1** and the other params must to be empty:

```javascript
{
	"uuid": "",
	"status": {
		"value": 1, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
		"text": "new connection"
	},
	"funcKey": "",
	"funcParams": [],
	"stateMessage": false
}
```

Server will response with:

```javascript
{
	"uuid": "63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 3, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
		"text": "connected"
	},
	"funcKey": "",
	"funcParams": [],
	"stateMessage": false
}
```

Then, client is able to continue broadcasting messages.

#### Join to existing room:

To join to an existing room, client must to send a message with the room's **uuid** and **status.value = 2**. The other params must to be empty:

```javascript
{
	"uuid": "f63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 2, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
		"text": "join connection"
	},
	"funcKey": "",
	"funcParams": [],
	"stateMessage": false
}
```

Server will response with:

```javascript
{
	"uuid": "63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 3, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
		"text": "connected"
	},
	"funcKey": "",
	"funcParams": [],
	"stateMessage": false
}
```

Then, client is able to receive messages.

#### Close a room (only room's owner):

To close to an existing room, client (only the room's owner) must to send a message with the room's **uuid** and **status.value = 4**. The other params must to be empty:

```javascript
{
	"uuid": "63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 4, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
		"text": "close connection"
	},
	"funcKey": "",
	"funcParams": [],
	"stateMessage": false
}
```

Server will response with:

```javascript
{
	"uuid": "63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 0, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = closed
		"text": "not connected"
	},
	"funcKey": "",
	"funcParams": [],
	"stateMessage": false
}
```

Then, all room's websockets will be closed.

#### Broadcast action (only room's owner):

To broadcast an action, client (only the room's owner) must to send a message with the room's **uuid**, **status.value = 4**, a **funcKey** and an array of string params **funcParams**. 
The other params must to be empty:

```javascript
{
	"uuid": "63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 3, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = closed
		"text": "connected"
	},
	"funcKey": "myKey",
	"funcParams": ["param1", "param2", "param3"],
	"stateMessage": false
}
```

Room will broacast the message to all connected clients and they could use **funcKey** and **funcParams** to execute actions.

#### Broacast state message (only room's owner):

An state message is one that is stored by the room. When a new client connection occurs in that room, all stored state messages will be sent to that client. 
It allow the new client to execute all these actions once connected.
To broadcast an state message, client (only the room's owner) must to send a message with the room's **uuid**, **status.value = 3**, 
a **funcKey**, an array of string params **funcParams** and **stateMessage = true**:

```javascript
{
	"uuid": "63ca67e-69bb-4a16-a71f-86a87acbe0b5",
	"status": {
		"value": 3, // Connection status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = closed
		"text": "connected"
	},
	"funcKey": "myKey",
	"funcParams": ["param1", "param2", "param3"],
	"stateMessage": true
}
```

Room will broadcast the message to all connected clients and they could use **funcKey** and **funcParams** to execute actions.

## How to use it (Cient side)

[**View Javascript SDK documentation**]()

## Authors

* **Pablo Acuña**


