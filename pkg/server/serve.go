package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"rc/shared"
)

func Serve() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8080")

	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	buffer := make([]byte, 1024)
	log.Println("Listening to udp://localhost:8080")

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			log.Println("error reading packet", err)
			continue
		}

		var state shared.NormalizedGamepad

		buf := bytes.NewReader(buffer[:n])
		err = binary.Read(buf, binary.LittleEndian, &state)

		if err != nil {
			log.Println("error decoding data", err)
			continue
		}

		log.Printf("Got message from %s: %v\n", clientAddr, state)
	}
}
