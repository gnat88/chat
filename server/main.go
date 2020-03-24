package main

import (
	"chat/server/broker"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":18000")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("server started at port 18000")
	broker.LoadWorlds("list.txt")
	go broker.ServerRoom.Run()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("new connection ")
		sess := broker.NewBroker(conn)
		go sess.Run()
	}
}
