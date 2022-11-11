package http

import (
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers Headers

	status ParseStatus
}

type Response struct {
}

type ParseStatus int

type Headers map[string]string

func (h Headers) Get(name string) string {
	k := strings.ToLower(name)
	return h[k]
}

func (h Headers) Add(name, value string) {
	k := strings.ToLower(name)
	h[k] = value
}

const (
	ParseBegin = ParseStatus(iota)
	ParseMethod
	ParsePathPre
	ParsePath
	ParseVersionPre
	ParseVersion
	ParseVersionPreDone
	ParseVersionDone
	ParseHeaderNamePre
	ParseHeaderName
	ParseHeaderValuePre
	ParseHeaderValue
	ParseHeaderValueDone
	ParseDonePre
	ParseDone
	ParseError
)

func (req *Request) Feed(buf []byte) (ParseStatus, int) {
	var p, i int
	status := ParseBegin
	if req.status != ParseBegin {
		status = req.status
	}
	var headerName, headerValue string
	for i = 0; i < len(buf); i++ {
		switch status {
		case ParseBegin:
			if isAlpha(buf[i]) {
				p = i
				status = ParseMethod
			}
		case ParseMethod:
			if !isAlpha(buf[i]) {
				req.Method = string(buf[p:i])
				status = ParsePathPre
			}
		case ParsePathPre:
			if isPrintable(buf[i]) {
				p = i
				status = ParsePath
			}
		case ParsePath:
			if !isPrintable(buf[i]) {
				req.Path = string(buf[p:i])
				status = ParseVersionPre
			}
		case ParseVersionPre:
			if isPrintable(buf[i]) {
				p = i
				status = ParseVersion
			}
		case ParseVersion:
			if buf[i] == '\r' {
				status = ParseVersionDone
				continue
			}
			fallthrough
		case ParseVersionDone:
			if buf[i] == '\n' {
				n := i
				if buf[i-1] == '\r' {
					n--
				}
				req.Version = string(buf[p:n])
				status = ParseHeaderNamePre
			}
		case ParseHeaderNamePre:
			if isAlpha(buf[i]) {
				p = i
				status = ParseHeaderName
			} else if buf[i] == '\r' {
				status = ParseDonePre
				continue
			} else if buf[i] == '\n' {
				status = ParseDone
				break
			}
		case ParseHeaderName:
			if isAlpha(buf[i]) || buf[i] == '-' {
				continue
			} else if buf[i] == ':' {
				headerName = string(buf[p:i])
				status = ParseHeaderValuePre
				continue
			} else if buf[i] == '\r' {
				status = ParseDonePre
				continue
			}
			fallthrough
		case ParseDonePre:
			if buf[i] == '\n' {
				status = ParseDone
				break
			}
		case ParseHeaderValuePre:
			if isPrintable(buf[i]) {
				p = i
				status = ParseHeaderValue
			}
		case ParseHeaderValue:
			if buf[i] == '\r' {
				status = ParseHeaderValueDone
				continue
			}
			fallthrough
		case ParseHeaderValueDone:
			if buf[i] == '\n' {
				n := i
				if buf[i-1] == '\r' {
					n--
				}
				headerValue = string(buf[p:n])
				if req.Headers == nil {
					req.Headers = Headers{}
				}
				req.Headers.Add(headerName, headerValue)

				status = ParseHeaderNamePre
			}
		default:
			status = ParseError
			break
		}
	}

	req.status = status

	return status, i
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isPrintable(c byte) bool {
	return c >= '!' && c <= '~'
}
