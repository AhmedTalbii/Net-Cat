package logic

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// this is a fucntion that scan the user massages and pass it to the (Msg) channel, and to the history file
func Handlemessage(conn net.Conn, name string) {
	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			Msg <- Messages{ ConSender:conn, NameS:name, Text:fmt.Sprintf("\033[41m%s has left our chat...\033[0m", name), Normal:true }
			return
		}
		text := scanner.Text()

		if text == "" {
			SendMessage(conn, "", "31", "Can't send Empty message\n")
			SendMessage(conn, "", "36", "["+UpdateTime()+"]"+"["+name+"]:")
			continue
		}

		if er := ValidMessage(text); er != nil {
			switch er.Error() {
			case "out ascii":
				SendMessage(conn, "", "31", "Can't send characters out of the ascii range between from (32 to 126)\n")
				SendMessage(conn, "", "36", "["+UpdateTime()+"]"+"["+name+"]:")
			case "large_msg":
				SendMessage(conn, "", "31", "message too large\n")
				SendMessage(conn, "", "36", "["+UpdateTime()+"]"+"["+name+"]:")
			}
			continue
		}

		f, err := os.OpenFile("assets/history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o777)
		if err != nil { return }

		m := fmt.Sprintf("\n\033[35m[%s][%s]:\033[0m%s", UpdateTime(), name, text)
		f.WriteString(m[1:] + "\n")
		f.Close()
		Msg <- Messages{ ConSender:conn, NameS:name, Text:text, Normal:false }
	}
}
