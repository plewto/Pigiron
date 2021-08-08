package op

/*
** Defines global OSC handler functions.
**
*/

import (
	"fmt"
	"strconv"
	"github.com/rakyll/portmidi"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/help"
)

var empty []string


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
	osc.AddHandler(server, "del-op", remoteDeleteOperator)
	osc.AddHandler(server, "del-all", remoteDeleteAllOperators)
	osc.AddHandler(server, "connect", remoteConnect)
	osc.AddHandler(server, "disconnect-child", remoteDisconnect)
	osc.AddHandler(server, "disconnect-all", remoteDisconnectAll)
	osc.AddHandler(server, "disconnect-parents", remoteDisconnectParents)
	osc.AddHandler(server, "reset-op", remoteReset)
	osc.AddHandler(server, "reset-all", remoteResetAll)  
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
	osc.AddHandler(server, "print-info", remotePrintInfo)
	osc.AddHandler(server, "print-config", remotePrintConfig)
	osc.AddHandler(server, "midi", remoteMIDIInsert)
	osc.AddHandler(server, "op", dispatchExtendedCommand)
	osc.AddHandler(server, "help", remoteHelp)
}

// remotePing() handler for /pig/ping
// diagnostic function.
// osc returns ACK
//
func remotePing(msg *goosc.Message)([]string, error) {
	var err error
	fmt.Printf("PING %s\n", msg.Address)
	return empty, err
}


// remoteExit() handler for /pig/exit
// osc returns ACK
//
func remoteExit(msg *goosc.Message)([]string, error) {
	var err error
	osc.Exit = true
	return empty, err
}

// remoteQueryMIDIInputs() handler for /pig/q-midi-inputs
// osc returns list of MIDI input devices
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


// remoteQueryMIDIOutputs() handler for /pig/q-midi-outputs
// osc returns list of MIDI output devices
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


// remoteBatchLoad() handler for /pig/batch
// /pig/batch <filename>
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
	filename := args[0].S
	err = osc.BatchLoad(filename)
	return empty, err
}

// remoteNewOperator() handler for /pig/new
// Creates new Operator
//
// This command has two forms:
//
// General form for non-io operators.
//
//     /pig/new <operator-type>, <name>
//
// For MIDIInput and MIDIOutput operators, a device must be specified.
//
//     /pig/new MIDIInput, <name>, <device>
//     /pig/new MIDIOutput, <name>, <device>
//
//     The device may either be an integer index or a sub-string of
//     the device's name.
//
// osc returns actual operator's name.
//
func remoteNewOperator(msg *goosc.Message)([]string, error) {
	if len(msg.Arguments) > 2 {
		args, err := ExpectMsg("sss", msg)
		if err != nil {
			return empty, err
		}
		return makeIOOperator(args)
	} else {
		// non-io operator  {opType, name}
		args, err := ExpectMsg("ss", msg)
		if err != nil {
			return empty, err
		}
		optype, name := args[0].S, args[1].S
		op, err := NewOperator(optype, name)
		if err != nil {
			return empty, err
		}
		return []string{op.Name()}, err
	}
}

// helper for remoteNewOperator
//
func makeIOOperator(args []ExpectValue)([]string, error) {
	optype, name, device := args[0].S, args[1].S, args[2].S
	switch optype {
	case "MIDIInput":
		op, err := NewMIDIInput(name, device)
		if err != nil {
			return empty, err
		}
		return []string{op.Name()}, err
	case "MIDIOutput":
		op, err := NewMIDIOutput(name, device)
		if err != nil {
			return empty, err
		}
		return []string{op.Name()}, err
	default:
		msg := "Expected operator type at index 0, got %s"
		err := fmt.Errorf(msg, optype)
		return empty, err
	}
}


// remoteDeleteOperator() handler for /pig/op-del
// Delete operator
// osc /pig/op-del <name>
// osc returns ACK
//
func remoteDeleteOperator(msg *goosc.Message)([]string, error) {
	var err error
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	err = DeleteOperator(args[0].O.Name())
	return empty, err
}

