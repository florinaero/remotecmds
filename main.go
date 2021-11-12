package main

import (
	"fmt"

	"github.com/florinaero/remotecmds/pkg/server"
)

func main() {
	fmt.Println("Start remotecmds...")
	server.HandleRequests()
	server.StartServer()
}