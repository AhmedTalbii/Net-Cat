package helpers

import (
	"net"
	"sync"
)

// usersInfo struct holds the information of active users
type usersInfo struct {
	sync.RWMutex
	info map[net.Conn]string
}

type Messages struct {
	ConSender net.Conn
	NameS     string
	Text      string
	Normal    bool
}

var (
	Msg    = make(chan Messages)
	users  = usersInfo{info: make(map[net.Conn]string)}
	MsgRLU sync.RWMutex
	Time   = UpdateTime()
)
