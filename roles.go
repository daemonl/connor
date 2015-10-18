package broker

import (
	"net"
	"time"
)

func tcpConnection(brokerAddress string) net.Conn {

	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", brokerAddress)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		return conn
	}
}

func ClientConnection(brokerAddress string) net.Conn {
	for {
		conn := tcpConnection(brokerAddress)
		n, err := conn.Read([]byte{0x00}) // Read a byte
		if err != nil || n != 1 {
			time.Sleep(time.Second)
			if conn != nil {
				conn.Close()
			}
			continue
		}
		return conn
	}
}

func WorkerConnection(brokerAddress string) net.Conn {
	return tcpConnection(brokerAddress)
}
