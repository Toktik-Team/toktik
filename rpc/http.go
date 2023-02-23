package rpc

import (
	"io"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	client http.Client
)

type HttpRequestOptClosure func(*http.Request)

func WithAuthorizationHeader(authorization string) HttpRequestOptClosure {
	return func(req *http.Request) {
		req.Header.Add("Authorization", authorization)
	}
}

func HttpRequest(method, url string, body io.Reader, opt ...HttpRequestOptClosure) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for _, item := range opt {
		item(req)
	}

	return client.Do(req)
}

func init() {
	client = http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}
