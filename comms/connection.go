package comms

import (
	"encoding/binary"
	"fmt"
	"net"
	"skata/common"
	"time"
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

// Listener listens for TCP connections and attempts to establish the node type
type Listener struct {
	*net.TCPListener
	ConnectionChan   <-chan *Connection
	internalConnChan chan *Connection
}

// NewListener is the factory method for creating a Listener
func NewListener(addr string) *Listener {
	listener := new(Listener)
	listenAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	baseListener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	listener.TCPListener = baseListener
	listener.internalConnChan = make(chan *Connection)
	listener.ConnectionChan = listener.internalConnChan
	go listener.ListenAndAccept()
	return listener
}

func (l *Listener) handleNewConnection(conn *net.TCPConn) {
	skataConn := new(Connection)
	skataConn.CreateFromTCPConn(conn)
	select {
	case response := <-skataConn.Pipe:
		hello, ok := response.(*SkataSignal)
		if !ok {
			conn.Close()
			return
		}
		skataConn.Source = hello.source
		l.internalConnChan <- skataConn
	case <-time.Tick(time.Second * 5):
		conn.Close()
		return
	}
}

func (l *Listener) ListenAndAccept() {
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go l.handleNewConnection(conn)
	}
}

// Connection is a high-level abstraction of writing to
// a TCP connection
type Connection struct {
	Source common.SkataNodeID
	conn   *net.TCPConn
	Pipe   chan SkataMessage
	closed bool
}

// NewConnection creates a Connection object and returns it
func NewConnection(address, port string, source common.SkataNodeID) (conn *Connection) {
	conn = new(Connection)
	conn.Source = source
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address+":"+port)
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	conn.CreateFromTCPConn(tcpConn)
	return
}

// CreateFromTCPConn takes an existing connection
// and injects it into the Connection object
func (c *Connection) CreateFromTCPConn(conn *net.TCPConn) {
	c.initConnection(conn)
	return
}

func (c *Connection) initConnection(conn *net.TCPConn) {
	c.conn = conn
	c.Pipe = make(chan SkataMessage)
	go c.commRoutine()
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
