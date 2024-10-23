package agent

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	headerDataLength = 4 // bytes
	defaultTimeout   = 5 * time.Second
)

type client struct {
	ip      string
	port    int
	timeout time.Duration
}

type ConnectionError struct {
	text string
	err  error
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("%v: %v", e.text, e.err)
}

func (e ConnectionError) Unwrap() error {
	return e.err
}

func (c client) Request(data []byte) ([]byte, error) {
	addr := fmt.Sprintf("%s:%d", c.ip, c.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)

	if err != nil {
		return nil, ConnectionError{
			text: "could not resolve TCP address",
			err:  err,
		}
	}

	dialer := net.Dialer{Timeout: c.timeout}
	conn, err := dialer.Dial(tcpAddr.Network(), tcpAddr.String())

	if err != nil {
		return nil, ConnectionError{
			text: "could not connect to the server agent",
			err:  err,
		}
	}
	defer conn.Close()

	if err = writeData(conn, data); err != nil {
		return nil, err
	}

	// read response data length
	dataLen, err := readDataLen(conn)

	if err != nil {
		return nil, err
	}

	var rData []byte
	buffer := make([]byte, 256)
	var rLen int

	for {
		len, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("could not read response: %v", err)
		}

		rData = append(rData, buffer[:len]...)
		rLen += len

		if rLen >= dataLen {
			break
		}
	}

	return rData, nil
}

func writeData(writer io.Writer, data []byte) error {
	// First, write sending data length to the two bytes
	header := make([]byte, headerDataLength)
	dataLen := len(data)
	binary.BigEndian.PutUint32(header, uint32(dataLen))

	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("could not send request header: %v", err)
	}

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("could not send request data: %v", err)
	}

	return nil
}

// readDataLen reads first 2 bytes where data length is stored
func readDataLen(reader io.Reader) (int, error) {
	header := make([]byte, headerDataLength)

	if _, err := reader.Read(header); err != nil {
		return 0, fmt.Errorf("could not read response data length: %v", err)
	}

	return int(binary.BigEndian.Uint32(header)), nil
}
