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

func GetDefaultHeaders(contentLen int) headers.Headers {
	headersSet := headers.NewHeaders()
	headersSet["Content-Length"] = strconv.Itoa(contentLen) 
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
}
