package main

import (
	"fmt"
	"reactor-model/single-reactor/client"
)

func sendMsg2Server(clt *client.Client) {
	clt.Connect("127.0.0.1", 9090)
	clt.SendMsg("timeSleep")
}

func main() {
	cltList := make([]*client.Client, 0)
	for i := 0; i < 100; i ++ {
		clt := client.NewClient()
		cltList = append(cltList, clt)
		go sendMsg2Server(clt)
	}
	var q string
	fmt.Scanf("%s\n", &q)
	if q == "q" {
		for _, clt := range cltList {
			clt.Alive <- true
			clt.Close()
		}
		return
	}
}
