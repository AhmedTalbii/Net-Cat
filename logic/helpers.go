package logic

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// format the current time
func UpdateTime() string {
    return time.Now().Format("2006-1-2 15:04:05")
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
	case !ContainASCIIchar(msg):
		return errors.New("out ascii")
	case len(msg) > 100:
		return errors.New("large_msg")
	}
	return nil
}

// this fucntion check if a string doesn't contain characters outside the ascii range
func ContainASCIIchar(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

// send formatted message to a connection with color
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
