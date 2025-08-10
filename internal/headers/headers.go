package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		//not enough data
		return 0, false, nil
	}
	if idx == 0 {
		//end of headers section
		//consume CRLF
		return 2, true, nil
	}

	line := string(data[:idx])
	line = strings.TrimSpace(line)

	key, val, found := strings.Cut(line, ":")
	if !found {
		return 0, false, fmt.Errorf("invalid header field: '%s'\n", line)
	}
	if key != strings.TrimRight(key, " ") { // my old way was probably risky: if strings.ContainsAny(key, " \t")
		return 0, false, fmt.Errorf("invalid header key: '%s'\n", key)
	}

	val = strings.TrimSpace(val)

	h[key] = val
	return idx + 2, false, nil
}
