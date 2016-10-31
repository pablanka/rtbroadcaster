# rtbroadcaster

It broadcasts messages from one websocket client to all connected clients. Golang based, using Gorilla websocket package.


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

End with an example of getting some data out of the system or using it for a little demo



## Authors

* **Pablo Acuña**


