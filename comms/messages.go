package comms

import (
	"skata/common"
)

// SkataMessageType is the overriden uint type for messages
type SkataMessageType uint

// Skata message types
const (
	Event SkataMessageType = iota
	Signal
	Request
	Response
	Custom
)

// SkataMessage is the base message structure
// from which all other message types should be composed
type SkataMessage struct {
	messageType SkataMessageType
	source      common.SkataNodeID
}

// SkataSignal is a signal that the receiver MUST treat as
// a command.
type SkataSignal struct {
	SkataMessage
	Signal int
}

// SkataEvent is an arbitrary event.
// This allows arbitrary data to be communicated by
// the scheduler or nodes and if the nodes CAN handle it,
// they SHOULD
type SkataEvent struct {
	SkataMessage
	EventName string
	Timestamp int
}

// SkataRequest is a request to a node.
// A node SHOULD respond to a request with a SkataResponse
type SkataRequest struct {
	SkataMessage
	RequestType string
	ID          string
}

// SkataResponse is the response to a SkataRequest and
// completes a two-way communication between nodes.
type SkataResponse struct {
	SkataMessage
	RequestID string
	Data      []byte
}
