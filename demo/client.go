package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"

	"github.com/daemonl/connor"
)

var brokerAddress string

func init() {
	flag.StringVar(&brokerAddress, "broker", "127.0.0.1:5556", "broker's client address")
}

func main() {
	go do("A")
	go do("B")
	do("C")
}

func logf(s string, p ...interface{}) {
	log.Printf(s+"\n", p...)
}

func do(name string) {
	for {
		conn := broker.ClientConnection(brokerAddress)
		// Establish a connection, retry indefinately.
		logf("Connected %s [%s]", name, brokerAddress)

		scanner := bufio.NewScanner(conn)
		for i := 0; i < 4; i++ {
			n, err := conn.Write([]byte(fmt.Sprintf("%d\n", i)))
			if err != nil {
				logf(err.Error())
				continue
			}
			logf("wrote %d bytes on %s", n, name)

			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			logf("%s GOT: %s", name, line)
		}
		conn.Close()
		logf("Connection Ended %s [%s]", name, brokerAddress)
	}
}
