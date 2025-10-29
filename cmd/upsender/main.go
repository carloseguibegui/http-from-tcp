package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	server, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("error resolving UDP address", err)
	}
	conn, err := net.DialUDP("udp", nil, server)
	if err != nil {
		log.Fatal("error dialing UDP", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		l, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error reading from stdin:", err)
			continue
		}
		_, err = conn.Write([]byte(l))
		if err != nil {
			log.Println("error writing to UDP connection:", err)
			continue
		}
	}
}
