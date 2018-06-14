package comms

// EventHandler is a handler specific to a SkataEvent
type EventHandler func(SkataEvent) error

// SignalHandler is a handler specific to a SkataSignal
type SignalHandler func(SkataSignal) error

// RequestHandler is a handler specific to a SkataRequest
type RequestHandler func(SkataRequest) error

// CustomHandler is a handler specific to a SkataCustom messages
// This allows custom messages to be passed.
type CustomHandler func(SkataCustom) error

// MessageHandler is a collection of handlers for each specific message type
type MessageHandler struct {
	EventHandlers   map[string]EventHandler
	SignalHandlers  map[SignalType]SignalHandler
	RequestHandlers map[RequestType]RequestHandler
	CustomHandlers  []CustomHandler
}

func (m *MessageHandler) handleMessage(msg SkataMessage) error {
	switch typedMsg := msg.(type) {
	case SkataEvent:
		handler, found := m.EventHandlers[typedMsg.EventName]
		if found {
			if err := handler(typedMsg); err != nil {
				return err
			}
		}
	case SkataSignal:
		handler, found := m.SignalHandlers[typedMsg.Signal]
		if found {
			if err := handler(typedMsg); err != nil {
				return err
			}
		}
	case SkataRequest:
		handler, found := m.RequestHandlers[typedMsg.Request]
		if found {
			if err := handler(typedMsg); err != nil {
				return err
			}
		}
	}
	return nil
}
