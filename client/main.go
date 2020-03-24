package main

import (
	sess2 "chat/client/sess"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:18000")
	if err != nil {
		panic(err)
	}
	sess := &sess2.Session{Conn: conn, CloseSig:make(chan struct{})}
	go sess.RunCmdHelper()
	go sess.Run()
	<- sess.CloseSig
}

