package logic

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// this is a function that write the pingWing Message to the user asking him for entering his name, then it scan the name and pass it to the map and to the channel, then it read the history data if exist
func WelcomeCLient(con net.Conn) {
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
		if err := Valid_Name(name); err != nil { fmt.Fprint(con, err); continue }
		break
	}

	users.Lock()
	users.info[con] = name
	users.Unlock()

	history, err := os.ReadFile("assets/history.txt")

	if err != nil { fmt.Println("Error reading history:", err); return
	} else { con.Write(history) }

	if len(name) != 0 { Msg <- Messages{ ConSender:con, NameS:name, Text:fmt.Sprintf("\033[42m%s has joined our chat...\033[0m", name), Normal:true } }
	Handlemessage(con, name)
	delete(users.info, con)
}