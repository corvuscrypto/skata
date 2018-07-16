package hub

import (
	"skata/common"
	"skata/comms"
)

// SkataNode is the representation of a piece of the Skata network
// that the program can recognize without confusion
type SkataNode struct {
	Pipe *comms.Connection
	ID   common.SkataNodeID
	Type common.SkataNodeType
}

// NewSkataNode is the factory for SkataNodes
func NewSkataNode(conn *comms.Connection) *SkataNode {
	node := new(SkataNode)
	node.Pipe = conn
	node.ID = conn.Source
	node.Type = node.ID.GetNodeType()
	return node
}
