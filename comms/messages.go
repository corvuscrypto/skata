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

// SkataMessage is the message interface that all messages should have
type SkataMessage interface {
	Type() SkataMessageType
}

// SkataMessageBase is the base message structure
// from which all other message types should be composed
type SkataMessageBase struct {
	source common.SkataNodeID
}

// SkataSignal is a signal that the receiver MUST treat as
// a command.
type SkataSignal struct {
	SkataMessageBase
	Signal int
}

// Type satisfies the message interface
func (s SkataSignal) Type() SkataMessageType {
	return Signal
}

// SkataEvent is an arbitrary event.
// This allows arbitrary data to be communicated by
// the scheduler or nodes and if the nodes CAN handle it,
// they SHOULD
type SkataEvent struct {
	SkataMessageBase
	EventName string
	Timestamp int
}

// Type satisfies the message interface
func (s SkataEvent) Type() SkataMessageType {
	return Event
}

// SkataRequest is a request to a node.
// A node SHOULD respond to a request with a SkataResponse
type SkataRequest struct {
	SkataMessageBase
	RequestType string
	ID          string
}

// Type satisfies the message interface
func (s SkataRequest) Type() SkataMessageType {
	return Request
}

// SkataResponse is the response to a SkataRequest and
// completes a two-way communication between nodes.
type SkataResponse struct {
	SkataMessageBase
	RequestID string
	Data      []byte
}

// Type satisfies the message interface
func (s SkataResponse) Type() SkataMessageType {
	return Response
}
