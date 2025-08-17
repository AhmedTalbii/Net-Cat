package helpers

import (
	"net"
	"os"
	"sync"
)

var file *os.File

// usersInfo struct holds the information of active users
type usersInfo struct {
	sync.RWMutex
	info map[net.Conn]string
}

type Messages struct {
	ConSender net.Conn
	NameS     string
	Text      string
}

var (
	users  = usersInfo{info: make(map[net.Conn]string)}
	Msg    = make(chan Messages)
	MsgRLU sync.RWMutex
	Time   = UpdateTime()
)
