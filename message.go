package rtbroadcaster

import (
	"encoding/json"
	"fmt"
)

type messageStatus struct {
	// connetion status: 0 = not connected, 1 = connected, 2 = closed
	Value int

	// status message
	Text string
}

type message struct {

	// room uuid
	UUID string

	// connetion status
	Status messageStatus

	// function to execute key
	FuncKey string

	// function to execute parameters
	FuncParams []string
}

func decodeMessageFromJSON(jsonMessage []byte) *message {
	var msg message
	err := json.Unmarshal(jsonMessage, &msg)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &msg
}

func encodeJSONFromMessage(_message *message) []byte {
	bytemsg, err := json.Marshal(_message)
	if err != nil {
		fmt.Println("error:", err)
	}
	return bytemsg
}
