//go:build windows

package client

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"time"
)

func Connect() {
	serverIp := os.Args[1:]

	if len(serverIp) == 0 {
		log.Println("server-ip is required")
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", serverIp[0])
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 60hz
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for range ticker.C {
		state := GetControllerState()

		if state == nil {
			continue
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
