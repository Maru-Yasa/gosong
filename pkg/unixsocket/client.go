package unixsocket

import (
	"fmt"
	"net"
	"time"

	"github.com/Maru-Yasa/gosong/pkg/proto/kvproto"
)

// Client represents a Unix socket client
type Client struct {
	sockFile string
	timeout  time.Duration
}

// NewClient creates a new Unix socket client
func NewClient(sockFile string) *Client {
	return &Client{
		sockFile: sockFile,
		timeout:  5 * time.Second,
	}
}

// SetTimeout sets the timeout for socket operations
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// Send sends a command to the Unix socket server
func (c *Client) Send(data map[string]string) (string, error) {
	// Encode the data using kvproto
	encoded := kvproto.Encode(data)

	// Connect to the Unix socket
	conn, err := net.DialTimeout("unix", c.sockFile, c.timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect to socket: %w", err)
	}
	defer conn.Close()

	// Send the encoded data
	_, err = conn.Write([]byte(encoded))
	if err != nil {
		return "", fmt.Errorf("failed to send data: %w", err)
	}

	// Read the response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(buf[:n]), nil
}

// Ping sends a ping command to the server
func (c *Client) Ping() (bool, error) {
	response, err := c.Send(map[string]string{
		"action": "ping",
	})
	if err != nil {
		return false, err
	}

	return response == "pong\n", nil
}