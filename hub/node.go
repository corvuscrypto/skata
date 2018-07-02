package hub

import (
	"skata/common"
)

// SkataNode is the representation of a piece of the Skata network
// that the program can recognize without confusion
type SkataNode struct {
	pipe
	ID   common.SkataNodeID
	Type common.SkataNodeType
}
