package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kwintti/httpfromtcp/internal/headers"
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
	if strings.HasPrefix(target, "/httpbin") {
		requestTarget := strings.TrimPrefix(target, "/httpbin")
		resp, err := http.Get("https://httpbin.org" + requestTarget)
		if err != nil {
			log.Println(err)
		}
		w.WriteStatusLine(200)
		resp.Header.Del("Content-Length")
		resp.Header.Add("Transfer-Encoding", "chunked")
		adjustedHeaders := headers.NewHeaders()
		for k := range resp.Header {
			adjustedHeaders[k] = resp.Header.Get(k)
		}
		w.WriteHeaders(adjustedHeaders)
		w.ChunckedFlush()
		buf := make([]byte, 999)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				w.WriteChunkedBody(buf[:n])
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(err)
				break
			}
		}
		w.WriteChunkedBodyDone()
		w.RespFullySent = true	
		resp.Body.Close()
		
		return
	}
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
