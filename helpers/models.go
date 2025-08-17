package helpers

import (
	"net"
	"os"
	"sync"
)

var file *os.File

// usersInfo struct holds the number of active users and their information
type usersInfo struct {
	sync.RWMutex
	info map[net.Conn]string
}

type Messages struct {
	ConSender net.Conn
	NameS     string
	Text      string
}

// broadcast channel is used to send messages to all connected clients
var (
	users  = usersInfo{info: make(map[net.Conn]string)}
	Msg    = make(chan Messages)
	MsgRLU sync.RWMutex
	Time   = UpdateTime()
)
