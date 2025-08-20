package logic

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// this fucntion does multiple tasks : create the history file if it dosn't exist, start listening on channel, accepts clients (<10), also handles clients/messages.
func HandleConnections(listner net.Listener) {
	file, err := os.Create("assets/history.txt")
	if err != nil { fmt.Println("Error creating file:", err); return }
	file.Close()
	go StartListeningChan()

	for {
		conn, err := listner.Accept()
		if err != nil { if errors.Is(err, net.ErrClosed) { fmt.Println("Problem in server : ", err); return }; continue }

		users.Lock()
		if len(users.info) >= 10 {
			fmt.Fprint(conn, "the chat room is currently full, please try again later\n")
			conn.Close()
			users.Unlock()
			continue
		}
		users.info[conn] = ""
		users.Unlock()

		go WelcomeCLient(conn)
	}
}
