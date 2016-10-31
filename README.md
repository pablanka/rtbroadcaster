# rtbroadcaster

It broadcasts messages from one websocket client to all connected clients. Golang based, using Gorilla websocket package.

## How it works

* **Manager:**
    
    Creates new client connections and rooms. It adds/removes clients to specific room.

* **room:**

    Maintains the set of active clients and broadcasts messages to the room's clients.

* **client:**

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

```
broadcastsMgr := rtbroadcaster.NewManager() // Creates broadcast manager
```

Create a new websocket client connection.

```
http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
	broadcastsMgr.CreateNewClient(w, r) // create a new socket client and manage it.
})
```

Full example:

```
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

### Message configuration

Messages are used for:

* **Create new broacast room:**

* **Join to existing room:**

* **Close a room (only room's owner):**

* **Send action (only room's owner):**


## How to use it (Cient side)

## Authors

* **Pablo Acuña**


