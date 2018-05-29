package common

import "time"

//go:generate stringer -type=SkataNodeType

// TimeGenerator is an interface that exposes a method
// to get the current Time
type TimeGenerator interface {
	Now() time.Time
}

// DefaultClock is just a wrapper
type DefaultClock struct{}

// Now method to satisfy the TimeGenerator interface
func (d DefaultClock) Now() time.Time {
	return time.Now()
}

// IDClock is an Overridable clock from which times are grabbed
var IDClock TimeGenerator = DefaultClock{}

// SkataNodeType is an enumeration type for identifying
// the type of skata node.
type SkataNodeType uint

func (s SkataNodeType) String() string {
	return ""
}

// Types of skata nodes
const (
	HeartNode SkataNodeType = iota
	SchedulerNode
	WorkerNode
	HubNode
)

// SkataNodeID is the node ID type
type SkataNodeID uint64

// GetNodeVersion parses the ID and returns the node version
func (id SkataNodeID) GetNodeVersion() int {
	return int(id >> 40)
}

// GetNodeType parses the ID and returns the node type
func (id SkataNodeID) GetNodeType() SkataNodeType {
	return SkataNodeType(uint64(id) & 0xFF)
}

// GetNodeCreationTime parses the ID and returns the generation
// time of the ID
func (id SkataNodeID) GetNodeCreationTime() time.Time {
	return time.Unix(int64((uint64(id)>>8)&0xFFFFFFFF), 0).UTC()
}

// GenerateID takes a node type and generates a new ID
func GenerateID(nodeType SkataNodeType) SkataNodeID {
	generationTime := IDClock.Now().Unix()
	var id uint64 = version & 0xFF
	id <<= 32
	id |= uint64(generationTime)
	id <<= 8
	id |= uint64(nodeType)
	return SkataNodeID(id)
}
