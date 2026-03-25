package server

import (
	"log"
	"net"
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

		log.Printf("Got message from %s: %s\n", clientAddr, string(buffer[:n]))

		_, err = conn.WriteToUDP(buffer[:n], clientAddr)
		if err != nil {
			log.Printf("Write error: %v", err)
		}
	}
}
