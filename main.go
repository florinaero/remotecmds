package main

import (
	"fmt"
	"github.com/florinaero/remotecmds/pkg/server"
)

func main() {
	fmt.Println("Start remotecmds...")
	if 1 > 0 {
		server.HandleRequests()
		server.StartServer()
	}
}
