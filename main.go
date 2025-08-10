package main

import (
	"fmt"
	"net"

	"net_cat/helpers"
)

func main() {
	conn, er := net.Listen("tcp", ":8080")
	if er != nil {
		fmt.Println("error listenning:", er)
		return
	}
	helpers.HandleConnections(conn)
}
