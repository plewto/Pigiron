package op

import (
	"fmt"
	"errors"

	midi "github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
)


type Operator interface {
	midi.ChannelSelector
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
	parents() map[string]Operator
	Parents() map[string]Operator
	children() map[string]Operator
	Children() map[string]Operator
	IsParentOf(child Operator) bool
	IsChildOf(parent Operator) bool
	circularTreeTest(depth int) bool
	Disconnect(child Operator) Operator
	Connect(child Operator) error
	DisconnectAll()
	Disjoin()

	// OSC
	OSCAddress() string
	FormatOSCAddress(command string) string

	// MIDI
	MIDIEnabled() bool
	SetMIDIEnabled(flag bool)
	// Accept(message gomidi.Message) bool
	// distribute(message gomidi.Message)
	// Send(message gomidi.Message)
}


type baseOperator struct {
	opType string
	name string
	channelSelector midi.ChannelSelector
	parentMap map[string]Operator
	childrenMap map[string]Operator
	midiEnabled bool
}

func initOperator(op *baseOperator, opType string, name string, mode midi.ChannelMode) {
	op.opType = opType
	op.name = name
	switch mode {
	case midi.SingleChannel:
		op.channelSelector = midi.MakeSingleChannelSelector()
	case midi.MultiChannel:
		op.channelSelector = midi.MakeMultiChannelSelector()
	default:
		op.channelSelector = midi.MakeNullChannelSelector()
	}
	op.parentMap = make(map[string]Operator)
	op.childrenMap = make(map[string]Operator)
	op.midiEnabled = true
}

func (op *baseOperator) OperatorType() string {
	return op.opType
}


func (op *baseOperator) Name() string {
	return op.name
}


func (op *baseOperator) setName(s string) {
	op.name = s
}

func (op *baseOperator) String() string {
	return fmt.Sprintf("%s  \"%s\"\n", op.opType, op.name)
}


func (op *baseOperator) commonInfo() string {
	s := fmt.Sprintf("%s  name: \"%s\"    %s\n", op.opType, op.name, op.channelSelector)
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


func (op *baseOperator) Info() string {
	return op.commonInfo()
}


func (op *baseOperator) Panic() {
	for _, child := range op.children() {
		child.Panic()
	}
}

func (op *baseOperator) Reset() {}


func (op *baseOperator) IsRoot() bool {
	return len(op.parentMap) == 0
}


func (op *baseOperator) IsLeaf() bool {
	return len(op.childrenMap) == 0
}

func (op *baseOperator) printTree(depth int) {
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

func (op *baseOperator) PrintTree() {
	op.printTree(0)
}

func (op *baseOperator) parents() map[string]Operator {
	return op.parentMap
}

func (op *baseOperator) children() map[string]Operator {
	return op.childrenMap
}

func (op *baseOperator) Parents() map[string]Operator {
	result := make(map[string]Operator)
	for key, pop := range op.parents() {
		result[key] = pop
	}
	return result
}


func (op *baseOperator) Children() map[string]Operator {
	result := make(map[string]Operator)
	for key, pop := range op.children() {
		result[key] = pop
	}
	return result
}

func (op *baseOperator) Disconnect(child Operator) Operator{
	delete(op.children(), child.Name())
	delete(child.parents(), op.Name())
	return child
}


func (op *baseOperator) IsParentOf(child Operator) bool {
	_, flag := op.children()[child.Name()]
	return flag
}

func (op *baseOperator) IsChildOf(parent Operator) bool {
	_, flag := op.parents()[parent.Name()]
	return flag
}


func (op *baseOperator) circularTreeTest(depth int) bool {
	if depth > config.MaxTreeDepth {
		return true
	} else {
		for _, c := range op.children() {
			return c.circularTreeTest(depth + 1)
		}
	}
	return false
}


func (op *baseOperator) Connect(child Operator) error {
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

func (op *baseOperator) DisconnectAll() {
	op.Panic()
	for _, child := range op.Children() {
		op.Disconnect(child)
	}
}

func (op *baseOperator) Disjoin() {
	for _, parent := range op.Parents() {
		parent.Disconnect(op)
	}
}

func (op *baseOperator) ChannelMode() midi.ChannelMode {
	return op.channelSelector.ChannelMode()
}

func (op *baseOperator) EnableChannel(channel int, flag bool) error {
	return op.channelSelector.EnableChannel(channel, flag)
}

func (op *baseOperator) SelectChannel(channel int) error {
	return op.channelSelector.SelectChannel(channel)
}

func (op *baseOperator) SelectedChannels() []int {
	return op.channelSelector.SelectedChannels()
}


func (op *baseOperator) ChannelSelected(channel int) bool {
	return op.channelSelector.ChannelSelected(channel)
}

func (op *baseOperator) DeselectAllChannels() {
	op.channelSelector.DeselectAllChannels()
}

func (op *baseOperator) OSCAddress() string {
	sfmt := "/%s/op/%s/"
	return fmt.Sprintf(sfmt, config.ApplicationOSCPrefix, op.Name())
}


func (op *baseOperator) FormatOSCAddress(command string) string {
	return fmt.Sprintf("%s%s", op.OSCAddress(), command)
}


func (op *baseOperator) MIDIEnabled() bool {
	return op.midiEnabled
}

func (op *baseOperator) SetMIDIEnabled(flag bool) {
	op.midiEnabled = flag
}


/* ******
// The Default behavior is to return true for all messages.
// Extending Operators should override for specific behavior.
//
func (op *baseOperator) Accept(msg gomidi.Message) bool {
	return true
}


// distribute ends a MIDI message to all child Operators.
//
func (op *baseOperator) distribute(msg gomidi.Message) {
	for _, child := range op.children() {
		child.Send(msg)
	}
}


// Send receives MIDI messages form parent Operators and transmits to all child Operators.
// Only messages for which Accept returns true are processed.
// The transmitted message need not be the same as the received message.
// Extending Operators should override.
//
func (op *baseOperator) Send(msg gomidi.Message) {
	if op.Accept(msg) && op.MIDIEnabled() {
		op.distribute(msg)
	}
}
***** */


