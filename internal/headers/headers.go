package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var rnSep = []byte("\r\n")

func isValidToken(token []byte) bool {
	for _, ch := range token {
		valid := false
		if ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' {
			valid = true
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			valid = true
		}
		if !valid {
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Get(name string) (string, bool) {
	str, found := h[strings.ToLower(name)]
	return str, found
}

func (h Headers) Set(name, value string) {
	name = strings.ToLower(name)
	if v, ok := h[name]; ok {
		h[name] = fmt.Sprintf("%s,%s", v, value)
	} else {
		h[name] = value
	}
}

func (h Headers) Replace(name, value string) {
	name = strings.ToLower(name)
	h[name] = value
}

func (h Headers) Delete(name string) {
	name = strings.ToLower(name)
	delete(h, name)
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range *h {
		cb(n, v)
	}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Parsing invalid field-line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Invalid field-name containing trailing whitespaces")
	}

	return string(name), string(value), nil

}

func (h Headers) Parse(data []byte) (int, bool, error) {
	// Set-Person: lane-loves-go, prime-loves-zig, tj-loves-ocaml
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rnSep)
		if idx == -1 {
			break
		}
		//empty header
		if idx == 0 {
			done = true
			read += len(rnSep)
			break
		}
		end := read + idx
		name, value, err := parseHeader(data[read:end])
		if err != nil {
			return 0, false, err
		}
		if !isValidToken([]byte(name)) {
			return 0, false, fmt.Errorf("Header not valid as token %s", name)
		}
		read += idx + len(rnSep)
		h.Set(name, value)
	}
	return read, done, nil
}
