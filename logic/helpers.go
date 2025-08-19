package helpers

import (
	"errors"
	"fmt"
	"net"
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
				if !message.Normal {
					SendMessage(con, "\n", "35", "[", UpdateTime(), "]", "[", message.NameS, "]:")
					SendMessage(con, "", "", message.Text)
				} else {
					if strings.HasSuffix(message.Text, "has joined our chat...") {
						SendMessage(con, "\n", "42", message.Text)
						// m := fmt.Sprintf("\n\033[42m" + message.Text + "\033[0m")
						// fmt.Fprint(con, m)
					} else if strings.HasSuffix(message.Text, "has left our chat...") && message.NameS == "" {
						continue
					} else if strings.HasSuffix(message.Text, "has left our chat...") {
						SendMessage(con, "\n", "41", message.Text)
						// fmt.Fprint(con, "\n\033[41m"+message.Text+"\033[0m")
					} else {
						SendMessage(con, "\n", "35", "[", UpdateTime(), "]", "[", message.NameS, "]:")
						SendMessage(con, "", "", message.Text)
						// m := fmt.Sprintf("\n\033[35m[%s][%s]:\033[0m%s", UpdateTime(), message.NameS, message.Text)
						// fmt.Fprint(con, m)
					}
				}
				// fmt.Fprint(con, "\n\033[36m["+UpdateTime()+"]"+"["+nameR+"]:\033[0m")
				SendMessage(con, "\n", "36", "[", UpdateTime(), "]", "[", nameR, "]:")
				continue
			}
			SendMessage(con, "", "36", "[", UpdateTime(), "][", message.NameS, "]:")
			// fmt.Fprint(con, "\033[36m["+UpdateTime()+"]"+"["+message.NameS+"]:\033[0m")
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
		return errors.New("\033[1;31minvalid name length (1-15), try again:\033[0m\n[ENTER YOUR NAME]: ")
	}

	for _, char := range name {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') {
			return errors.New("\033[1;31minvalid name, only letters allowed, try again:\033[0m\n[ENTER YOUR NAME]: ")
		}
	}

	users.RLock()
	defer users.RUnlock()
	for _, v := range users.info {
		if strings.EqualFold(v, name) {
			return errors.New("\033[1;31mname already exists, try again:\033[0m\n[ENTER YOUR NAME]: ")
		}
	}
	return nil
}

func ValidMessage(msg string) error {
	switch {
	case strings.HasSuffix(msg, "has joined our chat...") || strings.HasSuffix(msg, "has left our chat..."):
		return errors.New("The_user_send_joined_or_left")
	case !ContainASCIIchar(msg):
		return errors.New("out ascii")
	case len(msg) > 100:
		return errors.New("large_msg")
	}

	return nil
}

func ContainASCIIchar(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func SendMessage(con net.Conn, Nl, color string, parts ...string) {
	if color == "" {
		color = "0"
	}
	fmt.Fprint(con, Nl)
	fmt.Fprint(con, "\033["+color+"m")
	for _, part := range parts {
		fmt.Fprint(con, part)
	}
	fmt.Fprint(con, "\033[0m")
}
