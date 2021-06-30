package op

import (
	"fmt"
	"errors"

	"github.com/rakyll/portmidi"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
	// "github.com/plewto/pigiron/osc"
)

type baseOperator struct {
	opType string
	name string
	channelSelector midi.ChannelSelector
	parentMap map[string]Operator
	childrenMap map[string]Operator
	midiOutputEnabled bool
	dispatchTable map[string]func(*goosc.Message)([]string, error)
}

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

func (op *baseOperator) Reset() {
	op.SetMIDIOutputEnabled(true)
}

func (op *baseOperator) Close() {}

func (op *baseOperator) IsRoot() bool {
	return len(op.parentMap) == 0
}


func (op *baseOperator) IsLeaf() bool {
	return len(op.childrenMap) == 0
}

func (op *baseOperator) printTree(depth int) {
	switch {
		case depth > int(config.GlobalParameters.MaxTreeDepth):
		fmt.Printf("ERROR: MaxTreeDepth exceeded: %d\n", config.GlobalParameters.MaxTreeDepth)
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


// TODO: Validate
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


// Disconnets all child operators from op
//
func (op *baseOperator) DisconnectAll() {
	op.Panic()
	for _, child := range op.Children() {
		op.Disconnect(child)
	}
}

// Disconnects all child operators from op
// then recursivly call on child operators.
//
func (op *baseOperator) DisconnectTree() {
	op.Panic()
	children := op.Children()
	for _, child := range children {
		op.Disconnect(child)
		child.DisconnectTree()
	}
}
	
func (op *baseOperator) DisconnectParents() {
	for _, parent := range op.Parents() {
		parent.Disconnect(op)
	}
}


func (op *baseOperator) ChannelMode() midi.ChannelMode {
	return op.channelSelector.ChannelMode()
}

func (op *baseOperator) EnableChannel(channel midi.MIDIChannel, flag bool) error {
	return op.channelSelector.EnableChannel(channel, flag)
}

func (op *baseOperator) SelectChannel(channel midi.MIDIChannel) error {
	return op.channelSelector.SelectChannel(channel)
}

func (op *baseOperator) SelectedChannelIndexes() []midi.MIDIChannelIndex {
	return op.channelSelector.SelectedChannelIndexes()
}


func (op *baseOperator) ChannelIndexSelected(ci midi.MIDIChannelIndex ) bool {
	return op.channelSelector.ChannelIndexSelected(ci)
}

func (op *baseOperator) DeselectAllChannels() {
	op.channelSelector.DeselectAllChannels()
}

func (op *baseOperator) SelectAllChannels() {
	op.channelSelector.SelectAllChannels()
}


func (op *baseOperator) MIDIOutputEnabled() bool {
	return op.midiOutputEnabled
}

func (op *baseOperator) SetMIDIOutputEnabled(flag bool) {
	op.midiOutputEnabled = flag
}


func (op *baseOperator) Accept(event portmidi.Event) bool {
	return true
}


func (op *baseOperator) distribute(event portmidi.Event) {
	if op.MIDIOutputEnabled() {
		for _, child := range op.children() {
			child.Send(event)
		}
	}
}

func (op *baseOperator) Send(event portmidi.Event) {
	if op.Accept(event) {
		op.distribute(event)
	}
}


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

func (op *baseOperator) Commands() []string {
	var keys []string = make([]string, 0, len(op.dispatchTable))
	for k := range op.dispatchTable {
		keys = append(keys, k)
	}
	return keys
}

func (op *baseOperator) addCommandHandler(command string, handler func(*goosc.Message)([]string, error)) {
	op.dispatchTable[command] = handler
}
	

// osc /pig/op name, ping
// -> ACK
//
func (op *baseOperator) remotePing(*goosc.Message)([]string, error) {
	var err error
	var result = []string{op.Name(), "Ping"}
	return result, err
}

// osc /pig/op name, q-commands
// -> list of operator-local commands
//
func (op *baseOperator) remoteQueryCommands(*goosc.Message)([]string, error) {
	var err error
	return op.Commands(), err
}

