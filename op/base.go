package op

/*
** base.go defines baseOperator type which implements the Operator interface.
** All Operator types should extend baseOperator.
**
*/

import (
	"fmt"
	"errors"
	"github.com/rakyll/portmidi"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
)

/*
** baseOperator struct implements the Operator interface.
** All Operator classes should extend baseOperator.
*/
type baseOperator struct {
	opType string
	name string
	channelSelector midi.ChannelSelector
	parentMap map[string]Operator
	childrenMap map[string]Operator
	midiOutputEnabled bool
	dispatchTable map[string]func(*goosc.Message)([]string, error)
}

// op.initOperator() initializes the baseOperator
// This method should be call as part of the construction of all structs extending
// baseOperator.
//
func initOperator(op *baseOperator, opType string, name string, mode midi.ChannelMode) {
	op.opType = opType
	op.name = name
	switch mode {
	case midi.SingleChannel:
		op.channelSelector = midi.NewSingleChannelSelector()
	case midi.MultiChannel:
		op.channelSelector = midi.NewMultiChannelSelector()
	default:
		op.channelSelector = midi.NewNullChannelSelector()
	}
	op.parentMap = make(map[string]Operator)
	op.childrenMap = make(map[string]Operator)
	op.midiOutputEnabled = true
	op.dispatchTable = make(map[string]func(*goosc.Message)([]string, error))
	op.addCommandHandler("ping", op.remotePing)
	op.addCommandHandler("q-commands", op.remoteQueryCommands)
}


// op.OperatorType() returns name for this specific operator's type.
//
func (op *baseOperator) OperatorType() string {
	return op.opType
}

// op.Name() returns the unique operator's ID.
// No two Operator's should have the same name.
//
func (op *baseOperator) Name() string {
	return op.name
}

// op.setName() sets the operator's name.
// 
func (op *baseOperator) setName(s string) {
	op.name = s
}

func (op *baseOperator) String() string {
	return fmt.Sprintf("%-12s name: \"%s\"", op.opType, op.name)
}

