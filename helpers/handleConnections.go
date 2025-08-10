package helpers

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type usersInfo struct {
	sync.Mutex
	Number int
	info   map[net.Conn]string
}

var users = usersInfo{
	info: make(map[net.Conn]string),
}

func HandleConnections(conn net.Listener) {
	for {
		client, err := conn.Accept()
		if err != nil {
			fmt.Println("accepting went wrong : ", err)
			return
		}
		Ping_Win_Mess, err := os.ReadFile("/home/faaaziz/Desktop/net-cat/assets/pingwing.txt")
		if err != nil {
			log.Fatal(err)
		}
		client.Write(Ping_Win_Mess)
		var tab []byte
		client.Read(tab)
		name := strings.TrimSpace(string(tab))
		for !Valid_Name(name) {
			fmt.Fprintln(client, "Not a valid name try again ...")
			client.Read(tab)
			name = strings.TrimSpace(string(tab))
		}
		users.Lock()
		users.Number++
		users.info[client] = name
		users.Unlock()

	}
}

func Valid_Name(name string) bool {
	if users.Number != 0 {
		for _, Name := range users.info {
			if !strings.EqualFold(Name, name) {
				return false
			}
		}
	}

	return true
}