// remoteReset() handler for /pig/reset-op 
// Resets operator.
// osc /pig/reset-op <name>
// osc returns ACK
//
func remoteReset (msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	op.Reset()
	return empty, err
}

// remoteResetAll() handler for /pig/reset-all
// Resets all operators.
// osc /pig/reset-all
// osc returns ACK
//
func remoteResetAll(msg *goosc.Message)([]string, error) {
	var err error
	for _, op := range Operators() {
		op.Reset()
	}
	return empty, err
}

// remoteEnableMIDI() handler for /pig/enable-midi
// Enables/disables operator MIDI output.
// osc /pig/enable-midi <name>, <bool>
// osc returns bool.
//
func remoteEnableMIDI(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("ob", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	flag := args[1].B
	op.SetMIDIOutputEnabled(flag)
	rs := []string{fmt.Sprintf("%v", flag)}
	return rs, err
	return empty, err
}

// remoteQueryMIDIEnabled() handler for /pig/q-midi-enabled
// Gets state of operator midi-enabled flag.
// osc /pig/q-midi-enabled <name>
// osc returns bool
//
func remoteQueryMIDIEnabled(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	acc := []string{fmt.Sprintf("%v", op.MIDIOutputEnabled())}
	return acc, err
}


// remoteQueryChannelMode() handler for /pig/q-midi-channel-mode
// osc /pig/q-midi-channel-mode <name>
// osc returns the ChannelSelector mode
//
func remoteQueryChannelMode(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	mode := op.ChannelMode().String()
	return []string{mode}, err
}


// remoteQuerySelectedChannels() handler for /pig/q-channels
// osc /pig/q-channels name
// osc returns list of enabled MIDI channels.
//
func remoteQuerySelectedChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	clist := op.SelectedChannelIndexes()
	acc := make([]string, len(clist))
	for i, ci := range clist {
		acc[i] = fmt.Sprintf("%d", int(ci+1))
	}
	return acc, err
}


