package main

import (
	"fmt"
	"reactor-model/single-reactor/reactor"
)

func main() {
	r, err := reactor.NewReactor("127.0.0.1:9090")
	if err != nil {
		fmt.Printf("new reactor error|%v\n", err)
		return
	}
	r.Run()
}