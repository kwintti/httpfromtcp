package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
// 	"time"
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

		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection is closed")
	 }
}

func getLinesChannel(f net.Conn) <-chan string {
	linesChan := make(chan string)
	go func(){
		defer f.Close()
		defer close(linesChan)
		currentLine := ""
		slice := make([]byte, 8)
		for {
			//err := f.SetDeadline(time.Now().Add(time.Second))
			n, err := f.Read(slice)
			if n == 0 {
				toString := string(slice[:n])
				currentLine += toString
				linesChan <- currentLine
				break
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					toString := string(slice[:n])
					currentLine += toString
					linesChan <- currentLine
					break
				}
				log.Fatal(err)
			}
			toString := string(slice[:n])
			part, endPart, found := strings.Cut(toString, "\n")
			currentLine += part
			if found {
				linesChan <- currentLine
				currentLine = endPart
				continue
			}
		}
	}()
	return linesChan
}


