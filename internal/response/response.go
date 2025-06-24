package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kwintti/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK StatusCode = 200 
	BadRequest = 400 
	InternalServerError = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	if statusCode == 200 {
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return fmt.Errorf("Couldn't write: %v", err)
		}
	}
	if statusCode == 400 {
		_, err := w.Write([]byte("HTTP/1.1 400 BadRequest\r\n"))
		if err != nil {
			return fmt.Errorf("Couldn't write: %v", err)
		}
	}
	if statusCode == 500 {
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return fmt.Errorf("Couldn't write: %v", err)
		}
	}
	_, err := w.Write([]byte(""))
	if err != nil {
		return fmt.Errorf("Couldn't write: %v", err)
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headersSet := headers.NewHeaders()
	headersSet["Content-Length"] = strconv.Itoa(contentLen) 
	headersSet["Connection"] = "close"
	headersSet["Content-Type"] = "text/plain"

	return headersSet
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(k+": "+v+"\r\n"))
		if err != nil {
			return fmt.Errorf("Couldn't write header: %v", err)
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("Couldn't write header: %v", err)
	}

	return nil	
}
