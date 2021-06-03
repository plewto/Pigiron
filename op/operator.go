// Package op defines the primary Pigiron structure, called an Operator.
//
// Operators are defined by the PigOp interface with common behavior
// implemented by the Operator struct.
//
// Each Operator has zero or more parents (inputs) and zero or more
// children (outputs).  Various types of Operators are linked together into
// a "MIDI Process Tree". Cyclical trees are not allowed.
//
// The Operator struct corresponds to an 'abstract class' and is not used
// directly.  Instead several structs extend Operator for specific
// behaviors.
//
// Operators should not be directly constructed.  Use The factory function
// MakeOperator instead.   
//

package op

import (
	"fmt"
	"errors"

	gomidi "gitlab.com/gomidi/midi"
	"github.com/plewto/pigiron/config"
)


// PigOp interface extends ChannelSelector and defines all methods on Operators.
//
// See Operator struct for concrete implementation.
//
type PigOp interface {
	ChannelSelector
	OperatorType() string
	Name() string
	setName(s string)
	commonInfo() string
	Info() string
	Panic()
	Reset()

	// Node
	IsRoot() bool
	IsLeaf() bool
	printTree(depth int)
	PrintTree()
	parents() map[string]PigOp
	Parents() map[string]PigOp
	children() map[string]PigOp
	Children() map[string]PigOp
	IsParentOf(child PigOp) bool
	IsChildOf(parent PigOp) bool
	circularTreeTest(depth int) bool
	Disconnect(child PigOp) PigOp
	Connect(child PigOp) error
	DisconnectAll()
	Disjoin()

	// OSC
	OSCAddress() string
	FormatOSCAddress(command string) string

	// MIDI
	MIDIEnabled() bool
	SetMIDIEnabled(flag bool)
	Accept(message gomidi.Message) bool
	distribute(message gomidi.Message)
	Send(message gomidi.Message)
	Close()
}


// Operator struct implements PigOp interface.
// Operator is not used directly, instead it serves as the super struct
// for all other Operator types.
//
type Operator struct {
	opType string
	name string
	channelSelector ChannelSelector
	parentMap map[string]PigOp
	childrenMap map[string]PigOp
	midiEnabled bool
}

// initOperator initialize Operator values.
// Extending structs should call initOperator as part of their construction
// process. 
func initOperator(op *Operator, opType string, name string, mode ChannelMode) {
	op.opType = opType
	op.name = name
	switch mode {
	case SingleChannel:
		op.channelSelector = MakeSingleChannelSelector()
	case MultiChannel:
		op.channelSelector = MakeMultiChannelSelector()
	default:
		op.channelSelector = MakeNullChannelSelector()
	}
	op.parentMap = make(map[string]PigOp)
	op.childrenMap = make(map[string]PigOp)
	op.midiEnabled = true
}

// OperatorType returns string representation for this specific Operator type.
//
func (op *Operator) OperatorType() string {
	return op.opType
}


// Name returns an identifying name for an Operator.
// An Operator's name must be unique.  If a new Operator's name clashes
// with an existing name, the registry will reassign it to be unique.
// Other then possible changes made to avoid a clash, an Operator's name
// is fixed for the duration of it's existence.
//
func (op *Operator) Name() string {
	return op.name
}


// setName changes an Operator's name.
// setName should be called no more the once for any specific Operator.
//
func (op *Operator) setName(s string) {
	op.name = s
}

func (op *Operator) String() string {
	return fmt.Sprintf("%s  name: %s  %s\n", op.opType, op.name, op.channelSelector)
}


