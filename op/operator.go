package op

import (
	"github.com/rakyll/portmidi"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)

// Operator interface defines the primary pigiron object.
// Operators are linked together as nodes on a MIDI process tree.
// Each Operator has zero or more parents (inputs) and zero or more
// children (outputs).  When an Operator receives a MIDI message via it's
// Send method, it selectively forwards it to all of it's child Operators.
// The message may be modified prior to rebroadcasting it.
// 

type Operator interface {
	midi.ChannelSelector
	OperatorType() string
	Name() string
	setName(s string)
	commonInfo() string
	Info() string
	Panic()
	Reset()
	Close()

	// Node
	IsRoot() bool
	IsLeaf() bool
	printTree(depth int)
	PrintTree()
	parents() map[string]Operator
	Parents() map[string]Operator
	children() map[string]Operator
	Children() map[string]Operator
	IsParentOf(child Operator) bool
	IsChildOf(parent Operator) bool
	circularTreeTest(depth int) bool
	Connect(child Operator) error
	Disconnect(child Operator) Operator
	DisconnectAll()
	DisconnectTree()
	DisconnectParents()

	// OSC
	DispatchCommand(command string, msg *goosc.Message)([]string, error)
	Commands() []string
	addCommandHandler(command string, handler func(*goosc.Message)([]string, error))
	
	// MIDI
	MIDIOutputEnabled() bool
	SetMIDIOutputEnabled(flag bool)
	
	Accept(event portmidi.Event) bool
	distribute(event portmidi.Event)
	Send(event portmidi.Event)
}
