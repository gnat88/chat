package sess

import (
	"bufio"
	"chat/cmn"
	"chat/server/broker"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Session struct {
	Conn net.Conn
	Name string
	CloseSig chan struct{}
}


func (s *Session) Run() {
	defer s.Conn.Close()
	var err error
	header := make([]byte, 2)
	for {
		binary.LittleEndian.PutUint16(header, 0)
		_, err = io.ReadFull(s.Conn, header)
		if err != nil {
			panic(err)
		}
		msgLen := binary.LittleEndian.Uint16(header)
		buffer := make([]byte, msgLen)
		_, err = io.ReadFull(s.Conn, buffer)
		if err != nil{
			panic(err)
		}
		s.handler(buffer)
	}
}

func (s *Session) RunCmdHelper() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Printf("Type \"help\" for more information \n")
	fmt.Printf("$")
	for input.Scan() {
		line := input.Text()
		if line == "bye" {
			break
		}
		if line == "" {
			fmt.Printf("$")
			continue
		}
		line = strings.TrimSpace(line)
		s.cmd(line)
		fmt.Printf("\n$ ")
	}
}

func (s *Session) cmd(text string) {
	if text == "help" {
		fmt.Printf("Usage: /login [Name], /quit, help")
		return
	}
	
	if strings.HasPrefix(text, "/quit") {
		close(s.CloseSig)
		return
	}
	if strings.HasPrefix(text, "/login") {
		cmds := strings.Split(text, " ")
		if len(cmds) < 2 {
			fmt.Printf("login failed")
			return
		}
		s.Name = cmds[1]
		text = "/login " + cmds[1]
	}else {
		if s.Name == "" {
			fmt.Printf("not login")
			return
		}
	}
	
	msg := &cmn.Msg{
		Name: s.Name,
		Text: text,
	}
	
	bin := msg.Marshal()
	
	msgLen := len(bin)
	if msgLen > broker.MaxMessageLength {
		fmt.Printf("msg is too long")
		return
	}
	header := make([]byte, 2)
	binary.LittleEndian.PutUint16(header, uint16(msgLen))
	var err error
	_, err = s.Conn.Write(header)
	if err != nil {
		panic(err)
	}
	_, err = s.Conn.Write(bin)
	if err != nil {
		panic(err)
	}
}

func (s *Session) handler(d []byte) {
	msg := &cmn.Msg{}
	msg.UnMarshal(d)
	fmt.Printf("\n%s: %s\n$", msg.Name, msg.Text)
}
