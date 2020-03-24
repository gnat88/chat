package broker

import (
	"chat/cmn"
	"sync"
)


var ServerRoom *Room

const (
	roomChanSize = 1024
	historySize = 10
)

type User struct {
	Name string
	UpdatedAt int64
}
type Room struct {
	ID int
	Brokers sync.Map
	Chan chan []byte
	mu sync.Mutex
	histories []*cmn.Msg `json:"histories"`
}

func init() {
	ServerRoom = &Room{
		ID:    1,
		Brokers: sync.Map{},
		Chan:  make(chan []byte, roomChanSize),
		histories: make([]*cmn.Msg, 0),
	}
}

func (r *Room) GetLast() []*cmn.Msg{
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.histories
}

func (r *Room) PopLast(msg *cmn.Msg){
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.histories) > 10 {
		r.histories = append(r.histories[1:], msg)
	}else {
		r.histories = append(r.histories, msg)
	}
}


func (r *Room)Cast(msg *cmn.Msg) {
	bin := msg.Marshal()
	r.PopLast(msg)
	r.Chan <- bin
}

func (r *Room)Run() {
	for {
		select {
		 case msg := <- r.Chan:
		 	r.Brokers.Range(func(key, value interface{}) bool {
			    if ioWriter, ok := value.(*Broker); ok {
			    	ioWriter.Write(msg)
			    }
			    return true
		    })
		 	
		}
	}
}
