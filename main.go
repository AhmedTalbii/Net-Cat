package main

import (
	"fmt"
	"net"
	"os"

	logic "net_cat/logic"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 { fmt.Println("[USAGE]: ./TCPChat $port"); return }

	port := ":8989"
	if len(args) == 1 { port = ":" + args[0] }

	listner, err := net.Listen("tcp", port)
	if err != nil { fmt.Println("Error at listening :", err); return }
	defer listner.Close()

	addr := listner.Addr().String()
	_, por, _ := net.SplitHostPort(addr)
	fmt.Println("server is started at port :", por)

	logic.HandleConnections(listner)
}
