// Package drpc implements the Discord RPC protocol.
package drpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"os"
)

// Client is a Discord RPC client.
type Client struct {
	id   string
	conn net.Conn
}

// New creates a new Discord RPC client with the given application ID and connects it.
func New(id string) (*Client, error) {
	if id == "" {
		return nil, errors.New("drpc: application id is empty")
	}

	client := &Client{id: id}
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

// Connect connects the client to the Discord RPC server and sends a handshake.
// It returns a nil if the client is already connected.
func (c *Client) Connect() (err error) {
	if c.conn != nil {
		return nil
	}
	if c.conn, err = connect(); err != nil {
		return err
	}

	if err = c.Write(OpHandshake, Handshake{
		V:        "1",
		ClientID: c.id,
	}); err != nil {
		return err
	}

	_, err = c.Read()
	return err
}

// Close closes the client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
		c.conn = nil
	}
	return nil
}

// Write writes a message to the connection.
// It automatically reconnects if the connection is closed.
func (c *Client) Write(opcode Opcode, payload interface{}) error {
	if err := c.Connect(); err != nil {
		return err
	}

	v, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err = binary.Write(buf, binary.LittleEndian, opcode); err != nil {
		return err
	}
	if err = binary.Write(buf, binary.LittleEndian, uint32(len(v))); err != nil {
		return err
	}
	if _, err = buf.Write(v); err != nil {
		return err
	}

	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		if err := c.Close(); err != nil {
			return err
		}
		if err := c.Connect(); err != nil {
			return err
		}
	}

	return err
}

// Read reads a message from the connection.
// It automatically reconnects if the connection is closed.
func (c *Client) Read() (*Message, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}

	buf := make([]byte, 512)

	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if n <= 8 {
		return nil, errors.New("drpc: invalid message length")
	}

	msg := &Message{binary.LittleEndian.Uint32(buf[:4]), buf[8:n]}
	switch msg.Opcode {
	// case OpPing:
	// i think it is sent by the client, not the server
	case OpClose:
		if err = c.Close(); err != nil {
			return nil, err
		}
	}

	return msg, nil
}

// SetActivity sets the player's activity.
func (c *Client) SetActivity(activity Activity) (err error) {
	if err = c.Write(OpFrame, FrameSetActivity{
		FrameHeader: NewFrameHeader("SET_ACTIVITY"),
		Args: FrameSetActivityArgs{
			PID:      os.Getpid(),
			Activity: activity,
		},
	}); err != nil {
		return err
	}
	_, err = c.Read()
	return err
}

// func (c *Client) Subscribe(event string) (err error) {
// 	if err = c.Write(OpFrame, NewFrameHeaderWithEvent("SUBSCRIBE", event)); err != nil {
// 		return err
// 	}
// 	_, err = c.Read()
// 	return err
// }
