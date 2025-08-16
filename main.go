package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	"net_cat/helpers"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	port := ":8989"
	if len(args) == 1 {
		port = args[0]
	}
	listner, err := net.Listen("tcp", port)
	addr := listner.Addr().String()
	_, por, _ := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println("Error at listening :", err)
		return
	}
	defer listner.Close()
	fmt.Println("port is :", por)
	for {

		conn, err := listner.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("error in the server :", err)
				return
			}
			continue
		}
		go helpers.HandleConnections(conn)
	}
}
