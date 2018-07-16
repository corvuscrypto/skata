package common

type SkataConnectionConfig struct {
	HubAddress string
	DryTest    bool
}

// SkataConnection is a wrapper for some convenience
// that handles instance identification and connection
// upkeep
type SkataConnection struct {
	InstanceID SkataNodeID
}

// DefaultConnectionConfig is the default connection setting
var DefaultConnectionConfig = &SkataConnectionConfig{
	HubAddress: "",
	DryTest:    true,
}

// NewSkataConnection creates a new connection and assigns
// an ID based on the node type
func NewSkataConnection(nodeType SkataNodeType, config *SkataConnectionConfig) (conn *SkataConnection, err error) {
	if config == nil {
		config = DefaultConnectionConfig
	}
	conn = new(SkataConnection)
	conn.InstanceID = GenerateID(nodeType)
	if !config.DryTest {
	}
	return
}
