package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	for line := range getLinesChannel(f) {
		fmt.Printf("read: %s \n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)
	go func(){
		defer close(linesChan)
		defer f.Close()
		currentLine := ""
		slice := make([]byte, 8)
		for {
			n, err := f.Read(slice)
			if err != nil {
				if err == io.EOF {
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
		if currentLine != "" {
			linesChan <- currentLine
		}

	}()
	return linesChan
}


