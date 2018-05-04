package network

import (
	"bufio"
	"net"
)

// TODO - Make own client reader for non buffered reading
// Improve Packet send and recieve to handle headers

type Client struct {
	conn    net.Conn
	Send    chan Packet
	Recieve chan Response
	Error   chan error
}

func (c Client) sendPacket() {
	writer := bufio.NewWriter(c.conn)
	for {
		packet := <-c.Send
		encoded := packet.Encode()
		_, err := writer.Write(encoded)
		if err != nil {
			c.Error <- err
		}
		err = writer.Flush()
		if err != nil {
			c.Error <- err
		}
	}
}

func (c Client) recieveResponse() {
	reader := bufio.NewReader(c.conn)
	for {
		response := <-c.Recieve
		for {
			buf := make([]byte, 4096)
			read, err := c.conn.Read(buf)

			if err != nil {
				c.Error <- err
			}
			if read == 0 {
				break
			}
		}
		err := response.Decode(reader)
		if err != nil {
			c.Error <- err
		}
		c.Error <- nil
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewClient(address string) (*Client, error) {
	var client *Client
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return client, err
	}
	client = &Client{
		conn:    connection,
		Send:    make(chan Packet),
		Recieve: make(chan Response),
		Error:   make(chan error),
	}
	go client.sendPacket()
	go client.recieveResponse()

	return client, nil
}
