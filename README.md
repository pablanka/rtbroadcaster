# rtbroadcaster

It broadcasts messages from one websocket client to all connected clients. Golang based, using Gorilla websocket package.
There is a manager object that mamages all clientes and rooms. A new client connection request is 

## How it works

* **Manager:**
    
    Creates new client connections and rooms. It adds/removes clients to specific room.

* **room:**

    Maintains the set of active clients and broadcasts messages to the room's clients.

* **client:**

    Is an middleman between the websocket connection and its room. 
    There is only one room's owner. Only the room's owner can broadcast messages
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

## How to use it (Cient side)

## Authors

* **Pablo Acuña**


