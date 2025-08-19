package logic

// this fucntion listens for messages from (Msg) channel and sends them to the all clients, it also format join and left user and chat text.
func StartListeningChan() {
	MsgRLU.RLock()
	for message := range Msg {
		for con, nameR := range users.info {
			if len(nameR) == 0 {
				continue
			}
			if con != message.ConSender {
				if message.Normal {
					SendMessage(con, "\n", "", message.Text)
				} else {
					SendMessage(con, "\n", "35", "[", UpdateTime(), "]", "[", message.NameS, "]:")
					SendMessage(con, "", "", message.Text)
				}
				SendMessage(con, "\n", "36", "[", UpdateTime(), "]", "[", nameR, "]:")
				continue
			}
			SendMessage(con, "", "36", "[", UpdateTime(), "][", message.NameS, "]:")
		}
	}
	MsgRLU.RUnlock()
}
