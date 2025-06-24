package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/kwintti/httpfromtcp/internal/request"
	"github.com/kwintti/httpfromtcp/internal/response"
)



type Server struct {
	Listener net.Listener
	serving atomic.Bool
	handler Handler 
}

func Serve(port int, handler Handler) (*Server, error) {
	portString := strconv.Itoa(port)
	l, err := net.Listen("tcp", ":"+portString)
	if err != nil {
		return nil, fmt.Errorf("Error listening to the port: %v", err)
	}
	server := &Server{
		Listener:l,
		handler: handler,
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
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		WriteError(conn, HandlerError{StatusCode: 400, Message: "Bad Request\n"})
		return
	}
	var buf bytes.Buffer
	errorHand := s.handler(&buf, req)
	if errorHand != nil {
		WriteError(conn, *errorHand)
		return
	}

	response.WriteStatusLine(conn, 200)
	headers := response.GetDefaultHeaders(buf.Len())
	response.WriteHeaders(conn, headers)
	conn.Write(buf.Bytes())

}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

func WriteError(w io.Writer, handlerError HandlerError) {
	fmt.Fprintf(w, "HTTP/1.1 %v %v\r\n", handlerError.StatusCode, http.StatusText(int(handlerError.StatusCode))) 
	fmt.Fprint(w, "\r\n") 
	fmt.Fprintf(w, "%v", handlerError.Message) 
}
