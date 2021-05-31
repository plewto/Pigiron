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
// Operators should not be directly constructed.  Use The factory  function
// MakeOperator instead.   In addition to creating new operators
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
// OperatorType() returns type of operator as a string.
//
// Name() returns the Operator's name.
//    Each Operator has a unique name as assigned when it is added to the
//    registry.  Once assigned by the registry, an Operator's name may not
//    be changed. 
//
// setName() set's the operator's name.
//
// Info() returns strings indicating internal state for the Operator.
//
// Panic() transmits a MIDI panic.
//     The behavior of a Panic is dependent on the specific Operator type.
//     The default panic does nothing.
//
// Reset() sets Operator parameter's to an initial state.
//     The default reset does nothing.
//
// IsRoot() returns true if the operator has no parents.
//
// IsLeaf() returns true if the operator has no children.
//
// PrintTree() prints the structure of the MIDI process tree.
//    Typically PrintTree is called on a root operator.
//
// parents() returns the parents map.
//
// Parents() returns a copy of the parents map.
//
// children() returns the children map.
//
// Children() returns a copy of the children map.
//
// IsParentOf(child) returns true if operator is a parent of child.
//
// IsChildOf(parent) returns true if operator is a child of parent.
//
// Disconnect(child) removes child from operators children list.
//   It is not an error if the two operators are not currently connected. 
//   Returns the child operator.
//
// Connect(child) adds child to the operator.
//   If as a result of connecting child to parent the maximum tree depth is
//   exceeded, the two operators are disconnected and an error is returned.
//
// DisconnectAll() disconnects all children of the operator.
//
// Disjoin() disconnects the operator from all of it's parents.
//
// OSCAddress() returns the operator's OSC address prefix.
//
// FormatOSCAddress(command) returns the OSC address for the given command.
//
// MIDIEnabled() returns true if MIDI input is enabled.
//
// SetMIDIEnabled(flag) Enable/disable MIDI input.
//
type PigOp interface {
	ChannelSelector
	OperatorType() string
	Name() string
	setName(s string)
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
	
}


// Operator struct implements PigOp interface.
// The Operator struct is not used directly, instead it serves as the super
// struct for all other Operator types.
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


func (op *Operator) OperatorType() string {
	return op.opType
}

func (op *Operator) Name() string {
	return op.name
}

func (op *Operator) setName(s string) {
	op.name = s
}

func (op *Operator) String() string {
	return fmt.Sprintf("%s  name: %s  %s\n", op.opType, op.name, op.channelSelector)
}

func (op *Operator) Info() string {
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

func (op *Operator) Panic() {}

func (op *Operator) Reset() {}

func (op *Operator) IsRoot() bool {
	return len(op.parentMap) == 0
}

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

func (op *Operator) PrintTree() {
	op.printTree(0)
}

func (op *Operator) parents() map[string]PigOp {
	return op.parentMap
}

func (op *Operator) children() map[string]PigOp {
	return op.childrenMap
}

func (op *Operator) Parents() map[string]PigOp {
	result := make(map[string]PigOp)
	for key, pop := range op.parents() {
		result[key] = pop
	}
	return result
}

func (op *Operator) Children() map[string]PigOp {
	result := make(map[string]PigOp)
	for key, pop := range op.children() {
		result[key] = pop
	}
	return result
}

func (op *Operator) Disconnect(child PigOp) PigOp{
	delete(op.children(), child.Name())
	delete(child.parents(), op.Name())
	return child
}

func (op *Operator) IsParentOf(child PigOp) bool {
	_, flag := op.children()[child.Name()]
	return flag
}

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


func (op *Operator) ChannelMode() ChannelMode {
	return op.channelSelector.ChannelMode()
}

func (op *Operator) EnableChannel(channel int, flag bool) error {
	return op.channelSelector.EnableChannel(channel, flag)
}

func (op *Operator) SelectChannel(channel int) error {
	return op.channelSelector.SelectChannel(channel)
}


func (op *Operator) SelectedChannels() []int {
	return op.channelSelector.SelectedChannels()
}

func (op *Operator) ChannelSelected(channel int) bool {
	return op.channelSelector.ChannelSelected(channel)
}

func (op *Operator) DeselectAllChannels() {
	op.channelSelector.DeselectAllChannels()
}

func (op *Operator) DisconnectAll() {
	op.Panic()
	for _, child := range op.Children() {
		op.Disconnect(child)
	}
}


func (op *Operator) Disjoin() {
	for _, parent := range op.Parents() {
		parent.Disconnect(op)
	}
}

func (op *Operator) OSCAddress() string {
	sfmt := "%s/op/%s/"
	return fmt.Sprintf(sfmt, config.ApplicationOSCPrefix, op.Name())
}

func (op *Operator) FormatOSCAddress(command string) string {
	return fmt.Sprintf("%s%s", op.OSCAddress(), command)
}


func (op *Operator) MIDIEnabled() bool {
	return op.midiEnabled
}

func (op *Operator) SetMIDIEnabled(flag bool) {
	op.midiEnabled = flag
}
	
func (op *Operator) Accept(msg gomidi.Message) bool {
	return true
}

func (op *Operator) distribute(msg gomidi.Message) {
	for _, child := range op.children() {
		child.Send(msg)
	}
}


func (op *Operator) Send(msg gomidi.Message) {
	if op.Accept(msg) && op.MIDIEnabled() {
		op.distribute(msg)
	}
}
