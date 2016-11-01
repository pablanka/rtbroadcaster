package rtbroadcaster

import (
	"encoding/json"
	"fmt"
)

type messageStatus struct {
	// Connetion status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
	Value int

	// Status message
	Text string
}

type message struct {

	// Room uuid
	UUID string

	// Connetion status
	Status messageStatus

	// Function to execute key
	FuncKey string

	// Function to execute parameters
	FuncParams []string

	// If it should be saved as state message
	StateMessage bool

	// State message ID. To make a state message overwrite
	StateMessageID string
}

func decodeMessageFromJSON(jsonMessage []byte) *message {
	var msg message
	err := json.Unmarshal(jsonMessage, &msg)
	if err != nil {
		fmt.Println("decodeMessageFromJSON error:", err)
		fmt.Println("decodeMessageFromJSON error json:", string(jsonMessage))
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
