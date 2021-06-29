package op

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/rakyll/portmidi"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
)

var (
	empty []string
)

// Add general op-related handlers to global OSC server
//
func Init() {
	server := osc.GlobalServer
	osc.AddHandler(server, "ping", remotePing)
	osc.AddHandler(server, "exit", remoteExit)
	osc.AddHandler(server, "q-midi-inputs", remoteQueryMIDIInputs)
	osc.AddHandler(server, "q-midi-outputs", remoteQueryMIDIOutputs)
	osc.AddHandler(server, "batch", remoteBatchLoad)
	osc.AddHandler(server, "new", remoteNewOperator)
	osc.AddHandler(server, "delete-operator", remoteDeleteOperator)
	osc.AddHandler(server, "delete-all", remoteDeleteAllOperators)
	osc.AddHandler(server, "connect", remoteConnect)
	osc.AddHandler(server, "disconnect-child", remoteDisconnect)
	osc.AddHandler(server, "disconnect-all", remoteDisconnectAll)
	osc.AddHandler(server, "disconnect-parents", remoteDisconnectParents)
	osc.AddHandler(server, "reset-operator", remoteReset)
	osc.AddHandler(server, "reset-all", remoteResetAll)   // clash with 'reset'
	osc.AddHandler(server, "enable-midi", remoteEnableMIDI)
	osc.AddHandler(server, "q-midi-enabled", remoteQueryMIDIEnabled)
	osc.AddHandler(server, "q-channel-mode", remoteQueryChannelMode)
	osc.AddHandler(server, "q-channels", remoteQuerySelectedChannels)
	osc.AddHandler(server, "q-channel-selected", remoteQueryChannelSelected)
	osc.AddHandler(server, "select-channels", remoteSelectChannels)
	osc.AddHandler(server, "deselect-channels", remoteDeselectChannels)
	osc.AddHandler(server, "select-all-channels", remoteSelectAllChannels)
	osc.AddHandler(server, "deselect-all-channels", remoteDeselectAllChannels)
	osc.AddHandler(server, "invert-channels", remoteInvertChannelSelection)
	osc.AddHandler(server, "print-graph", remotePrintGraph)
        osc.AddHandler(server, "q-operator-types", remoteQueryOperatorTypes)
	osc.AddHandler(server, "q-operators", remoteQueryOperators)
	osc.AddHandler(server, "q-roots", remoteQueryRoots)
	osc.AddHandler(server, "q-graph", remoteQueryGraph)
	osc.AddHandler(server, "q-commands", remoteQueryCommands)
	osc.AddHandler(server, "q-children", remoteQueryChildren)
	osc.AddHandler(server, "q-parents", remoteQueryParents)

	
	osc.AddHandler(server, "op", dispatchExtendedCommand)
}

// osc /pig/ping -> ACK
// diagnostic function.
//
func remotePing(msg *goosc.Message)([]string, error) {
	var err error
	fmt.Printf("PING %s\n", msg.Address)
	return empty, err
}

// osc /pig/exit -> ACK
// Terminate application
//
func remoteExit(msg *goosc.Message)([]string, error) {
	var err error
	osc.Exit = true
	return empty, err
}

