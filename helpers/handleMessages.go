package helpers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func StartListeningChan() {
	MsgRLU.RLock()
	for message := range Msg {
		for con, nameR := range users.info {
			updateTime := UpdateTime()
			if len(nameR) == 0 {
				continue
			}
			if con != message.ConSender {
				if strings.HasSuffix(message.Text, "has joined our chat...") {
					fmt.Fprint(con, "\n"+message.Text)
				} else {
					fmt.Fprint(con, "\n["+updateTime+"]"+"["+message.NameS+"]:"+message.Text)
				}
				fmt.Fprint(con, "\n["+updateTime+"]"+"["+nameR+"]:")
				continue
			}
			fmt.Fprint(con, "["+updateTime+"]"+"["+message.NameS+"]:")
		}
	}
	MsgRLU.RUnlock()
}

func HandleMessage(conn net.Conn, name string) {
	scanner := bufio.NewScanner(conn)

	for {
		if !scanner.Scan() {
			fmt.Fprint(conn, "Error happend in scanning name")
			return
		}
		text := scanner.Text()

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
