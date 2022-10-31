package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address:    address,
		timeout:    timeout,
		input:      in,
		output:     out,
		connection: nil,
	}
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	input      io.Reader
	output     io.Writer
	connection net.Conn
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}
	t.connection = conn
	return nil
}

func (t *telnetClient) Close() error {
	if t.connection != nil {
		return t.connection.Close()
	}
	return nil
}

func (t *telnetClient) Send() error {
	if t.connection == nil {
		return fmt.Errorf("not connected")
	}
	s := bufio.NewScanner(t.input)
	for s.Scan() {
		_, err := t.connection.Write([]byte(fmt.Sprintf("%s\n", s.Text())))
		if err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}
	return nil
}

func (t *telnetClient) Receive() error {
	if t.connection == nil {
		return fmt.Errorf("not connected")
	}
	s := bufio.NewScanner(t.connection)
	for s.Scan() {
		_, err := t.output.Write([]byte(fmt.Sprintf("%s\n", s.Text())))
		if err != nil {
			return fmt.Errorf("receive error: %w", err)
		}
	}
	return nil
}

// P.S. Author's solution takes no more than 50 lines.
