package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main()  {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Coudln't resolve address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Couldn't create connection: %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Couldn't read input: %v", err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatalf("Couldn't write to UDP: %v", err)
		}
	}

	
}
