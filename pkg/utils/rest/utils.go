package rest

import (
	"io"
	"net/http"
)

func setHeaders(header http.Header, headers map[string]string) {
	if len(headers) == 0 {
		return
	}

	for key := range headers {
		header.Set(key, headers[key])
	}
}

type Options struct {
	Headers map[string]string
	Body    io.Reader
}
