package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
	"rc/shared"
)

func Serve() {
	serverIp := os.Args[1:]

	if len(serverIp) == 0 {
		log.Println("binding address is required")
		os.Exit(1)
	}

	addr, err := net.ResolveUDPAddr("udp", serverIp[0])
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	buffer := make([]byte, 1024)

	InitPins()

	log.Println("Listening to udp://" + serverIp[0])
	for {
		n, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			log.Println("error reading packet", err)
			continue
		}

		var state shared.NormalizedGamepad

		buf := bytes.NewReader(buffer[:n])
		err = binary.Read(buf, binary.LittleEndian, &state)

		log.Println("buffer size", n)
		if err != nil {
			log.Println("error decoding data", err)
			continue
		}

		ApplyControls(&state)
	}
}
