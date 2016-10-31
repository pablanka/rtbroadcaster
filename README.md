# rtbroadcaster

It broadcasts messages from one websocket client to all connected clients. Golang based, using Gorilla websocket package.
There is a manager that mamages all clientes and rooms. A

## How it works

* Manager:
    
    Creates new client connections and rooms. It adds/removes clients to specific room.

* room



## How to use it

Get the package with go get command:

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

## Authors

* **Pablo Acuña**


