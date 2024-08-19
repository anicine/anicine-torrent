package client

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

var (
	ErrNotFound = errors.New("not found")
	ErrBadData  = errors.New("bad data")
	ErrNoData   = errors.New("no data")
)

var (
	logger = slog.Default().WithGroup("[HTTP]")
	proxy  *http.Client
)

func SetProxy(client *http.Client) {
	if client != nil {
		proxy = client
	} else {
		proxy = http.DefaultClient
	}
}

type Args struct {
	Proxy    bool
	Method   string
	Endpoint *url.URL
	Headers  map[string]string
	Body     io.Reader
}