// commonInfo returns a string detailing the current internal state of an Operator.
// The result is used by the Info method.
//
func (op *Operator) commonInfo() string {
	s := fmt.Sprintf("%s  name: %s    %s\n", op.opType, op.name, op.channelSelector)
	s += fmt.Sprintf("\tMIDI enabled: %v\n", op.MIDIEnabled())
	s += fmt.Sprintf("\tOSC address: '%s'\n", op.OSCAddress())
	s += "\tparents: "
	if op.IsRoot() {
		s += "<none>\n"
	} else {
		s += "\n"
		for name, _ := range op.parentMap {
			s += fmt.Sprintf("\t\t%s\n", name)
		}
	}
	s += "\tchildren: "
	if op.IsLeaf() {
		s += "<none>\n"
	} else {
		s += "\n"
		for name, _ := range op.childrenMap {
			s += fmt.Sprintf("\t\t%s\n", name)
		}
	}
	return s
}


// Info returns a string detailing the internal state of an Operator.
// Extending classes should append specific details the result of
// commonInfo.   By default Info simply returns the result of commonInfo.
//
func (op *Operator) Info() string {
	return op.commonInfo()
}


// Panic halts MIDI playback and kills all notes.
// The default behavior is to propagate a Panic to all child Operators
// without actually effecting MIDI.  Extending Operators should implement
// as needed.
//
func (op *Operator) Panic() {
	for _, child := range op.children {
		child.Panic()
	}
}

// Reset restores an Operator to an initial state.
// By default reset does nothing.  Extending Operators should
// implement as needed.
// 
func (op *Operator) Reset() {}


// IsRoot returns true for Operators without parents.
//
func (op *Operator) IsRoot() bool {
	return len(op.parentMap) == 0
}


// IsLeaf returns true for Operators without children.
//
func (op *Operator) IsLeaf() bool {
	return len(op.childrenMap) == 0
}

func (op *Operator) printTree(depth int) {
	switch {
	case depth > config.MaxTreeDepth:
		fmt.Printf("ERROR: MaxTreeDepth exceeded: %d\n", config.MaxTreeDepth)
		return 
	case depth == 0:
		fmt.Printf("%s\n", op.Name())
	default:
		pad := ""
		for i := 0; i <= depth; i++ {
			pad += "  "
		}
		fmt.Printf(" %s%s\n", pad, op.name)
	}
	for _, child := range op.childrenMap {
		child.printTree(depth + 1)
	}
}

// PrintTree prints a representation of the MIDI process tree.
//
func (op *Operator) PrintTree() {
	op.printTree(0)
}

func (op *Operator) parents() map[string]PigOp {
	return op.parentMap
}

func (op *Operator) children() map[string]PigOp {
	return op.childrenMap
}

// Parents returns a copy of an Operator's parent map.
//
func (op *Operator) Parents() map[string]PigOp {
	result := make(map[string]PigOp)
	for key, pop := range op.parents() {
		result[key] = pop
	}
	return result
}


// Children returns a copy of an Operator's children map.
//
func (op *Operator) Children() map[string]PigOp {
	result := make(map[string]PigOp)
	for key, pop := range op.children() {
		result[key] = pop
	}
	return result
}

// Disconnect removes the connection between an Operator and a child.
// It is not an error to Disconnect two non-connected Operators.
//
func (op *Operator) Disconnect(child PigOp) PigOp{
	delete(op.children(), child.Name())
	delete(child.parents(), op.Name())
	return child
}


// IsParentOf returns true if an Operator is a parent of child.
//
func (op *Operator) IsParentOf(child PigOp) bool {
	_, flag := op.children()[child.Name()]
	return flag
}

// IsChildOf returns true if an Operator is a child of parent.
//
func (op *Operator) IsChildOf(parent PigOp) bool {
	_, flag := op.parents()[parent.Name()]
	return flag
}


func (op *Operator) circularTreeTest(depth int) bool {
	if depth > config.MaxTreeDepth {
		return true
	} else {
		for _, c := range op.children() {
			return c.circularTreeTest(depth + 1)
		}
	}
	return false
}