// osc /pig/q-midi-inputs
// -> ACK list of MIDI input devices
//
func remoteQueryMIDIInputs(msg *goosc.Message)([]string, error) {
	var err error
	ids := midi.InputIDs()
	acc := make([]string, len(ids))
	for i, id := range ids {
		info := portmidi.Info(id)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	return acc, err
}


// osc /pig/q-midi-outputs
// -> ACK list of MIDI output devices
//
func remoteQueryMIDIOutputs(msg *goosc.Message)([]string, error) {
	var err error
	ids := midi.OutputIDs()
	acc := make([]string, len(ids))
	for i, id := range ids {
		info := portmidi.Info(id)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	return acc, err
}


// osc /pig/batch filename
// 
//
func remoteBatchLoad(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("s", msg)
	if err != nil {
		fmt.Print(config.GlobalParameters.ErrorColor)
		fmt.Printf("ERROR: %s\n", msg.Address)
		fmt.Printf("ERROR: %s\n", err)
		fmt.Print(config.GlobalParameters.TextColor)
		return empty, err
	}
	filename := fmt.Sprintf("%s", args[0])
	err = osc.BatchLoad(filename)
	return empty, err
}


// osc /pig/new optype name
//     /pig/new MidiInput name device
//     /pig/new MidiOutput name device
// -> name
// Not used for MIDIInput or MIDIOutput
//
func remoteNewOperator(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("ss", msg)
	if err != nil {
		return empty, err
	}
	otype, name := args[0], args[1]
	switch otype {
	case "MIDIInput":
		return remoteNewMIDIInput(args)
	case "MIDIOutput":
		return remoteNewMIDIOutput(args)
	default:
		op, err := NewOperator(otype, name)
		if err != nil {
			return empty, err
		}
		return StringSlice(op.Name()), err
	}
}

// args ["input", name, device]
// -> name
//
func remoteNewMIDIInput(args []string)([]string, error) {
	args, err := Expect("sss", args)
	if err != nil {
		return empty, err
	}
	name, device := args[1], args[2]
	op, err := NewMIDIInput(name, device)
	if err != nil {
		return args, err
	}
	return StringSlice(op.Name()), err
}

// args ["output", name, device]
// -> name
//
func remoteNewMIDIOutput(args []string)([]string, error) {
	args, err := Expect("sss", args)
	if err != nil {
		return empty, err
	}
	name, device := args[1], args[2]
	op, err := NewMIDIOutput(name, device)
	if err != nil {
		return args, err
	}
	return StringSlice(op.Name()), err
}


// osc /pig/delete-operator name
// -> ACK
//
func remoteDeleteOperator(msg *goosc.Message)([]string, error) {
	var err error
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	DeleteOperator(args[0])
	return empty, err
}


// osc /pig/reset name
// -> ACK
//
func remoteReset (msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	op.Reset()
	return empty, err
}

// osc /pig/reset-all
// -> ACK
//
func remoteResetAll(msg *goosc.Message)([]string, error) {
	var err error
	for _, op := range Operators() {
		op.Reset()
	}
	return empty, err
}


// osc /pig/enable-midi name bool
// -> bool
//
func remoteEnableMIDI(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("ob", msg)
	if err != nil {
		return empty, err
	}
	name := args[0]
	flag, _ := strconv.ParseBool(args[1])
	op, _ := GetOperator(name)
	op.SetMIDIOutputEnabled(flag)
	rs := []string{fmt.Sprintf("%v", flag)}
	return rs, err
	return empty, err
}

// osc /pig/q-midi-enabled name
// -> bool
//
func remoteQueryMIDIEnabled(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	acc := []string{fmt.Sprintf("%v", op.MIDIOutputEnabled())}
	return acc, err
}

// osc /pig/q-midi-channel-mode name
// -> mode
//
func remoteQueryChannelMode(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	mode := op.ChannelMode().String()
	return []string{mode}, err
}

// osc /pig/q-channels name
// -> list
//
func remoteQuerySelectedChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	clist := op.SelectedChannelIndexes()
	acc := make([]string, len(clist))
	for i, ci := range clist {
		acc[i] = fmt.Sprintf("%d", int(ci+1))
	}
	return acc, err
}

// osc /pig/select-channels name [channels ....]
// -> list of enabled channels
//
func remoteSelectChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("oi", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	args = ToStringSlice(msg.Arguments)
	for _, s := range args[1:] {
		n, err := strconv.Atoi(s)
		if err != nil {
			msg := "Expected MIDI channel, got '%v'"
			err = errors.New(fmt.Sprintf(msg, s))
			return empty, err
		}
		if n < 1 || 16 < n {
			msg := "Expected MIDI channel, got '%v'"
			err = errors.New(fmt.Sprintf(msg, s))
			return empty, err
		}
		op.EnableChannel(midi.MIDIChannel(n), true)
	}
	clist := op.SelectedChannelIndexes()
	acc := make([]string, len(clist))
	for i, ci := range clist {
		acc[i] = fmt.Sprintf("%d", ci+1)
	}
	return acc, err
}

// osc /pig/deselect-channels name [channels ....]
// -> list of enabled channels
//
func remoteDeselectChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("oi", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	args = ToStringSlice(msg.Arguments)
	for _, s := range args[1:] {
		n, err := strconv.Atoi(s)
		if err != nil {
			msg := "Expected MIDI channel, got '%v'"
			err = errors.New(fmt.Sprintf(msg, s))
			return empty, err
		}
		if n < 1 || 16 < n {
			msg := "Expected MIDI channel, got '%v'"
			err = errors.New(fmt.Sprintf(msg, s))
			return empty, err
		}
		op.EnableChannel(midi.MIDIChannel(n), false)
	}
	clist := op.SelectedChannelIndexes()
	acc := make([]string, len(clist))
	for i, ci := range clist {
		acc[i] = fmt.Sprintf("%d", ci+1)
	}
	return acc, err
}

// osc /pig/select-all-channels name
// -> ACK
//
func remoteSelectAllChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	for i:=1; i<17; i++ {
		op.EnableChannel(midi.MIDIChannel(i), true)
	}
	return empty, err
}

