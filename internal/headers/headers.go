package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
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
	if len(fieldName) == 0 {
		return 0, false, fmt.Errorf("field-name is missing")
	}
	if fieldName[len(fieldName)-1] == ' ' {
		return 0, false, fmt.Errorf("Spaces between field-name and colon: %v", err)
	}
	trimmedFieldValue := string(bytes.Trim(fieldValue, " "))
	trimmedFieldName := string(bytes.TrimLeft(fieldName, " "))

	
	trimmedFieldName = strings.ToLower(trimmedFieldName)

	re := regexp.MustCompile("[a-zA-Z0-9!#$%&'*+-.^_`|~]")
	for _, c := range trimmedFieldName {
		match := re.MatchString(string(c))
		if !match {
			return 0, false, fmt.Errorf("field-name contains invalid character")
		}
	}

	val, ok := h[trimmedFieldName]
	if ok {
		h[trimmedFieldName] = val + ", " + trimmedFieldValue
			return BytesRead, false, nil
	}
	
	h[trimmedFieldName] = trimmedFieldValue
	
	return BytesRead, false, nil
}
