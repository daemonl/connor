package connor

import (
	"fmt"
	"io"
	"net"
	"time"
)

// Logger is nil, set it to enable logging
var Logger interface {
	Log(string)
}

func logf(m string, p ...interface{}) {
	if Logger != nil {
		Logger.Log(fmt.Sprintf(m, p...))
	}
}

// Loop dial infinitely loops until it has a connection
func LoopDial(address string) net.Conn {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", address)
		if err != nil {
			logf(err.Error())
			time.Sleep(time.Second)
			continue
		}
		return conn
	}
}

func TinyHandshakeDial(address string) net.Conn {
	for {
		conn := LoopDial(address)
		err := TinyHandshake(conn, 0x01)
		if err != nil {
			logf(err.Error())
			continue
		}
		return conn
	}
}

func LoopListen(bind string, chanConn chan net.Conn) {
	var listener net.Listener
	var err error
	for {
		logf("Listen on %s", bind)
		listener, err = net.Listen("tcp", bind)
		if err != nil {
			logf(err.Error())
			time.Sleep(time.Second)
			continue
		}

		for {
			conn, err := listener.Accept()
			if err != nil {
				logf(err.Error())
				continue
			}
			go func(conn net.Conn) {
				logf("Got Connection on %s", bind)
				chanConn <- conn
			}(conn)
		}
	}

}

// TinyHandshakeListen binds to an address and sends the connection to the
// chanel, after successful handshake
func TinyHandshakeListen(bind string, chanConn chan net.Conn) {
	localChanConn := make(chan net.Conn)
	go LoopListen(bind, localChanConn)
	for {
		conn := <-localChanConn
		go func(conn net.Conn) {
			err := TinyHandshake(conn, 0x01)
			if err != nil {
				logf(err.Error())
				return
			}
			chanConn <- conn
		}(conn)
	}
}

// Handshake writes one byte, and ensures that byte comes back. It
// is identical on both sides, the byte should be hardcoded as a
// version, not an echo
func TinyHandshake(conn net.Conn, msg byte) error {
	logf("Handshake %s", conn.RemoteAddr())
	n, err := conn.Write([]byte{msg})
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("could not write handshake byte")
	}
	in := []byte{0x00}
	n, err = conn.Read(in) // Read a byte
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("could not read handshake byte")
	}
	if in[0] != msg {
		return fmt.Errorf("handskake byte didn't match %x != %x", in, []byte{msg})
	}
	return nil

}

// BindConnections uses io.Copy in both directions, and closes both
// connections on return
func BindConnections(a net.Conn, b net.Conn) {
	go copyThenClose(a, b)
	go copyThenClose(b, a)
}

func copyThenClose(connA net.Conn, connB net.Conn) {
	io.Copy(connA, connB)
	// The error or number of bytes isn't really important
	// Close both if either fail / close.
	connA.Close()
	connB.Close()
}

type FuncLogger func(string)

func (self FuncLogger) Log(msg string) {
	self(msg)
}
