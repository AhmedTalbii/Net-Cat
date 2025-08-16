package helpers

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var (
	Con []Client
	Msg chan Message
)

type Client struct {
	Name string
	Con  net.Conn
}

type Message struct {
	IpSender, Message string
}

func HandleMessage(conn net.Conn, name string) {
	scanner := bufio.NewScanner(conn)
	var f bool = false

	for {
		time := UpdateTime()
		if !f {
			fmt.Fprint(conn, "["+time+"]"+"["+name+"]:")
			f = true
		}
		if !scanner.Scan() {
			fmt.Fprint(conn, "Error happend in scanning name")
			return
		}
		text := scanner.Text()
		for _, con := range Con {
			updateTime := UpdateTime()
			if con.Con != conn {
				fmt.Fprint(con.Con, "\n["+updateTime+"]"+"["+name+"]:"+text)
				fmt.Fprint(con.Con, "\n["+updateTime+"]"+"["+con.Name+"]:")

			} else {
				fmt.Fprint(conn, "["+updateTime+"]"+"["+name+"]:")
			}
		}
	}
}

func UpdateTime() string {
	return fmt.Sprint(time.Now().Year(), "-", int(time.Now().Month()), "-", time.Now().Day(), " ", time.Now().Hour(), ":", time.Now().Minute(), ":", time.Now().Second())
}
