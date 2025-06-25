package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kwintti/httpfromtcp/internal/request"
	"github.com/kwintti/httpfromtcp/internal/response"
	"github.com/kwintti/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request){
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
	bodyText := []byte(`
	<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
	`)
	w.WriteStatusLine(400)
	headers := response.GetDefaultHeaders(len(bodyText))
	w.WriteHeaders(headers)
	w.WriteBody(bodyText)
	return
	}
	if target == "/myproblem" {
	bodyText := []byte(`
	<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`)
	w.WriteStatusLine(500)
	headers := response.GetDefaultHeaders(len(bodyText))
	w.WriteHeaders(headers)
	w.WriteBody(bodyText)
	return
	}

	bodyText := []byte(`
	<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
	`)
	w.WriteStatusLine(200)
	headers := response.GetDefaultHeaders(len(bodyText))
	w.WriteHeaders(headers)
	w.WriteBody(bodyText)
	return
}
