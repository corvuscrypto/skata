package comms

import (
	"encoding/binary"
	"net"
)

// MaxReadBytes is the max bytes to read from the stream
const MaxReadBytes = 1 << 10

// Acts as a frame
func createPacket(message SkataMessage) []byte {
	typeByte := byte(message.Type())
	data := message.Serialize()
	return append([]byte{typeByte}, data...)
}

func parsePacket(data []byte) SkataMessage {
	switch SkataMessageType(data[0]) {
	case Signal:
		signal := new(SkataSignal)
		signal.Deserialize(data[1:])
		return signal
	}
	return nil
}

// Connection is a high-level abstraction of writing to
// a TCP connection
type Connection struct {
	conn   *net.TCPConn
	Pipe   chan SkataMessage
	closed bool
}

// NewConnection creates a Connection object and returns it
func NewConnection(address, port string) (conn *Connection) {
	conn = new(Connection)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address+":"+port)
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	conn.conn = tcpConn
	conn.Pipe = make(chan SkataMessage)
	go conn.commRoutine()
	return
}

// Close wrapper
func (c *Connection) Close() {
	c.closed = true
	c.conn.Close()
}

func (c *Connection) Write(msg SkataMessage) error {
	data := createPacket(msg)
	dataLength := len(data)
	dataLengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(dataLengthBytes, uint64(dataLength))
	// append length to beginning of
	data = append(dataLengthBytes, data...)
	expectedWriteLength := len(data)
	writtenBytes := 0
	for {
		bytesWritten, err := c.conn.Write(data[writtenBytes:])
		if err != nil {
			return err
		}
		writtenBytes += bytesWritten
		if writtenBytes == expectedWriteLength {
			return nil
		}
	}
}

// Send a packet of data
func (c *Connection) commRoutine() {
	for {
		dataLength := make([]byte, 8)
		_, err := c.conn.Read(dataLength)
		if err != nil {
			if c.closed {
				return
			}
			panic(err)
		}
		packetLength := binary.BigEndian.Uint64(dataLength)
		packet := make([]byte, int(packetLength))
		var bytesRead int
		for bytesRead != int(packetLength) {
			n, err := c.conn.Read(packet[bytesRead:])
			if err != nil {
				if c.closed {
					return
				}
				panic(err)
			}
			bytesRead += n
		}
		msg := parsePacket(packet)
		c.Pipe <- msg
	}
}
