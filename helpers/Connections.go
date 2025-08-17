package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func HandleConnections(listner net.Listener) {
	defer listner.Close()
	var err error

	file, err = os.Create("assets/history.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	go StartListeningChan()

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
		go HandleClient(conn)
	}
}

func HandleClient(conn net.Conn) {
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

	users.Lock()
	users.info[conn] = name
	users.Unlock()

	history, err := os.ReadFile("assets/history.txt")
	if err != nil {
		fmt.Println("Error reading history:", err)
	}
	conn.Write(history)

	if len(name) != 0 {
		Msg <- Messages{
			ConSender: conn,
			NameS:     name,
			Text:      fmt.Sprintf("%s has joined our chat...", name),
		}
	}

	HandleMessage(conn, name)
	delete(users.info, conn) 
}

func Valid_Name(name string) error {
	if name == "" || len(name) > 15 {
		return errors.New("invalid name length (1-15), try again:\n[ENTER YOUR NAME]: ")
	}

	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return errors.New("invalid name, only letters allowed, try again:\n[ENTER YOUR NAME]: ")
		}
	}

	users.RLock()
	defer users.RUnlock()
	for _, v := range users.info {
		if strings.EqualFold(v, name) {
			return errors.New("name already exists, try again:\n[ENTER YOUR NAME]: ")
		}
	}
	return nil
}
