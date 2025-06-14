package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/quic-go/quic-go"
)

func SendBinaryString(conn interface{}, message string) error {
	b := NewAction("BinaryFrame")
	buf, err := b.EncodeString(message)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}
	switch c := conn.(type) {
	case net.Conn:
		// Send the buffer over the connection
		if _, err := c.Write(buf); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	case quic.Stream:
		if _, err := c.Write(buf); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	default:
		// Handle unsupported connection types
		return fmt.Errorf("unsupported connection type: %T", conn)
	}
	// Successful
	return nil
}

func ReceiveBinaryString(conn interface{}) (string, error) {
	// Header size
	const headerSize = 2

	// Create a buffer to read the first 2 bytes (the length of the message)
	lenBuf := make([]byte, headerSize)

	switch c := conn.(type) {
	case net.Conn:
		// Read exactly 2 bytes for the message length
		if _, err := io.ReadFull(c, lenBuf); err != nil {
			return "", fmt.Errorf("failed to read message length from net.Conn: %w", err)
		}
	case quic.Stream:
		if _, err := io.ReadFull(c, lenBuf); err != nil {
			return "", fmt.Errorf("failed to read message length from quic.Stream: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported connection type: %T", conn)
	}
	b := NewAction("BinaryFrame")
	m, messageBuf, err := b.DecodeString(lenBuf)
	if err != nil {
		return "", fmt.Errorf("failed to decode message length: %w", err)
	}

	switch c := conn.(type) {
	case net.Conn:
		if _, err := io.ReadFull(c, messageBuf); err != nil {
			return "", fmt.Errorf("failed to read message from net.Conn: %w", err)
		}
	case quic.Stream:
		if _, err := io.ReadFull(c, messageBuf); err != nil {
			return "", fmt.Errorf("failed to read message from quic.Stream: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported connection type: %T", conn)
	}

	// Convert the message buffer to a string and return it
	return m, nil
}

func SendBinaryTransportString(conn interface{}, message string, transport byte) error {
	b := NewAction("BinaryFrame")
	buf, err := b.Encode(message, transport)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}
	switch c := conn.(type) {
	case net.Conn:
		// Send the buffer over the connection
		if _, err := c.Write(buf); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	case quic.Stream:
		if _, err := c.Write(buf); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	default:
		// Handle unsupported connection types
		return fmt.Errorf("unsupported connection type: %T", conn)
	}
	fmt.Println("SSS222")
	// Successful
	return nil
}

func ReceiveBinaryTransportString(conn interface{}) (string, byte, error) {
	// // Header size
	// const headerSize = 3

	// // Create a buffer to read the first 2 bytes (the length of the message)
	// lenBuf := make([]byte, headerSize)
	var lenBuf [2]byte
	fmt.Println("S")
	switch c := conn.(type) {
	case net.Conn:
		// Read exactly 2 bytes for the message length
		if _, err := io.ReadFull(c, lenBuf[:]); err != nil {
			return "", 0, fmt.Errorf("failed to read message length from net.Conn: %w", err)
		}
	case quic.Stream:
		if _, err := io.ReadFull(c, lenBuf[:]); err != nil {
			return "", 0, fmt.Errorf("failed to read message length from quic.Stream: %w", err)
		}
	default:
		return "", 0, fmt.Errorf("unsupported connection type: %T", conn)
	}
	fmt.Println("SSS111")
	messageLength := binary.BigEndian.Uint16(lenBuf[:])
	if messageLength == 0 {
		fmt.Println("Errr1")
		return "", 0, fmt.Errorf("message length is zero")
	}

	// Allocate a buffer for the message
	messageBuf := make([]byte, messageLength+1) // +1 for transport byte

	// Read the actual message
	switch c := conn.(type) {
	case net.Conn:
		if _, err := io.ReadFull(c, messageBuf); err != nil {
			fmt.Printf("Errr2: %v\n", err)
			return "", 0, fmt.Errorf("failed to read message from net.Conn: %w", err)
		}
	case quic.Stream:
		if _, err := io.ReadFull(c, messageBuf); err != nil {
			fmt.Println("Errr23")
			return "", 0, fmt.Errorf("failed to read message from quic.Stream: %w", err)
		}
	default:
		return "", 0, fmt.Errorf("unsupported connection type: %T", conn)
	}
	fmt.Println("SS")
	// Decode the message
	b := NewAction("BinaryFrame")
	message, transport, err := b.Decode(append(lenBuf[:], messageBuf...))
	if err != nil {
		fmt.Println("Errr")
		return "", 0, fmt.Errorf("failed to decode message: %w", err)
	}
	fmt.Println("SSS")
	return message, transport, nil
}

// SendPort sends the port number as a 2-byte big-endian unsigned integer.
func SendBinaryInt(conn net.Conn, port uint16) error {
	// Create a 2-byte slice to hold the port number
	// buf := make([]byte, 2)
	var buf [2]byte

	// Encode the port number as a big-endian 2-byte unsigned integer
	binary.BigEndian.PutUint16(buf[:], port)

	// Send the 2-byte buffer over the connection
	if _, err := conn.Write(buf[:]); err != nil {
		return fmt.Errorf("failed to send port number %d: %w", port, err)
	}

	// Successful
	return nil
}

// ReceivePort reads a 2-byte big-endian unsigned integer directly from the connection
func ReceiveBinaryInt(conn net.Conn) (uint16, error) {
	var port uint16
	// Use binary.Read to read the port directly from the connection
	err := binary.Read(conn, binary.BigEndian, &port)
	if err != nil {
		return 0, fmt.Errorf("failed to read port number from connection: %w", err)
	}

	// Successful
	return port, nil
}

func SendBinaryByte(conn interface{}, message byte) error {
	// Create a 1-byte buffer and send the message
	messageBuf := [1]byte{message}

	switch c := conn.(type) {
	case net.Conn:
		if _, err := c.Write(messageBuf[:]); err != nil {
			return fmt.Errorf("failed to read message from net.Conn: %w", err)
		}
	case quic.Stream:
		if _, err := c.Write(messageBuf[:]); err != nil {
			return fmt.Errorf("failed to read message from net.Conn: %w", err)
		}
	default:
		return fmt.Errorf("unsupported connection type: %T", conn)
	}

	// Successful
	return nil
}

func ReceiveBinaryByte(conn net.Conn) (byte, error) {
	var messageBuf [1]byte

	switch c := conn.(type) {
	case net.Conn:
		if _, err := io.ReadFull(c, messageBuf[:]); err != nil {
			return 0, fmt.Errorf("failed to read message from net.Conn: %w", err)
		}
	case quic.Stream:
		if _, err := io.ReadFull(c, messageBuf[:]); err != nil {
			return 0, fmt.Errorf("failed to read message from quic.Stream: %w", err)
		}
	default:
		return 0, fmt.Errorf("unsupported connection type: %T", conn)
	}
	// Convert the message buffer to a string and return it
	return messageBuf[0], nil
}
