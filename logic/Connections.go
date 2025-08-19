package helpers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
)

// this fucntion does multiple tasks : create the history file if it dosn't exist, start listening on channel, accepts clients (<10), also handles clients/messages.
func HandleConnections(listner net.Listener) {
	file, err := os.Create("assets/history.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	file.Close()
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
		if len(users.info) >= 10 {
			fmt.Fprint(conn, "the chat room is currently full, please try again later\n")
			conn.Close()
			users.Unlock()
			continue
		}
		users.info[conn] = ""
		users.Unlock()

		go HandleCLient(conn)
	}
}

// this is a function that write the pingWing Message to the user asking him for entering his name, then it scan the name and pass it to the map and to the channel, then it read the history data if exist
func HandleCLient(con net.Conn) {
	defer con.Close()
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
		"     `-'       `--'\n"

	con.Write([]byte("\033[0;93m" + Ping_Win_Mess + "\033[0m" + "[ENTER YOUR NAME]:"))

	name := ""
	scanner := bufio.NewScanner(con)
	for scanner.Scan() {
		name = scanner.Text()
		if err := Valid_Name(name); err != nil {
			fmt.Fprint(con, err)
			continue
		}
		break
	}

	users.Lock()
	users.info[con] = name
	users.Unlock()

	history, err := os.ReadFile("assets/history.txt")

	if err != nil {
		fmt.Println("Error reading history:", err)
		return
	} else {
		con.Write(history)
	}

	if len(name) != 0 {
		Msg <- Messages{
			ConSender: con,
			NameS:     name,
			Text:      fmt.Sprintf("%s has joined our chat...", name),
			Normal:    true,
		}
	}
	Handlemessage(con, name)
	delete(users.info, con)
}

// this is a fucntion that scan the user massages and pass it to the (Msg) channel, and to the history file
func Handlemessage(conn net.Conn, name string) {
	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			Msg <- Messages{
				ConSender: conn,
				NameS:     name,
				Text:      fmt.Sprintf("%s has left our chat...", name),
				Normal:    true,
			}
			return
		}
		text := scanner.Text()
		IsNormal := true

		if text == "" {
			SendMessage(conn, "", "31", "Can't send Empty message\n")
			SendMessage(conn, "", "36", "["+UpdateTime()+"]"+"["+name+"]:")
			continue
		} else if er := ValidMessage(text); er != nil {
			switch er.Error() {
			case "The_user_send_joined_or_left" :
				IsNormal = false
			case "out ascii":
				SendMessage(conn, "", "31", "can't send characters out of the ascii range\n")
				SendMessage(conn, "", "36", "["+UpdateTime()+"]"+"["+name+"]:")
			case "large_msg":
				SendMessage(conn, "", "31", "message too large\n")
				SendMessage(conn, "", "36", "["+UpdateTime()+"]"+"["+name+"]:")
			}
			continue
		}

		f, err := os.OpenFile("assets/history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o777)
		if err != nil {
			return
		}
		m := fmt.Sprintf("\n\033[35m[%s][%s]:\033[0m%s", UpdateTime(), name, text)
		f.WriteString(m[1:] + "\n")
		f.Close()
		Msg <- Messages{
			ConSender: conn,
			NameS:     name,
			Text:      text,
			Normal:    IsNormal,
		}
	}
}
