// Package api implements an API client for the makerbotd server
package api

import (
	"context"
	"net"
	"net/http"
)

const (
	// DefaultUnixSocketPath is the default path that makerbotd will listen to on a UNIX socket
	DefaultUnixSocketPath = "/var/run/makerbot.sock"
)

// NewClientSocket creates an API client that connects to a
// makerbotd server via a UNIX domain socket.
func NewClientSocket(socket string) *Client {
	httpc := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socket)
			},
		},
	}

	return &Client{
		http:    httpc,
		baseURL: "http://unix",
	}
}

// NewClientTCP creates an API client that connects to a
// makerbotd server via a TCP connection.
//
// `base` should be the base URL, e.g. "http://localhost:6969"
func NewClientTCP(base string) *Client {
	return &Client{
		http:    &http.Client{},
		baseURL: base,
	}
}
