package smf

import (
	"fmt"
	"testing"
	"github.com/rakyll/portmidi"
)

func TestChannelMessage(t *testing.T) {
	fmt.Print()
	var err error
	var cmsg *ChannelMessage

	_, err = NewChannelMessage(SysexStatus, 0, 0, 0)
	if err == nil {
		t.Fatal("NewChannelMessage did not detect invalid status byte")
	}

	cmsg, err = NewChannelMessage(NoteOffStatus, 1, 60, 64)
	if err != nil {
		t.Fatal("NewChannelMessage returned incorrect error")
	}
	if cmsg.Status() != NoteOffStatus {
		t.Fatal("ChannelMessage.Status() did not return expected value.")
	}
	if cmsg.ChannelByte() != 1 {
		t.Fatal("ChannelMessage.ChannelByte() did not return expected value.")
	}
	bytes := cmsg.Bytes()
	if len(bytes) != 3 {
		t.Fatalf("ChannelMessage.Bytes() length incorrect, expected 3, got %d", len(bytes))
	}
	if bytes[1] != 60 || bytes[2] != 64 {
		msg := "ChannelMessage bytes contains incorrect values, expected [x 60, 64], got %v"
		t.Fatalf(msg, bytes)
	}
	var pmEvent portmidi.Event = cmsg.ToPortmidiEvent()
	if pmEvent.Status != int64(bytes[0]) {
		msg := "ToPortmidiEvent conversion has wrong status, expected %d, got %d"
		t.Fatalf(msg, bytes[0], pmEvent.Status)
	}
	if pmEvent.Data1 != int64(bytes[1]) {
		msg := "ToPortmidiEvent conversion has wrong data1, expected %d, got %d"
		t.Fatalf(msg, bytes[1], pmEvent.Data1)
	}
	if pmEvent.Data2 != int64(bytes[2]) {
		msg := "ToPortmidiEvent conversion has wrong data2, expected %d, got %d"
		t.Fatalf(msg, bytes[2], pmEvent.Data2)
	}
}
