package reactor

import (
	"fmt"
	"net"
	"github.com/google/uuid"
)

type ConnFd struct {
	Id string
	Conn net.Conn
	buf chan string
}

func NewConnFd(conn net.Conn, buf chan string) *ConnFd {
	id := uuid.New().String()
	return &ConnFd{
		Id:   id,
		Conn: conn,
		buf:  buf,
	}
}

func (f *ConnFd) Read() {
	for {
		var buf [1024]byte
		n, err := f.Conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed: %v\n", err)
		}
		msg := string(buf[:n])
		if msg == "bye" {
			fmt.Printf("connect: %s exit\n", f.Id)
			return
		}
		fmt.Printf("recive rsp: %v\n", msg)
		msg = f.Id + "|" + msg
		f.buf <- msg
	}
}