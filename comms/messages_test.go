package comms

import (
	"skata/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSkataSignal(t *testing.T) {
	signal := new(SkataSignal)
	signal.Signal = Hello
	signal.source = common.GenerateID(common.HubNode)

	signalBytes := signal.Serialize()

	newSignal := new(SkataSignal)
	err := newSignal.Deserialize(signalBytes)
	assert.NoError(t, err)
	assert.Equal(t, newSignal, signal)
}

func TestSkataEvent(t *testing.T) {
	event := new(SkataEvent)
	event.EventName = "test"
	event.Timestamp = time.Now()
	event.source = common.GenerateID(common.HubNode)

	eventBytes := event.Serialize()

	newEvent := new(SkataEvent)
	err := newEvent.Deserialize(eventBytes)
	assert.NoError(t, err)
	assert.Equal(t, newEvent, event)
}

func TestSkataRequest(t *testing.T) {
	request := new(SkataRequest)
	request.source = common.GenerateID(common.HubNode)
	request.Request = Status
	request.ID = "1234"

	requestBytes := request.Serialize()

	newRequest := new(SkataRequest)
	err := newRequest.Deserialize(requestBytes)
	assert.NoError(t, err)
	assert.Equal(t, newRequest, request)
}

func TestSkataResponse(t *testing.T) {
	response := new(SkataResponse)
	response.source = common.GenerateID(common.HubNode)
	response.RequestID = "1234"
	response.Data = []byte("test")

	responseBytes := response.Serialize()

	newResponse := new(SkataResponse)
	err := newResponse.Deserialize(responseBytes)
	assert.NoError(t, err)
	assert.Equal(t, newResponse, response)
}

func TestSkataCustom(t *testing.T) {
	message := new(SkataCustom)
	message.source = common.GenerateID(common.HubNode)
	message.Data = []byte("test")

	messageBytes := message.Serialize()

	newMessage := new(SkataCustom)
	err := newMessage.Deserialize(messageBytes)
	assert.NoError(t, err)
	assert.Equal(t, newMessage, message)
}
