package httpclient

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
)

func NewMockHTTPClient(handler http.Handler, mockError error) *http.Client {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				if mockError != nil {
					return nil, mockError
				}
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli
}