// remoteSelectChannels() handler for /pig/select-channels
// osc /pig/select-channels <name>, <channels, ....>
// osc returns list of enabled channels.
//
func remoteSelectChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("oi", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	cargs := ToStringSlice(msg.Arguments)
	for _, s := range cargs[1:] {
		n, err := strconv.Atoi(s)
		if err != nil {
			msg := "Expected MIDI channel, got '%v'"
			err = fmt.Errorf(msg, s)
			return empty, err
		}
		if n < 1 || 16 < n {
			msg := "Expected MIDI channel, got '%v'"
			err = fmt.Errorf(msg, s)
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

// remoteDeselectChannels() handler for /pig/deselect-channels 
// osc /pig/deselect-channels <name>, <channels, ....>
// osc returns list of enabled channels.
//
func remoteDeselectChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("oi", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	cargs := ToStringSlice(msg.Arguments)
	for _, s := range cargs[1:] {
		n, err := strconv.Atoi(s)
		if err != nil {
			msg := "Expected MIDI channel, got '%v'"
			err = fmt.Errorf(msg, s)
			return empty, err
		}
		if n < 1 || 16 < n {
			msg := "Expected MIDI channel, got '%v'"
			err = fmt.Errorf(msg, s)
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


// remoteSelectAllChannels() handler for /pig/select-all-channels
// osc /pig/select-all-channels <name>
// osc returns ACK
//
func remoteSelectAllChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	for i:=1; i<17; i++ {
		op.EnableChannel(midi.MIDIChannel(i), true)
	}
	return empty, err
}



// remoteDeselectAllChannels() handler for /pig/deselect-all-channels
// osc /pig/deselect-all-channels <name>
// osc returns ACK
//
func remoteDeselectAllChannels(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	for i:=1; i<17; i++ {
		op.EnableChannel(midi.MIDIChannel(i), false)
	}
	return empty, err
}


// remoteInvertChannelSelection() handler for /pig/invert-channels
// Inverts MIDI channel selection.
// osc /pig/invert-channels <name>
// osc returns list of enabled channel
//
func remoteInvertChannelSelection(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	for i:=0; i<16; i++ {
		flag := op.ChannelIndexSelected(midi.MIDIChannelNibble(i))
		op.EnableChannel(midi.MIDIChannel(i+1), flag)
	}
	clist := op.SelectedChannelIndexes()
	acc := make([]string, len(clist))
	for i, ci := range clist {
		acc[i] = fmt.Sprintf("%d", int(ci+1))
	}
	return acc, err
}


// remoteQueryChannelSelected() handler for /pig/q-channel-selected
// Checks if specific MIDI channel is selected.
// osc /pig/q-channel-selected <name>,  <channel>
// osc returns bool
//
func remoteQueryChannelSelected(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("oc", msg)
		if err != nil {
		return empty, err
	}
	op := args[0].O
	c := args[1].C
	ci := midi.MIDIChannelNibble(c-1)
	flag := op.ChannelIndexSelected(ci)
	acc := []string{fmt.Sprintf("%v", flag)}
	return acc, err
}
	

// remoteQueryOperators() handler for /pig/q-operators
// osc /pig/q-operators
// osc returns list of all operator names.
//
func remoteQueryOperators(msg *goosc.Message)([]string, error) {
	var err error
	ops := Operators()
	acc := make([]string, len(ops))
	for i, op := range ops {
		acc[i] = fmt.Sprintf("%s", op)
	}
	return acc, err
}


// remoteQueryRoots() handler for /pig/q-roots
// osc /pig/q-roots
// osc returns names for all root operators.
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
	
// remoteQueryOperatorTypes() handler for /pig/q-operator-types
// osc /pig/q-operator-types
// osc returns list of available operator types.
//
func remoteQueryOperatorTypes(msg *goosc.Message)([]string, error) {
	var err error
	return OperatorTypes(false), err
}

// remoteDeleteAllOperators() handler for /pig/del-all
// Deletes all operators.
// osc /pig/del-all
// osc returns ACK.
//
func remoteDeleteAllOperators(msg *goosc.Message)([]string, error) {
	var err error
	ClearRegistry()
	return empty, err
}


		
// remoteConnect() handler for /pig/connect
// Connects two or more operators
// osc /pig/connect <parent>, <child-1> <,child-2, child-3, ...>
// osc returns ACK
//
// If more then one child is specified, they are connected in sequence.
// parent -> child-1 -> child-2 ... -> child-n
//
func remoteConnect(msg *goosc.Message)([]string, error) {
	var err error
	var parent, child Operator
	args := ToStringSlice(msg.Arguments)
	if len(args) < 2 {
		err = fmt.Errorf("Expected at least parent & child operator pair.")
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



// remoteDisconnect() handler for /pig/disconnect-child
// osc /pig/disconnect-child <parent>, <child>
// osc returns ACK
//
func remoteDisconnect(msg *goosc.Message)([]string, error) {
	var err error
	// var parent, child Operator
	args, err := ExpectMsg("oo", msg)
	if err != nil {
		return empty, err
	}
	parent := args[0].O
	child := args[1].O
	parent.Disconnect(child)
	child.Panic()
	return empty, err
}


// remoteDisconnectAll() handler for /pig/disconnect-all
// Disconnects all children from parent.
// osc /pig/disconnect-all <parent>
// osc returns ACK
//
func remoteDisconnectAll(msg *goosc.Message)([]string, error) {
	var err error
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	parent := args[0].O
	parent.DisconnectAll()
	return empty, err
}


// remoteDisconnectParents() handler for /pig/disconnect-parents
// Disconnects all parents from operator.
// osc /pig/disconnect-parents <name>
// osc returns ACK.
//
func remoteDisconnectParents(msg *goosc.Message)([]string, error) {
	var err error
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	op.DisconnectParents()
	op.Panic()
	return empty, err
}

// remotePrintGraph() handler for /pig/print-graph
// Prints a graphic representation of the MIDI process tree.
// osc /pig/print-graph
// osc returns ACK.
//
func remotePrintGraph(msg *goosc.Message)([]string, error) {
	var err error
	for _, op := range RootOperators() {
		op.PrintTree()
		fmt.Println()
	}
	return empty, err
}

// remoteQueryGraph() handler for /pig/q-forest
// osc /pig/q-forest
// osc returns list of all operator connections.
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
			
// remoteQueryCommands() handler for /pig/q-commands
// Prints list of all osc commands.
// osc /pig/q-commands
// osc returns list of all commands.
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

// remoteQueryChildren() handler for /pig/q-children
// Prints list of operator's children.
// osc /pig/q-children <name>
// osc returns list of children.
//
func remoteQueryChildren(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	children := op.children()
	acc := make([]string, len(children))
	i := 0
	for key, _ := range children {
		acc[i] = key
		i++
	}
	return acc, err
}

// remoteQueryParents() handler for /pig/q-parents
// Prints list of operator's parents.
// osc /pig/q-parents <name>
// osc returns list of parents
//
func remoteQueryParents(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	parents := op.parents()
	acc := make([]string, len(parents))
	i := 0
	for key, _ := range parents {
		acc[i] = key
		i++
	}
	return acc, err
}

// remotePrintInfo() handler for /pig/print-info
// Prints operator's info.
// osc /pig/print-info <name>
// osc returns ACK.
//
func remotePrintInfo(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("o", msg)
	if err != nil {
		return empty, err
	}
	fmt.Println(args[0].O.Info())
	return empty, err
}

// remotePrintConfig() handler for /pig/print-config
// Prints global configuration values.
// osc /pig/print-config
// osc returns ACK.
//
func remotePrintConfig(msg *goosc.Message)([]string, error) {
	var err error
	config.PrintConfig()
	return empty, err
}

// remoteMIDIInsert handler for /pig/midi
// Sends MIDI events to specific operator.
// List of bytes can not include meta events and can not use running-status.
// osc /pig/midi <name>, <bytes, ....>
// osc returns ACK
//
func remoteMIDIInsert(msg *goosc.Message)([]string, error) {
	template := "o"
	for i := 0; i < len(msg.Arguments)-1; i++ {
		template += "i"
	}
	args, err := ExpectMsg(template, msg)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return empty, err
	}
	op := args[0].O
	bargs := args[1:]
	bytes := make([]byte, len(bargs))
	for i := 0; i < len(bargs); i++ {
		bytes[i] = byte(bargs[i].I)
	}
	var events []*midi.UniversalEvent
	events, err = midi.BytesToEvents(bytes)
	if err != nil {
		return empty, err
	}
	for _, event := range events {
		if !event.IsMetaEvent() {
			op.Send(event.PortmidiEvent())
		}
	}
	return empty, err
}
				
	
// dispatchExtendedCommand() handler for /pig/op
// Sends command to specific operator.  The general form is
// osc /pig/op <name>, <sub-command>  <,argument-1, argument-2, ..., argument-n>
// osc returns result of the sub-command.
//
func dispatchExtendedCommand(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("os", msg)
	if err != nil {
		return empty, err
	}
	op := args[0].O
	command := args[1].S
	result, rerr := op.DispatchCommand(command, msg)
	return result, rerr
}

func remoteHelp(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("s", msg)
	if err != nil {
		return empty, err
	}
	var topic = args[0].S
	var text = ""
	text, err = help.Help(topic)
	if err != nil {
		return empty, err
	}
	text = fmt.Sprintf("\n\n%s", text)
	return []string{text}, err
}
		
