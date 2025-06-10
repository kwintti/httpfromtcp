package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string


func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	before, _, found := bytes.Cut(data, []byte("\r\n"))
	if !found {
		return 0, false, nil
	}
	if len(before) == 0 && found {
		return 2, true, nil
	}
	BytesRead := len(before) + 2
	fieldName, fieldValue, _ := bytes.Cut(before, []byte(":"))
	if fieldName[len(fieldName)-1] == ' ' {
		return 0, false, fmt.Errorf("Spaces between field-name and colon: %v", err)
	}
	trimmedFieldValue := bytes.Trim(fieldValue, " ")
	trimmedFieldName := bytes.TrimLeft(fieldName, " ")
	h[string(trimmedFieldName)] = string(trimmedFieldValue)
	
	return BytesRead, false, nil
}
