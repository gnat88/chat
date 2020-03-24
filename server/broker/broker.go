package broker

import (
	"chat/cmn"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	HeaderLength = 2
	MaxMessageLength = 65535
)
var (
	ErrMsgHeadParse = errors.New("read message header error")
)

type Broker struct {
	conn         net.Conn
	roomID int
}

func NewBroker(conn net.Conn) *Broker {
	p := new(Broker)
	p.conn = conn
	return p
}

func (p *Broker) Run() {
	for {
		header := make([]byte, HeaderLength)
		if _, err := io.ReadFull(p.conn, header); err != nil {
			log.Printf("%v", err)
			return
		}
		
		msgLength := binary.LittleEndian.Uint16(header)
		if msgLength > MaxMessageLength || msgLength <= 0 {
			log.Printf("%v", fmt.Errorf("msg length invalid"))
			return
		}
		bin := make([]byte, msgLength)
		if _, err := io.ReadFull(p.conn, bin); err != nil {
			log.Printf("%v", err)
			return
		}
		log.Printf("%v\n", string(bin))
		msg := &cmn.Msg{}
		msg.UnMarshal(bin)
		p.handler(msg)
		log.Printf("%+v", msg)
		
	}
}

func (p *Broker) Write(args []byte) error {
	msgLen := uint16(len(args))
	header := make([]byte, HeaderLength)
	binary.LittleEndian.PutUint16(header, msgLen)
	
	var err error
	_, err = p.conn.Write(header)
	if err != nil {
		return err
	}
	
	_, err = p.conn.Write(args)
	if err != nil {
		return err
	}
	
	return nil
}

func (p *Broker) Close() {
	p.conn.Close()
}

func (p *Broker) handler(msg *cmn.Msg) {
	if msg.Name == "" {
		log.Printf("name is empty")
		return
	}
	msg.Text = Replace(msg.Text)
	if strings.HasPrefix(msg.Text, "/login") {
		ServerRoom.Brokers.Store(msg.Name, p)
		lastTen := ServerRoom.GetLast()
		for _, v := range lastTen {
			p.Write(v.Marshal())
		}
		return
	}
	ServerRoom.Cast(msg)
}
