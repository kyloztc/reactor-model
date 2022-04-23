package reactor

import (
	"fmt"
	"time"
)

type ReactorReq string

const (
	TimeSleepReq ReactorReq = "timeSleep"
)

var HandlerMap map[ReactorReq]Handler

type Handler interface {
	IsFree() bool
	Handle(connId string, resBuf chan string)
	SetResBuf(ch chan string)
}

type baseHandler struct {
	resBuf chan string
	isFree bool
}

type TimeSleepHandler struct {
	baseHandler
}

func newTimeSleepHandler() *TimeSleepHandler {
	return &TimeSleepHandler{baseHandler{
		isFree: true,
	}}
}

func (t *TimeSleepHandler) IsFree() bool {
	return t.isFree
}

func (t *TimeSleepHandler) Handle(connId string, resBuf chan string) {
	startTime := time.Now().Format("2006-01-02 15:04:05")
	time.Sleep(time.Second * 2)
	endTime := time.Now().Format("2006-01-02 15:04:05")
	info := fmt.Sprintf("time sleep from: %s to %s", startTime, endTime)
	fmt.Printf("time sleep finish\n")
	info = connId + "|" + info
	fmt.Printf("info: %v\n", info)

	t.resBuf <- info
	return
}

func (t *TimeSleepHandler) SetResBuf(ch chan string) {
	t.resBuf = ch
}

func init() {
	HandlerMap = make(map[ReactorReq]Handler)
	HandlerMap[TimeSleepReq] = newTimeSleepHandler()
}
