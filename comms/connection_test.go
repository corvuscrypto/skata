package comms

import (
	"io"
	"net"
	"skata/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPacketWrapping(t *testing.T) {
	signal := new(SkataSignal)
	signal.Signal = Hello
	signal.source = common.GenerateID(common.HubNode)

	packet := createPacket(signal)

	message := parsePacket(packet)
	assert.Equal(t, message, signal)
}

type testRoundTripper struct {
	*net.TCPConn
}

func (t *testRoundTripper) run() {
	addr, _ := net.ResolveTCPAddr("tcp", ":9000")
	listener, _ := net.ListenTCP("tcp", addr)
	conn, err := listener.AcceptTCP()
	if err != nil {
		panic(err)
	}
	t.TCPConn = conn
	for {
		data := make([]byte, 1024)
		n, err := t.Read(data)
		if err == io.EOF {
			return
		}
		t.Write(data[:n])
	}
}

func setupRoundTripper() *testRoundTripper {
	rt := new(testRoundTripper)
	go rt.run()
	return rt
}

func TestRoundTripConnection(t *testing.T) {
	signal := new(SkataSignal)
	signal.Signal = Hello
	signal.source = common.GenerateID(common.HubNode)
	rt := setupRoundTripper()
	time.Sleep(1)
	connection := NewConnection("localhost", "9000")

	connection.Write(signal)

	msg := <-connection.Pipe

	assert.Equal(t, msg, signal)

	connection.Close()
	rt.Close()
}
