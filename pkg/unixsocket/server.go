package unixsocket

import (
	"fmt"
	"net"
	"os"

	"github.com/Maru-Yasa/gosong/pkg/proto/kvproto"
)

// Server represents a Unix socket server
type Server struct {
	sockFile string
	listener net.Listener
}

// NewServer creates a new Unix socket server
func NewServer(sockFile string) (*Server, error) {
	// Remove existing socket file if it exists
	_ = os.Remove(sockFile)

	return &Server{
		sockFile: sockFile,
	}, nil
}

// Start starts the Unix socket server
func (s *Server) Start() error {
	var err error
	s.listener, err = net.Listen("unix", s.sockFile)
	if err != nil {
		return fmt.Errorf("failed to create unix socket: %w", err)
	}

	return nil
}

// Accept accepts incoming connections and handles them with the provided handler
func (s *Server) Accept(handler func(map[string]string) (string, error)) error {
	if s.listener == nil {
		return fmt.Errorf("server not started")
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		go s.handleConnection(conn, handler)
	}
}

// handleConnection handles a single connection
func (s *Server) handleConnection(conn net.Conn, handler func(map[string]string) (string, error)) {
	defer conn.Close()

	// Read data from the connection
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	// Decode the data using kvproto
	data := kvproto.Decode(string(buf[:n]))

	// Handle the command with the provided handler
	response, err := handler(data)
	if err != nil {
		return
	}

	// Send the response back
	_, _ = conn.Write([]byte(response))
}

// Close closes the server and removes the socket file
func (s *Server) Close() error {
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil {
			return err
		}
	}

	return os.Remove(s.sockFile)
}