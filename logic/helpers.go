package helpers

import (
	"errors"
	"fmt"
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
// format the current time
func UpdateTime() string {
	return fmt.Sprint(time.Now().Year(), "-", int(time.Now().Month()), "-", time.Now().Day(), " ", time.Now().Hour(), ":", time.Now().Minute(), ":", time.Now().Second())
}

// this fucntion check if the name is contain only alphabetical character, otherwise it return error
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
