package request

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const BufferSize = 8

type ParserState int

const (
	initialized ParserState = iota
	done
)

type Request struct {	
	RequestLine RequestLine
	ParserState
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
			return nil, fmt.Errorf("Couldn't parse to the buffer")
		}
		copy(buf[0:len(buf)-read], buf[read:])
		ReadToIndex -= read
		buf = buf[:ReadToIndex]
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

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserState == initialized { 
		parsed, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesRead == 0 {
			return 0, nil
		}
		r.RequestLine = parsed
		r.ParserState = done 
		return bytesRead, nil
	}
	if r.ParserState == done {
		return 0, fmt.Errorf("error, trying to read data in a done state") 
	}
	return 0, fmt.Errorf("unknown state")

}
