package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	Conn net.Conn
	Alive chan bool
}

func NewClient() *Client {
	return &Client{
		Alive: make(chan bool, 1),
	}
}

func (c *Client) Connect(address string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		fmt.Printf("conn server failed|err: %v\n", err)
		return
	}
	c.Conn = conn
	go c.readRsp()
}

func (c *Client) Close() {
	c.Conn.Write([]byte("bye"))
	c.Conn.Close()
}

func (c *Client) SendMsg(msg string) {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("send failed: %v\n", err)
		return
	}
}

func (c *Client) ReadFromStdIn() {
	input := bufio.NewReader(os.Stdin)
	for {
		s, _ := input.ReadString('\n')
		s = strings.TrimSpace(s)
		if strings.ToUpper(s) == "Q" {
			c.Alive <- true
			fmt.Printf("get q\n")
			return
		}
		c.SendMsg(s)
	}
}

func (c *Client) readRsp() {
	for {
		select {
		case <- c.Alive:
			break
		default:
			var buf [1024]byte
			n, err := c.Conn.Read(buf[:])
			if errors.Is(err, net.ErrClosed) {
				break
			}
			if err != nil {
				fmt.Printf("read failed: %v\n", err)
				continue
			}
			fmt.Printf("recive rsp: %v\n", string(buf[:n]))
		}
	}
}