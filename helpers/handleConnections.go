package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

// HandleConnections accepts incoming connections and manages them
// Create or open the history file for storing chat history
// Accept incoming connections
// Handle each client connection concurrently
func HandleConnections(listner net.Listener) {
	defer listner.Close()
	var err error

	file, err = os.Create("assets/history.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("Problem in server : ", err)
				return
			}
			continue
		}

		users.Lock()
		if len(users.info) >= 2 {
			fmt.Fprint(conn, "the chat room is currently full, please try again later\n")
			conn.Close()
			users.Unlock()
			continue
		}
		users.info[conn] = ""
		users.Unlock()

		go ClientInfo(conn)

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
func ClientInfo(conn net.Conn) {
	defer conn.Close()
	Ping_Win_Mess := "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		" |    `.       | `' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     `-'       `--'\n" +
		"[ENTER YOUR NAME]:"
	conn.Write([]byte(Ping_Win_Mess))

	name := ""
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		name = scanner.Text()
		if err := Valid_Name(name); err != nil {
			fmt.Fprint(conn, err)
			continue
		}
		break
	}

	users.RLock()
	users.info[conn] = name
	users.RUnlock()

	history, err := os.ReadFile("assets/history.txt")
	if err != nil {
		fmt.Println("Error reading history:", err)
	}
	conn.Write(history)

	broadcast <- Messages{
		sender:  users.info[conn],
		message: fmt.Sprintf("%s has joined our chat...\n", users.info[conn]),
	}
}

// Valid_Name checks if the provided name is valid (non-empty, under 15 characters, and not already taken)
// Ensure the name is unique (case-insensitive check)
func Valid_Name(name string) error {
	if name == "" || len(name) > 15 {
		return errors.New("invalid name length (1-15), try again:\n[ENTER YOUR NAME]: ")
	}

	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return errors.New("invalid name, only letters allowed, try again:\n[ENTER YOUR NAME]: ")
		}
	}

	for _, v := range users.info {
		if strings.EqualFold(v, name) {
			return errors.New("name already exists, try again:\n[ENTER YOUR NAME]: ")
		}
	}
	return nil
}
