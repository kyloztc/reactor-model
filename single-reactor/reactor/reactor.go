package reactor

import (
	"fmt"
	"net"
	"strings"
	"time"
)


type Reactor struct {
	listener net.Listener
	SocketMap map[string]*ConnFd
	ReqBuf chan string
	ResBuf chan string
	Dispatcher map[ReactorReq]Handler
}

func NewReactor(address string) (*Reactor, error) {
	reactor := new(Reactor)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("listen address %s error|%v\n", address, err)
		return nil, err
	}
	reactor.listener = listener
	reactor.ReqBuf = make(chan string, 10)
	reactor.ResBuf = make(chan string, 10)
	reactor.SocketMap = make(map[string]*ConnFd)
	handlerMap := make(map[ReactorReq]Handler)
	for key, value := range HandlerMap {
		value.SetResBuf(reactor.ResBuf)
		handlerMap[key] = value
	}
	reactor.Dispatcher = handlerMap
	return reactor, nil
}

func (r *Reactor) accept() {
	fmt.Printf("reactor listening...\n")
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			fmt.Printf("accept error|%v\n", err)
			continue
		}
		connFd := NewConnFd(conn, r.ReqBuf)
		r.SocketMap[connFd.Id] = connFd
		fmt.Printf("connect establish|id: %v\n", connFd.Id)
		go connFd.Read()
	}
}

func (r *Reactor) Run() {
	go r.accept()
	for {
		select {
		case req := <- r.ReqBuf:
			infos := strings.Split(req, "|")
			id := infos[0]
			reqTask := infos[1]
			handler := r.Dispatcher[ReactorReq(reqTask)]
			handler.Handle(id, r.ResBuf)
		}
		select {
		case res := <- r.ResBuf:
			infos := strings.Split(res, "|")
			id := infos[0]
			msg := infos[1]
			connFd := r.SocketMap[id]
			fmt.Printf("send rsp to: %v\n", id)
			connFd.Conn.Write([]byte(msg))
		}

		time.Sleep(time.Second)
	}
}