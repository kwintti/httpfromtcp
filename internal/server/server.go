package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
)



type Server struct {
	Listener net.Listener
	serving atomic.Bool
}

func Serve(port int) (*Server, error) {
	portString := strconv.Itoa(port)
	l, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, fmt.Errorf("Error listening to the port: %v", err)
	}
	server := &Server{
		Listener:l,
	}
	server.serving.Store(true)
	go server.listen()
	return server, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if !s.serving.Load(){
				break
			}
			fmt.Printf("Connection failed\n")
		}
		go func(c net.Conn){
			s.handle(c)	
		}(conn)
	}
}

func (s *Server) Close() error{
	err := s.Listener.Close()
	if err != nil {
		return fmt.Errorf("Couldn't close listener: %v", err)
	}
	s.serving.Store(false)
	return nil
}

func (s *Server) handle(conn net.Conn) {
	resp := []byte(
	"HTTP/1.1 200 OK\r\n" +
	"Content-Type: text/plain\r\n\r\n" + 
	"Hello World!")
	conn.Write(resp)
	conn.Close()
}
