package main

import (
	"fmt"
	"log"
	"net"
	"github.com/kwintti/httpfromtcp/internal/request"
 )

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection have been accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("RequestFromReader failed: %v\n", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)	
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)	
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)	
		fmt.Println("Connection is closed")
	 }
}
