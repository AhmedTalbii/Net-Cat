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
	defer listner.Close()
	var err error

	file, err = os.Create("assets/history.txt")
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
		if len(users.info) >= 2 {
			fmt.Fprint(conn, "the chat room is currently full, please try again later\n")
			conn.Close()
			users.Unlock()
			continue
		}
		users.info[conn] = ""
		users.Unlock()

		// this is a function that write the pingWing Message to the user asking him for entering his name, then it scan the name and pass it to the map and to the channel, then it read the history data if exist
		handleCLient := func(con net.Conn) {
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
				"     `-'       `--'\n" +
				"[ENTER YOUR NAME]:"
			con.Write([]byte(Ping_Win_Mess))

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
				conn.Write(history)
			}

			if len(name) != 0 {
				Msg <- Messages{
					ConSender: con,
					NameS:     name,
					Text:      fmt.Sprintf("%s has joined our chat...", name),
				}
			}
			// this is a fucntion that scan the user massages and pass it to the (Msg) channel, and to the history file
			handlemessage := func(conn net.Conn, name string) {
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
					if text == "" {
						fmt.Fprintf(conn, "["+UpdateTime()+"]"+"["+name+"]: ")
						continue
					}

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
			handlemessage(con, name)
			delete(users.info, con)
		}
		go handleCLient(conn)

	}
}
