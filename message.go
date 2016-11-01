package rtbroadcaster

import (
	"encoding/json"
	"fmt"
)

type messageStatus struct {
	Value int    `json:"value"` // Connetion status: 0 = not connected, 1 = new, 2 = join, 3 = connected, 4 = close
	Text  string `json:"text"`  // Status message
}

type message struct {
	UUID         string        `json:"uuid"`         // Room uuid
	Status       messageStatus `json:"status"`       // Connection status
	FuncKey      string        `json:"funcKey"`      // Function to execute key
	FuncParams   []string      `json:"funcParams"`   // Function to execute parameters
	StateMessage bool          `json:"stateMessage"` // If it should be saved as state message
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