// osc /pig/deselect-all-channels name
// -> ACK
//
func remoteDeselectAllChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	for i:=1; i<17; i++ {
		op.EnableChannel(midi.MIDIChannel(i), false)
	}
	return empty, err
}

// osc /pig/invert-channels name
// -> list of selected channels
//
func remoteInvertChannelSelection(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	for i:=0; i<16; i++ {
		flag := op.ChannelIndexSelected(midi.MIDIChannelIndex(i))
		op.EnableChannel(midi.MIDIChannel(i+1), flag)
	}
	clist := op.SelectedChannelIndexes()
	acc := make([]string, len(clist))
	for i, ci := range clist {
		acc[i] = fmt.Sprintf("%d", int(ci+1))
	}
	return acc, err
}
	
// osc /pig/q-channel-selected name c
// -> bool
//
func remoteQueryChannelSelected(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("oc", msg)
		if err != nil {
		return empty, err
	}
	op, _ := GetOperator(args[0])
	s, _ := strconv.Atoi(args[1])
	ci := midi.MIDIChannelIndex(s-1)
	flag := op.ChannelIndexSelected(ci)
	acc := []string{fmt.Sprintf("%v", flag)}
	return acc, err
}
	

// osc /pig/q-operators
// -> list
//
func remoteQueryOperators(msg *goosc.Message)([]string, error) {
	var err error
	ops := Operators()
	acc := make([]string, len(ops))
	for i, op := range ops {
		acc[i] = fmt.Sprintf("%s, %s", op.OperatorType(), op.Name())
	}
	return acc, err
}


// osc /pig/q-roots
// -> list
//
func remoteQueryRoots(msg *goosc.Message)([]string, error) {
	var err error
	ops := RootOperators()
	acc := make([]string, len(ops))
	for i, op := range ops {
		acc[i] = fmt.Sprintf("%s, %s", op.OperatorType(), op.Name())
	}
	return acc, err
}
	
// osc /pig/q-operator-types
// -> list
//
func remoteQueryOperatorTypes(msg *goosc.Message)([]string, error) {
	var err error
	return OperatorTypes(false), err
}

