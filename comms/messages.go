package comms

import (
	"encoding/binary"
	"skata/common"
	"time"
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
	Serialize() []byte
	Deserialize([]byte) error
}

// SkataMessageBase is the base message structure
// from which all other message types should be composed
type SkataMessageBase struct {
	source common.SkataNodeID
}

// SignalType is the type alias for defining signals
type SignalType uint8

// Defined Signal types
const (
	Hello SignalType = iota
)

// SkataSignal is a signal that the receiver MUST treat as
// a command.
type SkataSignal struct {
	SkataMessageBase
	Signal SignalType
}

// Type satisfies the message interface
func (s SkataSignal) Type() SkataMessageType {
	return Signal
}

// Serialize Satisfies the message interface
func (s *SkataSignal) Serialize() (data []byte) {
	data = make([]byte, 9)
	binary.BigEndian.PutUint64(data, uint64(s.source))
	data[8] = byte(s.Signal)
	return
}

// Deserialize Satisfies the message interface
func (s *SkataSignal) Deserialize(data []byte) (err error) {
	s.source = common.SkataNodeID(binary.BigEndian.Uint64(data[:8]))
	s.Signal = SignalType(data[8])
	return
}

// SkataEvent is an arbitrary event.
// This allows arbitrary data to be communicated by
// the scheduler or nodes and if the nodes CAN handle it,
// they SHOULD
type SkataEvent struct {
	SkataMessageBase
	Timestamp time.Time
	EventName string
}

// Type satisfies the message interface
func (s SkataEvent) Type() SkataMessageType {
	return Event
}

// Serialize Satisfies the message interface
func (s *SkataEvent) Serialize() (data []byte) {
	data = make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(s.source))
	timeBytes, _ := s.Timestamp.MarshalBinary()
	timeLengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeLengthBytes, uint64(len(timeBytes)))
	timeBytes = append(timeLengthBytes, timeBytes...)
	data = append(data, timeBytes...)
	data = append(data, []byte(s.EventName)...)
	return
}

// Deserialize Satisfies the message interface
func (s *SkataEvent) Deserialize(data []byte) (err error) {
	s.source = common.SkataNodeID(binary.BigEndian.Uint64(data[:8]))
	timestampLength := binary.BigEndian.Uint64(data[8:16])
	timestamp := new(time.Time)
	if err = timestamp.UnmarshalBinary(data[16 : 16+timestampLength]); err != nil {
		return
	}
	s.Timestamp = *timestamp
	s.EventName = string(data[16+timestampLength:])
	return
}

// RequestType is a type alias for defining requests
type RequestType uint8

// Defined request types
const (
	Status RequestType = iota
)

// SkataRequest is a request to a node.
// A node SHOULD respond to a request with a SkataResponse
type SkataRequest struct {
	SkataMessageBase
	Request RequestType
	ID      string
}

// Serialize Satisfies the message interface
func (s *SkataRequest) Serialize() (data []byte) {
	data = make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(s.source))
	data = append(data, byte(s.Request))
	data = append(data, []byte(s.ID)...)
	return
}

// Deserialize Satisfies the message interface
func (s *SkataRequest) Deserialize(data []byte) (err error) {
	s.source = common.SkataNodeID(binary.BigEndian.Uint64(data[:8]))
	s.Request = RequestType(data[8])
	s.ID = string(data[9:])
	return
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

// Serialize Satisfies the message interface
func (s *SkataResponse) Serialize() (data []byte) {
	data = make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(s.source))
	dataLengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLengthBytes, uint64(len(s.Data)))
	data = append(data, dataLengthBytes...)
	data = append(data, s.Data...)
	data = append(data, []byte(s.RequestID)...)
	return
}

// Deserialize Satisfies the message interface
func (s *SkataResponse) Deserialize(data []byte) (err error) {
	s.source = common.SkataNodeID(binary.BigEndian.Uint64(data[:8]))
	dataLength := binary.BigEndian.Uint64(data[8:16])
	s.Data = data[16 : 16+dataLength]
	s.RequestID = string(data[16+dataLength:])
	return
}

// Type satisfies the message interface
func (s SkataResponse) Type() SkataMessageType {
	return Response
}

// SkataCustom is just a custom message wrapper
type SkataCustom struct {
	SkataMessageBase
	Data []byte
}

// Type satisfies the message interface
func (s SkataCustom) Type() SkataMessageType {
	return Custom
}

// Serialize Satisfies the message interface
func (s *SkataCustom) Serialize() (data []byte) {
	data = make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(s.source))
	dataLengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLengthBytes, uint64(len(s.Data)))
	data = append(data, dataLengthBytes...)
	data = append(data, s.Data...)
	return
}

// Deserialize Satisfies the message interface
func (s *SkataCustom) Deserialize(data []byte) (err error) {
	s.source = common.SkataNodeID(binary.BigEndian.Uint64(data[:8]))
	dataLength := binary.BigEndian.Uint64(data[8:16])
	s.Data = data[16 : 16+dataLength]
	return
}