// op.commonInfo() returns the base string for a bases operator's internal state.
// The Info() method for extending structures should use the result of commonInfo()
// for the bulk of their results.
//
func (op *baseOperator) commonInfo() string {
	s := fmt.Sprintf("%s  name: \"%s\"    %s\n", op.opType, op.name, op.channelSelector)
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

// op.Info() returns a string representation of the operator's internal state.
// By default Info() simple returns the result of commonInfo().  Extending
// operator's may add additional details.
//
func (op *baseOperator) Info() string {
	return op.commonInfo()
}

// op.Panic() should initiate a MIDI all-sounds off
// The action of Panic() is dependent on the specific operator type.
// Most operator types should simply pass the Panic to all of their child
// operators.   Play type operators should halt playback.   MIDIOutput type
// operators should immediately transmit all-note off messages.
//
func (op *baseOperator) Panic() {
	for _, child := range op.children() {
		child.Panic()
	}
}


// op.Reset() returns state of Operator to an initial state.
//
func (op *baseOperator) Reset() {
	op.SetMIDIOutputEnabled(true)
}

// op.Close() is used at Operator destruction time to release resources.
// It is mostly used to release MIDI devices on program shutdown.  An operator
// should not be used after calling it's Close method.
//
func (op *baseOperator) Close() {}


// op.IsRoot() returns true iff operator has no inputs.
//
func (op *baseOperator) IsRoot() bool {
	return len(op.parentMap) == 0
}

// op.IsLeaf() returns true iff operator has no outputs.
//
func (op *baseOperator) IsLeaf() bool {
	return len(op.childrenMap) == 0
}


func printTree(op Operator, depth int) {
	if int64(depth) > config.GlobalParameters.MaxTreeDepth {
		return
	}
	if depth == 0 {
		fmt.Printf("%s\n", op.Name())
	} else {
		pad := ""
		for i := 0; i < depth; i++ {
			pad += "   "
		}
		fmt.Printf("%s%s\n", pad, op.Name())
	}
	for _, child := range op.Children() {
		printTree(child, depth+1)
	}
}


// op.PrintTree() prints the structure of the MIDI process tree.
//
func (op *baseOperator) PrintTree() {
	printTree(op, 0)
}


// op.parents() returns a list of the operator's parents.
//
func (op *baseOperator) parents() map[string]Operator {
	return op.parentMap
}

// op.children() returns a list of the operator's children.
//
func (op *baseOperator) children() map[string]Operator {
	return op.childrenMap
}

// op.Parents() returns copy of the operator's parents list.
// It is safe to modify the contents of the list.
//
func (op *baseOperator) Parents() map[string]Operator {
	result := make(map[string]Operator)
	for key, pop := range op.parents() {
		result[key] = pop
	}
	return result
}

// op.Children() returns a copy of the operator's children list.
// It is safe to modify the contents of the list.
//
func (op *baseOperator) Children() map[string]Operator {
	result := make(map[string]Operator)
	for key, pop := range op.children() {
		result[key] = pop
	}
	return result
}

// op.Disconnect() disconnects the child from this operator.
// It is not an error if the child is not currently connected to the operator.
// Returns the child operator.
//
func (op *baseOperator) Disconnect(child Operator) Operator {
	delete(op.children(), child.Name())
	delete(child.parents(), op.Name())
	return child
}

// op.IsParentOf() returns true iff the operator is a parent of child.
//
func (op *baseOperator) IsParentOf(child Operator) bool {
	_, flag := op.children()[child.Name()]
	return flag
}

// op.IsChildOf() returns true iff operator is a child of parent.
//
func (op *baseOperator) IsChildOf(parent Operator) bool {
	_, flag := op.parents()[parent.Name()]
	return flag
}


// TODO: Validate
// op.circularTreeTest() returns true if the operator tree exceeds the maximum tree depth.
// The maximum tree depth is established during configuration.
// Trees exceeding this value are assumed to be circular.
//
func (op *baseOperator) circularTreeTest(depth int) bool {
	if depth > int(config.GlobalParameters.MaxTreeDepth) {
		return true
	} else {
		for _, c := range op.children() {
			return c.circularTreeTest(depth + 1)
		}
	}
	return false
}

// op.Connect() connects child as a child of operator.
// Duplicate connections are silently ignored.
//
// Returns non-nil error if the connection should cause a circular tree.
//
func (op *baseOperator) Connect(child Operator) error {
	op.Disconnect(child)
	op.children()[child.Name()] = child
	child.parents()[op.Name()] = op
	var err error
	if op.circularTreeTest(0) {
		fstr := "Maximum tree depth exceeded at %s -> %s, MaxTreeDepth = %d"
		msg := fmt.Sprintf(fstr, op.Name(), child.Name(), config.GlobalParameters.MaxTreeDepth)
		err = errors.New(msg)
		op.Disconnect(child)
	}
	return err
}


// op.DisconnectAll() disconnects all children from the operator.
//
func (op *baseOperator) DisconnectAll() {
	op.Panic()
	for _, child := range op.Children() {
		op.Disconnect(child)
	}
}

// op.DisconnectTree() recursively disconnects all decedents operators.
//
func (op *baseOperator) DisconnectTree() {
	op.Panic()
	children := op.Children()
	for _, child := range children {
		op.Disconnect(child)
		child.DisconnectTree()
	}
}

// op.DisconnectParents() removes all parents from the operator.
//
func (op *baseOperator) DisconnectParents() {
	for _, parent := range op.Parents() {
		parent.Disconnect(op)
	}
}


// op.ChannelMode() see midi.ChannelSelector interface.
//
func (op *baseOperator) ChannelMode() midi.ChannelMode {
	return op.channelSelector.ChannelMode()
}


// op.EnableChannel() see midi.ChannelSelector interface.
//
func (op *baseOperator) EnableChannel(channel midi.MIDIChannel, flag bool) error {
	return op.channelSelector.EnableChannel(channel, flag)
}


// op.SelectChannel() see midi.ChannelSelector interface.
//
func (op *baseOperator) SelectChannel(channel midi.MIDIChannel) error {
	return op.channelSelector.SelectChannel(channel)
}

// op.SelectedChannelIndexes() see midi.ChannelSelector interface.
//
func (op *baseOperator) SelectedChannelIndexes() []midi.MIDIChannelNibble {
	return op.channelSelector.SelectedChannelIndexes()
}


// op.ChannelIndexSelected() see midi.ChannelSelector interface.
//
func (op *baseOperator) ChannelIndexSelected(ci midi.MIDIChannelNibble ) bool {
	return op.channelSelector.ChannelIndexSelected(ci)
}

// op.DeselectAllChannels() see midi.ChannelSelector interface.
//
func (op *baseOperator) DeselectAllChannels() {
	op.channelSelector.DeselectAllChannels()
}

// op.SelectAllChannels() see midi.ChannelSelector interface.
//
func (op *baseOperator) SelectAllChannels() {
	op.channelSelector.SelectAllChannels()
}

// op.MIDIOutputEnabled() returns true if received MIDI messages are re-transmitted.
//
func (op *baseOperator) MIDIOutputEnabled() bool {
	return op.midiOutputEnabled
}

// op.SetMIDIOutputEnabled() enable/disable MIDI output.
//
func (op *baseOperator) SetMIDIOutputEnabled(flag bool) {
	op.midiOutputEnabled = flag
}

// op.Accept() returns true if the operator is to re-transmit the MIDI event.
// Extending classes should override as needed.
// The default always returns true.
//
func (op *baseOperator) Accept(event portmidi.Event) bool {
	return true
}

// op.distribute() transmits MIDI event to all child operators.
//
func (op *baseOperator) distribute(event portmidi.Event) {
	if op.MIDIOutputEnabled() {
		for _, child := range op.children() {
			child.Send(event)
		}
	}
}

// op.Send() receives MIDI events and selectively modifies and re-transmits.
// Extending classes should override as needed.
// The default is to re-transmit all events unchanged.
//
func (op *baseOperator) Send(event portmidi.Event) {
	if op.Accept(event) {
		op.distribute(event)
	}
}

// op.DispatchCommand() processes incoming OSC messages.
//
// Returns:
//  1. list (possibly empty) of results.
//  2. error
//
func (op *baseOperator) DispatchCommand(command string, msg *goosc.Message)([]string, error) {
	var err error
	var result []string
	handler, flag := op.dispatchTable[command]
	if !flag {
		msg := "Invalid command for %s operator %s.  command: '%s'"
		err = errors.New(fmt.Sprintf(msg, op.OperatorType(), op.Name(), command))
		return result, err
	}
	result, err = handler(msg)
	return result, err
}

// op.Commands() returns list of OSC command addresses.
//
func (op *baseOperator) Commands() []string {
	var keys []string = make([]string, 0, len(op.dispatchTable))
	for k := range op.dispatchTable {
		keys = append(keys, k)
	}
	return keys
}

// op.AddCommandHandler() adds a new handler function for a specific command string.
// Args:
//   command - simple string
//   handler - func(*go-osc.Message)([]string, error)
//
func (op *baseOperator) addCommandHandler(command string, handler func(*goosc.Message)([]string, error)) {
	op.dispatchTable[command] = handler
}
	

// op.remotePing() command handler for 'ping' command.
//     osc /pig/op name, ping
//     -> ACK
//
func (op *baseOperator) remotePing(*goosc.Message)([]string, error) {
	var err error
	var result = []string{op.Name(), "Ping"}
	return result, err
}

// op.remoteQueryCommands() command handler for 'q-commands' command.
//    osc /pig/op name, q-commands
//    -> ACK list of local commands.
//
func (op *baseOperator) remoteQueryCommands(*goosc.Message)([]string, error) {
	var err error
	return op.Commands(), err
}

