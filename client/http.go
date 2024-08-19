package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

func bodyReader(body io.Reader) io.Reader {
	if body == nil {
		return nil
	}

	if buf, ok := body.(*bytes.Buffer); ok {
		return bytes.NewReader(buf.Bytes())
	}

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, body)
	return bytes.NewReader(buf.Bytes())
}

func Do(ctx context.Context, args Args) (*http.Response, error) {
	if args.Endpoint == nil {
		return nil, ErrBadData
	}
	var (
		client *http.Client = http.DefaultClient
		req                 = new(http.Request)
		resp                = new(http.Response)
		err    error
	)

	if args.Proxy {
		if proxy != nil {
			client = proxy
		} else {
			client = http.DefaultClient
		}
	}

	for i := 0; i < 5; i++ {
		if i > 2 {
			client = http.DefaultClient
		}

		req, err = http.NewRequestWithContext(ctx, args.Method, args.Endpoint.String(), bodyReader(args.Body))
		if err != nil {
			logger.Error("cannot make request", "link", args.Endpoint, "error", err)
			continue
		}

		req.Header.Add("User-Agent", UserAgent)
		for k, v := range args.Headers {
			req.Header.Add(k, v)
		}

		resp, err = client.Do(req)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil, context.Canceled
			}
			logger.Error("cannot get response", "link", args.Endpoint, "error", err)
			continue
		}

		logger.Info("accepted response", "link", args.Endpoint, "code", resp.StatusCode)

		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
		case http.StatusNotModified:
			return resp, nil
		case http.StatusNotFound:
			defer resp.Body.Close()
			return nil, ErrNotFound
		case http.StatusBadRequest:
			defer resp.Body.Close()
			return nil, ErrNotFound
		case http.StatusFound:
			endpoint, err := url.Parse(resp.Header.Get("Location"))
			if err != nil {
				return nil, ErrNotFound
			}
			args.Endpoint = endpoint
			continue
		case http.StatusMovedPermanently:
			endpoint, err := url.Parse(resp.Header.Get("Location"))
			if err != nil {
				return nil, ErrNotFound
			}
			args.Endpoint = endpoint
			continue
		case http.StatusForbidden:
			time.Sleep(350 * time.Millisecond)
			continue
		case http.StatusTooManyRequests:
			time.Sleep(30 * time.Second)
			continue
		}
	}
	if resp != nil {
		resp.Body.Close()
	}

	logger.Error("failed to complete the operation", "link", args.Endpoint, "error", err)

	return nil, ErrNoData
}
