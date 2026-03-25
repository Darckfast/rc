//go:build windows

package client

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

func Connect() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	for {
		state := GetControllerState()

		if state == nil {
			break
		}
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, state)
		if err != nil {
			log.Println("error encoding binary", err)
			continue
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("Send failed: %v", err)
			continue
		}
	}
}