// Connect makes a connection between an Operator and a child Operator.
// Returns an error if a circular-tree is produced as a result of the
// connection.
//
func (op *Operator) Connect(child PigOp) error {
	op.Disconnect(child)
	op.children()[child.Name()] = child
	child.parents()[op.Name()] = op
	var err error
	if op.circularTreeTest(0) {
		fstr := "Maximum tree depth exceeded at %s -> %s, MaxTreeDepth = %d"
		msg := fmt.Sprintf(fstr, op.Name(), child.Name(), config.MaxTreeDepth)
		err = errors.New(msg)
		op.Disconnect(child)
	}
	return err
}

// DisconnectAll removes the connection between an Operator and all of it's children.
//
func (op *Operator) DisconnectAll() {
	op.Panic()
	for _, child := range op.Children() {
		op.Disconnect(child)
	}
}

// Disjoin removes the connection between an Operator and all of it's parents.
//
func (op *Operator) Disjoin() {
	for _, parent := range op.Parents() {
		parent.Disconnect(op)
	}
}

// ChannelMode returns an Operator's ChannelSelector mode.
//
func (op *Operator) ChannelMode() ChannelMode {
	return op.channelSelector.ChannelMode()
}

// EnableChannel enables/disables given MIDI channel.
// Returns error if channel is invalid.
//
func (op *Operator) EnableChannel(channel int, flag bool) error {
	return op.channelSelector.EnableChannel(channel, flag)
}

// SelectChannel is the same as EnabledChannel(channel, true)
// Returns error if channel is invalid.
//
func (op *Operator) SelectChannel(channel int) error {
	return op.channelSelector.SelectChannel(channel)
}

// SelecrtedChannels returns a slice of currently enabled MIDI channels.
//
func (op *Operator) SelectedChannels() []int {
	return op.channelSelector.SelectedChannels()
}


// ChannelSelected returns true if channel is currently selected.
// Returns false for all out of bounds channel numbers.
//
func (op *Operator) ChannelSelected(channel int) bool {
	return op.channelSelector.ChannelSelected(channel)
}

// DeselectAllChannels disables all MIDI channels.
// For SingleChannel mode, channel 1 is selected.
//
func (op *Operator) DeselectAllChannels() {
	op.channelSelector.DeselectAllChannels()
}

// OSCAddress returns an Operator's base OSC address.
//
func (op *Operator) OSCAddress() string {
	sfmt := "/%s/op/%s/"
	return fmt.Sprintf(sfmt, config.ApplicationOSCPrefix, op.Name())
}


// FormatOSCAddress creates an Operator's OSC address for specific command.
//
func (op *Operator) FormatOSCAddress(command string) string {
	return fmt.Sprintf("%s%s", op.OSCAddress(), command)
}


// MIDIEnabled returns true if MIDI processing is enabled.
//
func (op *Operator) MIDIEnabled() bool {
	return op.midiEnabled
}

// SetMIDIEnabled disables/enables MIDI processing.
//
func (op *Operator) SetMIDIEnabled(flag bool) {
	op.midiEnabled = flag
}


// Accept returns true if a MIDI message should be processed.
// The Default behavior is to return true for all messages.
// Extending Operators should override for specific behavior.
//
func (op *Operator) Accept(msg gomidi.Message) bool {
	return true
}


// distribute ends a MIDI message to all child Operators.
//
func (op *Operator) distribute(msg gomidi.Message) {
	for _, child := range op.children() {
		child.Send(msg)
	}
}


// Send receives MIDI messages form parent Operators and transmits to all child Operators.
// Only messages for which Accept returns true are processed.
// The transmitted message need not be the same as the received message.
// Extending Operators should override.
//
func (op *Operator) Send(msg gomidi.Message) {
	if op.Accept(msg) && op.MIDIEnabled() {
		op.distribute(msg)
	}
}


// Close frees resources an Operator may be using.
// Most Operators simply ignore Close.  The MIDIInput and MIDIOutput
// Operators should close any open MIDI devices.
//
func (op *Operator) Close() {}
	
