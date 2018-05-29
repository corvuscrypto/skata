package common

import "github.com/streadway/amqp"

type SkataConnectionConfig struct {
	RabbitMQAddress string
	DryTest         bool
}

// SkataConnection is a wrapper for some convenience
// that handles instance identification and connection
// upkeep
type SkataConnection struct {
	InstanceID   SkataNodeID
	MQConnection *amqp.Connection
}

// DefaultConnectionConfig is the default connection setting
var DefaultConnectionConfig = &SkataConnectionConfig{
	RabbitMQAddress: "",
	DryTest:         true,
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
		conn.MQConnection, err = amqp.Dial(config.RabbitMQAddress)
	}
	return
}
