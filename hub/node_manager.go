package hub

import "skata/comms"

// NodeManager manages the nodes that are connected to the hub
type NodeManager struct {
	Nodes    []*SkataNode
	Listener *comms.Listener
}

// NewNodeManager creates a NodeManager and returns it
func NewNodeManager(listenAddr string) *NodeManager {
	manager := new(NodeManager)
	manager.Listener = comms.NewListener(listenAddr)
	manager.Nodes = []*SkataNode{}
	go manager.waitForConnections()
	return manager
}

func (n *NodeManager) waitForConnections() {
	go n.Listener.ListenAndAccept()
	for {
		connection := <-n.Listener.ConnectionChan
		node := NewSkataNode(connection)
		n.Nodes = append(n.Nodes, node)
	}
}