// osc /pig/delete-all-operators
// -> ACK
//
func remoteDeleteAllOperators(msg *goosc.Message)([]string, error) {
	var err error
	ClearRegistry()
	return empty, err
}


		
// osc /pig/connect parent child [more...]
// -> ACK
//
func remoteConnect(msg *goosc.Message)([]string, error) {
	var err error
	var parent, child Operator
	args := ToStringSlice(msg.Arguments)
	if len(args) < 2 {
		err = errors.New(fmt.Sprintf("Expected at least parent/child operator pair."))
		return empty, err
	}
	for i:=1; i<len(args); i++ {
		parent, err = GetOperator(args[i-1])
		if err != nil {
			return empty, nil
		}
		child, err = GetOperator(args[i])
		if err != nil {
			return empty, err
		}
		parent.Connect(child)
	}
	return empty, err
}

// osc /pig/disconnect-child parent child
// -> ACK | ERROR
//
func remoteDisconnect(msg *goosc.Message)([]string, error) {
	var err error
	var parent, child Operator
	args, err := ExpectMsg("oo", msg)
	if err != nil {
		return empty, err
	}
	parentName, childName := args[0], args[1]
	parent, _ = GetOperator(parentName)
	child, _ = GetOperator(childName)
	parent.Disconnect(child)
	child.Panic()
	return empty, err
}


// osc /pig/disconnect-all parent
// -> ACK | ERROR
//
func remoteDisconnectAll(msg *goosc.Message)([]string, error) {
	var err error
	var parent Operator
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	name := args[0]
	parent, _ = GetOperator(name)
	parent.DisconnectAll()
	return empty, err
}

// osc /pig/disconnect-parents name
// -> ACK | ERROR
//
func remoteDisconnectParents(msg *goosc.Message)([]string, error) {
	var err error
	var op Operator
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	name := args[0]
	op, _ = GetOperator(name)
	op.DisconnectParents()
	op.Panic()
	return empty, err
}

// osc /pig/print-graph
// -> ACK
//
func remotePrintGraph(msg *goosc.Message)([]string, error) {
	var err error
	for _, op := range RootOperators() {
		op.PrintTree()
	}
	return empty, err
}

// osc /pig/q-forest
// -> list
//
func remoteQueryGraph(msg *goosc.Message)([]string, error) {
	var err error
	var seen = make(map[string]bool)
	var acc []string
	for _, op := range Operators() {
		name := op.Name()
		_, flag := seen[name]
		if !flag {
			seen[name]=true
			if !op.IsLeaf() {
				s := fmt.Sprintf("%s -> [", name)
				for _, c := range op.Children() {
					s += fmt.Sprintf("%s, ", c.Name())
				}
				s = s[0:len(s)-2] + "]"
				acc = append(acc, s)
			} else {
				acc = append(acc, op.Name())
			}
		}
	}
	return acc, err
}
			
// osc /pig/q-commands
// -> list
//
func remoteQueryCommands(msg *goosc.Message)([]string, error) {
	var err error
	acc := osc.GlobalServer.Commands()
	for _, op := range Operators() {
		name := op.Name()
		for _, cmd := range op.Commands() {
			acc = append(acc, fmt.Sprintf("/pig/op %s, %s", name, cmd))
		}
	}
	return acc, err
}

// psc /pig/q-children name
// -> list
//
func remoteQueryChildren(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	name := args[0]
	op, _ := GetOperator(name)
	children := op.children()
	acc := make([]string, len(children))
	i := 0
	for key, _ := range children {
		acc[i] = key
		i++
	}
	return acc, err
}

// psc /pig/q-parents name
// -> list
//
func remoteQueryParents(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	name := args[0]
	op, _ := GetOperator(name)
	parents := op.parents()
	acc := make([]string, len(parents))
	i := 0
	for key, _ := range parents {
		acc[i] = key
		i++
	}
	return acc, err
}
	


// osc  /pig/op  [name, command, arguments...]
//
func dispatchExtendedCommand(msg *goosc.Message)([]string, error) {
	var err error
	var args []string
	var op Operator
	args, err = ExpectMsg("os", msg)
	if err != nil {
		return empty, err
	}
	name, command := args[0], args[1]
	op, _ = GetOperator(name)
	result, rerr := op.DispatchCommand(command, ToStringSlice(msg.Arguments))
	return result, rerr
}


