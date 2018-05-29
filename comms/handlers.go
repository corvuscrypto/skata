package comms

// EventHandler is a handler specific to a SkataEvent
type EventHandler func(SkataEvent) error

// SignalHandler is a handler specific to a SkataSignal
type SignalHandler func(SkataSignal) error

// RequestHandler is a handler specific to a SkataRequest
type RequestHandler func(SkataRequest) error

// MessageHandler is a collection of handlers for each specific message type
type MessageHandler struct {
	eventHandlers   map[string]EventHandler
	SignalHandlers  map[string]SignalHandler
	RequestHandlers map[string]RequestHandler
}
