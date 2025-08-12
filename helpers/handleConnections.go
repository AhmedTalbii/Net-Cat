package helpers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var file *os.File

// usersInfo struct holds the number of active users and their information
type usersInfo struct {
	sync.Mutex
	Number int // number of all the accepted connections
	info   map[net.Conn]string
}

var users = usersInfo{
	info: make(map[net.Conn]string),
}

// broadcast channel is used to send messages to all connected clients
var broadcast = make(chan string)

// HandleConnections accepts incoming connections and manages them
// Create or open the history file for storing chat history
// Accept incoming connections
// Handle each client connection concurrently
func HandleConnections(conn net.Listener) {
	var err error

	file, err = os.Create("history.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

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

// ClientInfo handles the individual client interactions, including name validation and broadcasting messages
// Load the initial message for new clients
// Read and validate the client’s name
// Keep asking for a valid name until provided
// Update user count and store user info
// Allow only 10 users in the chat room at a time
// Load and send chat history to the new client
// Broadcast to all users that a new user has joined
// Reject connection if the chat room is full
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

	for !Valid_Name(name) {
		fmt.Fprint(client, "Invalid Name try again\n[ENTER YOUR NAME]:")
		name, err = reader.ReadString('\n')
		if err != nil {
			client.Close()
			return
		}
		name = strings.TrimSpace(name)
	}

	users.Lock()
	users.Number++
	users.Unlock()

	if users.Number <= 10 {
		users.Lock()
		users.info[client] = name
		users.Unlock()

		history, err := os.ReadFile("history.txt")
		if err != nil {
			fmt.Println("Error reading history:", err)
		}
		client.Write(history)

		broadcast <- fmt.Sprintf("%s has joined our chat...\n", users.info[client])

	} else {

		users.Lock()
		users.Number--
		users.Unlock()
		fmt.Fprint(client, "the chat room is currently full, please try again later\n")
		client.Close()
	}
}

// Valid_Name checks if the provided name is valid (non-empty, under 15 characters, and not already taken)
// Ensure the name is unique (case-insensitive check)
func Valid_Name(name string) bool {
	if name == "" || len(name) > 15 {
		return false
	}

	for _, v := range users.info {
		if strings.EqualFold(v, name) {
			return false
		}
	}
	return true
}
