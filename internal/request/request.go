package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/kwintti/httpfromtcp/internal/headers"
)

const BufferSize = 8

type ParserState int

const (
	initialized ParserState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	done
)

type Request struct {	
	RequestLine RequestLine
	ParserState
	Headers headers.Headers 
	Body []byte
}

type RequestLine struct {
	HttpVersion string
	RequestTarget string
	Method string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BufferSize, BufferSize)

	var ReadToIndex = 0
	var NewRequest Request

	NewRequest.ParserState = initialized

	for {
		if NewRequest.ParserState == done {
			break
		}
		if ReadToIndex == len(buf) {
			bufBigger := make([]byte, len(buf)*2,len(buf)*2)
			copy(bufBigger, buf)
			buf = bufBigger
		}
		n, err := reader.Read(buf[ReadToIndex:])
		if err != nil {
			return nil, fmt.Errorf("Reading from the reader failed")
		}
		ReadToIndex += n
		read, err := NewRequest.parse(buf[:ReadToIndex])
		if err != nil {
			return nil, fmt.Errorf("Couldn't parse to the buffer: %v", err)
		}
		copy(buf[0:len(buf)-read], buf[read:])
		ReadToIndex -= read
	}

	return &NewRequest, nil
}

func parseRequestLine(line []byte) (RequestLine, int, error) {
	before, _, found := bytes.Cut(line, []byte("\r\n"))
	if !found {
		return RequestLine{}, 0, nil
	}
	bytesRead := len(before) + 2
	parsed := bytes.Split(before, []byte(" "))
	if len(parsed) != 3 {
		return RequestLine{},bytesRead, fmt.Errorf("length of requestline is not 3")
	}
	matched, err := regexp.MatchString("^[A-Z]+$", string(parsed[0]))
	if err != nil {
		return RequestLine{}, bytesRead, err
	}
	if !matched {
		return RequestLine{},bytesRead, fmt.Errorf("Method is not capitalized")
	}
	result := strings.Compare(string(parsed[2]), "HTTP/1.1")
	if !(result == 0) {
		return RequestLine{},bytesRead, fmt.Errorf("http version is not matching")
	}
	return RequestLine{
		HttpVersion: string(parsed[2][5:]),
		RequestTarget: string(parsed[1]),
		Method: string(parsed[0]),
	},bytesRead, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.ParserState == initialized { 
		parsed, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesRead == 0 {
			return 0, nil
		}
		r.RequestLine = parsed
		r.ParserState = requestStateParsingHeaders 
		return bytesRead, nil
	}
	if r.ParserState == requestStateParsingHeaders {
		if r.Headers == nil {
			r.Headers = headers.NewHeaders()
		}
		n, d, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse header: %v", err)
		}
		if d {
			r.ParserState = requestStateParsingBody 
		}
		return n, nil
	}

	if r.ParserState == requestStateParsingBody {
		val, ok:= r.Headers.Get("Content-Length")
		if !ok {
			r.ParserState = done
			return 0, nil
		}
		contLength, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("Couldn't convert to int: %v", err)
		}

		bytesLeftToRead := contLength - len(r.Body)
		slice := min(len(data), bytesLeftToRead)

		r.Body = append(r.Body, data[:slice]...)  

		if len(r.Body) > contLength {
			return 0, fmt.Errorf("Mismatch in Content-Length")
		}

		if len(r.Body) == contLength {
			r.ParserState = done
			return slice, nil
		}

		if len(r.Body) < contLength {
			return slice, nil
		}
	}

	return 0, nil 
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesRead := 0
	for r.ParserState != done {
		n, err := r.parseSingle(data[totalBytesRead:])
		if err != nil {
			return 0, fmt.Errorf("Couldn't parse sinlge: %v", err)
		}
		totalBytesRead += n
		if n == 0 {
			break
		}
	}
	return totalBytesRead, nil
}
