package response

import (
	"bytes"
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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.StatusOfWriter == Init {
		if statusCode == 200 {
			w.Statusline = "HTTP/1.1 200 OK\r\n"
		}
		if statusCode == 400 {
			w.Statusline = "HTTP/1.1 400 BadRequest\r\n"
		}
		if statusCode == 500 {
			w.Statusline = "HTTP/1.1 500 Internal Server Error\r\n"
		}
	} else {
		return fmt.Errorf("Writer in the wrong state")
	}
	w.StatusOfWriter = Statusline 
	return nil 
}

func GetDefaultHeaders(n int) headers.Headers {
	headersSet := headers.NewHeaders()
	headersSet["Content-Length"] = strconv.Itoa(n) 
	headersSet["Connection"] = "close"
	headersSet["Content-Type"] = "text/html"

	return headersSet
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.StatusOfWriter == Statusline {
		w.Headers = headers

	} else {
		return fmt.Errorf("Writer state is wrong")
	}
	w.StatusOfWriter = Headers
	return nil	
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.StatusOfWriter == Headers {
		w.Body = p
	} else {
		return 0, fmt.Errorf("Writer state is wrong")
	}

	w.StatusOfWriter = Body
	return len(p), nil
}

func (w *Writer) Flush() error {
	if w.StatusOfWriter == Body {
		_, err := w.Buf.Write([]byte(w.Statusline))
		if err != nil {
			return err
		}
		for k, v := range w.Headers {
			_, err = w.Buf.Write([]byte(k + ": " + v + "\r\n"))
			if err != nil {
				return err
			}
		}
		_, err = w.Buf.Write([]byte("\r\n"))
		if err != nil {
			return err
		}
		_, err = w.Buf.Write(w.Body)
		if err != nil {
			return err
		}

	} else {
		return fmt.Errorf("Writer state is wrong")
	}
	w.StatusOfWriter = done
	
	return nil
}

type StatusOfWriter int

const (
	Init StatusOfWriter = iota
	Statusline 
	Headers
	Body
	done
)

type Writer struct {
	StatusOfWriter StatusOfWriter
	Statuscode StatusCode 
	Statusline string
	Headers headers.Headers
	Body []byte 
	Buf io.Writer
	RespFullySent bool
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	trimmed := bytes.TrimRight(p, "\r\n")
	length := len(trimmed)
	lengthOfChunk := fmt.Sprintf("%x", length)
	lengthOfP := []byte(lengthOfChunk + "\r\n")
	_, err := w.ChunkWriter(lengthOfP)
	if err != nil {
		return 0, err 
	}
	_, err = w.ChunkWriter(append(trimmed, []byte("\r\n")...))
	if err != nil {
		return 0, err 
	}

	return len(p), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	_, err := w.ChunkWriter([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}
	return 0, nil 
}

func (w *Writer) ChunkWriter(c []byte) (int, error) {
	n, err := w.Buf.Write(c)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) ChunckedFlush() error {
	_, err := w.Buf.Write([]byte(w.Statusline))
	if err != nil {
		return err
	}
	for k, v := range w.Headers {
		headerLine := []byte(k + ": " + v + "\r\n")
		_, err = w.Buf.Write(headerLine)
		if err != nil {
			return err
		}
	}
	finalCRLF := []byte("\r\n")
	_, err = w.Buf.Write(finalCRLF)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for k, v := range h {
		_, err := w.Buf.Write([]byte(k + ": " + v + "\r\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
