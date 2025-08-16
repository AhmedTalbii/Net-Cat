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
	Number int // number of all the accepted connections
	info   map[net.Conn]string
}

type Messages struct {
	sender  string
	message string
}

// broadcast channel is used to send messages to all connected clients
var (
	users     = usersInfo{info: make(map[net.Conn]string)}
	broadcast = make(chan Messages)
)
