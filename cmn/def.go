package cmn

import (
	"encoding/binary"
)

type Msg struct {
	Name string
	Text string
}

func (m *Msg)Marshal() []byte {
	buff := make([]byte, 2 + len(m.Name) + len(m.Text))
	binary.LittleEndian.PutUint16(buff, uint16(len(m.Name)))
	copy(buff[2:], m.Name+m.Text)
	return buff
}

func (m *Msg)UnMarshal(d []byte) error  {
	nameLen := binary.LittleEndian.Uint16(d)
	m.Name = string(d[2:2+nameLen])
	m.Text = string(d[2+nameLen:])
	return nil
}


