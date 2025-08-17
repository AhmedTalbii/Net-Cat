package helpers

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// this fucntion listens for messages from (Msg) channel and sends them to the all clients, it also format join and left user and chat text.
func StartListeningChan() {
	MsgRLU.RLock()
	for message := range Msg {
		for con, nameR := range users.info {
			if len(nameR) == 0 {
				continue
			}
			f, err := os.OpenFile("assets/history.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o777)
			if err != nil {
				return
			}

			if con != message.ConSender {
				if strings.HasSuffix(message.Text, "has joined our chat...") {
					m := fmt.Sprintf("\n\033[42m" + message.Text + "\033[0m")

					fmt.Fprint(con, m)
				} else if strings.HasSuffix(message.Text, "has left our chat...") && message.NameS == "" {
					continue
				} else if strings.HasSuffix(message.Text, "has left our chat...") {
					fmt.Fprint(con, "\n\033[41m"+message.Text+"\033[0m")
				} else {
					m := fmt.Sprintf("\n\033[35m[%s][%s]:\033[0m%s", UpdateTime(), message.NameS, message.Text)
					f.WriteString(m[1:] + "\n")
					f.Close()
					fmt.Fprint(con, m)
				}

				fmt.Fprint(con, "\n\033[36m["+UpdateTime()+"]"+"["+nameR+"]:\033[0m")
				continue
			}
			fmt.Fprint(con, "\033[36m"+UpdateTime()+"]"+"["+message.NameS+"]:\033[0m")
		}
	}
	MsgRLU.RUnlock()
}

// format the current time
func UpdateTime() string {
	return fmt.Sprint(time.Now().Year(), "-", int(time.Now().Month()), "-", time.Now().Day(), " ", time.Now().Hour(), ":", time.Now().Minute(), ":", time.Now().Second())
}

// this fucntion check if the name is contain only alphabetical character, otherwise it return error
func Valid_Name(name string) error {
	if name == "" || len(name) > 15 {
		return errors.New("\033[1;31m invalid name length (1-15), try again:\n[ENTER YOUR NAME]: \033[0m")
	}

	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return errors.New("\033[1;31m invalid name, only letters allowed, try again:\n[ENTER YOUR NAME]: \033[0m")
		}
	}

	users.RLock()
	defer users.RUnlock()
	for _, v := range users.info {
		if strings.EqualFold(v, name) {
			return errors.New("\033[1;31m name already exists, try again:\n[ENTER YOUR NAME]: \033[0m")
		}
	}
	return nil
}
