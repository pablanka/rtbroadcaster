package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Pablanka/rtbroadcaster"
)

var addr = flag.String("addr", ":8080", "http service address")

var playerURL = "/Users/vixonic/PROJECTS/demos/360remote_control"

func main() {
	flag.Parse()

	broadcastsMgr := rtbroadcaster.NewManager() // Creates broadcast manager

	// Handle requests
	http.Handle("/", http.FileServer(http.Dir(playerURL)))
	http.HandleFunc("/presentation", func(w http.ResponseWriter, r *http.Request) {
		broadcastsMgr.CreateNewClient(w, r)
	})
	log.Println("Server running")
	// Serve and listen
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
