package helpers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func StartListeningChan() {
	MsgRLU.RLock()
	for message := range Msg {
		for con, nameR := range users.info {
			// updateTime := UpdateTime()
			if len(nameR) == 0 {
				continue
			}
			if con != message.ConSender {
				if strings.HasSuffix(message.Text, "has joined our chat...") {
					fmt.Fprint(con, "\n"+message.Text)
				} else if strings.HasSuffix(message.Text, "has left our chat...") && message.NameS == "" {
					continue
				} else if strings.HasSuffix(message.Text, "has left our chat...") {
					fmt.Fprint(con, "\n"+message.Text)
				} else {
					fmt.Fprint(con, "\n["+UpdateTime()+"]"+"["+message.NameS+"]:"+message.Text)
				}
				fmt.Fprint(con, "\n["+UpdateTime()+"]"+"["+nameR+"]:")
				continue
			}
			fmt.Fprint(con, "["+UpdateTime()+"]"+"["+message.NameS+"]:")
		}
	}
	MsgRLU.RUnlock()
}

func HandleMessage(conn net.Conn, name string) {
	scanner := bufio.NewScanner(conn)

	for {
		if !scanner.Scan() {
			fmt.Fprint(conn, "Error happend in scanning name")
			Msg <- Messages{
				ConSender: conn,
				NameS:     name,
				Text:      fmt.Sprintf("%s has left our chat...", name),
			}
			return
		}
		text := scanner.Text()
		f, err := os.OpenFile("assets/history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o777)
		if err != nil {
			return
		}
		f.WriteString("[" + UpdateTime() + "]" + "[" + name + "]: " + text + "\n")
		f.Close()

		Msg <- Messages{
			ConSender: conn,
			NameS:     name,
			Text:      text,
		}
	}
}

func UpdateTime() string {
	return fmt.Sprint(time.Now().Year(), "-", int(time.Now().Month()), "-", time.Now().Day(), " ", time.Now().Hour(), ":", time.Now().Minute(), ":", time.Now().Second())
}
