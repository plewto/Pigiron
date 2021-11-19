package smf

import (
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
)

// Event struct combines MIDI message with delta time.
//
type Event struct {
	deltaTime uint64
	message gomidi.Message
}

func (ev *Event) String() string {
	return fmt.Sprintf("Δt %8d : %s", ev.deltaTime, ev.message)
}

// Event.Length() method returns byte count of event message.
//
func (ev *Event) Length() int {
	return len(ev.message.Data)
}

func (ev *Event) Dump() {
	fmt.Println(ev.String())
}

func (ev *Event) DeltaTime() uint64 {
	return ev.deltaTime
}

func (ev *Event) Message() gomidi.Message {
	return ev.message
}

