package helpers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type usersInfo struct {
	sync.Mutex
	Number int // number of all the accepted connections
	info   map[net.Conn]string
}

var users = usersInfo{
	info: make(map[net.Conn]string),
}
var broadcast chan string

func HandleConnections(conn net.Listener) {
	for {
		defer conn.Close()
		client, err := conn.Accept()
		if err != nil {
			fmt.Println("accepting went wrong : ", err)
			return
		}

		go ClientInfo(client)

	}
}

func ClientInfo(client net.Conn) {
	Ping_Win_Mess, err := os.ReadFile("assets/pingwing.txt")
	if err != nil {
		fmt.Println("Ping_win_message: ", err)
		return
	}
	client.Write(Ping_Win_Mess)
	reader := bufio.NewReader(client)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		client.Close()
		return
	}

	name = strings.TrimSpace(name)

	users.Lock()
	users.Number++
	users.Unlock()

	if users.Number <= 10 {
		users.Lock()
		users.info[client] = name
		broadcast <- fmt.Sprintf("%s has joined our chat...", users.info[client])
		users.Unlock()

	} else {
		users.Lock()
		users.Number--
		users.Unlock()
		client.Close()
	}
}
