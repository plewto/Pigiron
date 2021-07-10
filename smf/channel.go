package smf

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

// Implements RealtimeMIDIMessage
//
type ChannelMessage struct {
	bytes []byte
}

func NewChannelMessage(status StatusByte, channelByte byte, data1 byte, data2 byte) (*ChannelMessage, error) {
	var err error
	var m *ChannelMessage
	if !isChannelStatus(byte(status)) {
		msg := "ERROR NewChannelMessage, illegal status byte 0x%x  '%s'"
		err = fmt.Errorf(msg, byte(status), status)
		return m, err
	}
	count, _ := channelStatusByteCount[status]
	bytes := make([]byte, count)
	bytes[0] = byte(status) | channelByte
	bytes[1] = data1
	if count > 2 {
		bytes[2] = data2
	}
	m = &ChannelMessage{bytes}
	return m, err
}

func (m *ChannelMessage) Status() StatusByte {
	sb := m.bytes[0]
	return StatusByte(sb & 0xF0)
}

func (m *ChannelMessage) ChannelByte() byte {
	return m.bytes[0] & 0x0F
}

func (m *ChannelMessage) Bytes() []byte {
	return m.bytes
}

func (m *ChannelMessage) String() string {
	s, _ := statusTable[m.Status()]
	c := m.ChannelByte() + 1
	acc := fmt.Sprintf("%s  chan: %2d  data: ", s, c)
	for _, b := range m.bytes[1:] {
		acc += fmt.Sprintf("%3d ", b)
	}
	return acc
}

func (m *ChannelMessage) Dump() {
	fmt.Print("[")
	for _, b := range m.bytes {
		fmt.Printf("%02x ", b)
	}
	fmt.Printf("] %s\n", m)
}

func (m *ChannelMessage) ToPortmidiEvent() portmidi.Event {
	var time portmidi.Timestamp = portmidi.Timestamp(0)
	var status int64 = int64(m.bytes[0])
	var d1, d2 int64
	d1 = int64(m.bytes[1])
	if len(m.bytes) > 2 {
		d2 = int64(m.bytes[2])
	} else {
		d2 = 0
	}
	var sysex = make([]byte, 0)
	pme := portmidi.Event{time, status, d1, d2, sysex}
	return pme
}
	
		
	

		
