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
	addr, err := net.ResolveUDPAddr("udp", "10.0.0.35:8080")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for range ticker.C {
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

		conn.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Printf("Send failed: %v", err)
			continue
		}
	}
}
