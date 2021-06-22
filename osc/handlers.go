package osc

/* 
 * response.go contains global OSC response functions.
 *
*/


import (
	"fmt"

	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/op"
	"github.com/plewto/pigiron/midi"
)

func (s *OSCServer) sendError(err error, msg *goosc.Message) {
	signalError()
	fmt.Printf("OSC ERROR: %s\n", msg.Address)
	fmt.Printf("%v\n", err)
	errmsg := toSlice(err)
	s.client.Error(msg.Address, errmsg)
}


// osc /pig/ping -> ACK
// diagnostic function.
//
func (s *OSCServer) remotePing(msg *goosc.Message) {
	ClearError()
	fmt.Printf("%s\n", msg.Address)
	s.client.Ack(msg.Address, empty)
}

// osc /pig/exit -> ACK
// Terminate application
//
func (s *OSCServer) remoteExit(msg *goosc.Message) {
	ClearError()
	s.client.Ack(msg.Address, empty)
	Exit = true
}


// osc /pig/new-operator opType name
// -> ACK actual-name
// -> ERROR
//
func (s *OSCServer) remoteNewOperator(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString, xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		var newOp op.Operator
		opType := args[0]
		opName := args[1]
		newOp, err = op.NewOperator(opType, opName)
		if err != nil {
			s.sendError(err, msg)
		} else {
			s.client.Ack(msg.Address, toSlice(newOp.Name()))
		}
	}
}


// osc /pig/new-midi-output <device> <name>
// -> ACK actual-name
// -> ERROR
//
func (s *OSCServer) remoteNewMIDIOutput(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString, xpString}	
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		var newOp op.Operator
		dev := args[0]
		name := args[1]
		newOp, err = op.NewMIDIOutput(name, dev)
		if err != nil {
			s.sendError(err, msg)
		} else {
			s.client.Ack(msg.Address, toSlice(newOp.Name()))
		}
	}
}

// osc /pig/new-midi-input <device> <name>
// -> ACK actual-name
// -> ERROR
//
func (s *OSCServer) remoteNewMIDIInput(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString, xpString}	
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		var newOp op.Operator
		dev := args[0]
		name := args[1]
		newOp, err = op.NewMIDIInput(name, dev)
		if err != nil {
			s.sendError(err, msg)
		} else {
			s.client.Ack(msg.Address, toSlice(newOp.Name()))
		}
	}
}		

// osc /pig/delete-operator <name>
// -> ACK
// -> ERROR
//
func (s *OSCServer) remoteDeleteOperator(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		name := args[0]
		op.DeleteOperator(name)
		s.client.Ack(msg.Address, args)
	}
}

// osc /pig/connect <parent-name> <child-name>
// -> ACK
// -> ERROR
//
func (s *OSCServer) remoteConnect(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString, xpString}	
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		parent, err1 := op.GetOperator(args[0])
		if err1 != nil {
			s.sendError(err1, msg)
			return
		}
		child, err2 := op.GetOperator(args[1])
		if err2 != nil {
			s.sendError(err2, msg)
			return
		}
		err3 := parent.Connect(child)
		if err3 != nil {
			s.sendError(err3, msg)
		} else {
			s.client.Ack(msg.Address, empty)
		}
	}
}

// osc /pig/disconnect <parent> <child>
// -> ACK
// -> ERROR
//
func (s *OSCServer) remoteDisconnect(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString, xpString}	
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		parent, err1 := op.GetOperator(args[0])
		if err1 != nil {
			s.sendError(err1, msg)
			return
		}
		child, err2 := op.GetOperator(args[1])
		if err2 != nil {
			s.sendError(err2, msg)
			return
		}
		parent.Disconnect(child)
		s.client.Ack(msg.Address, empty)
	}
}

// osc /pig/disconnect-all <operator>
// -> ACK
// -> ERROR
//
func (s *OSCServer) remoteDisconnectAll(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	}
	p, err2 := op.GetOperator(args[0])
	if err2 != nil {
		s.sendError(err2, msg)
	} else {
		p.DisconnectAll()
		s.client.Ack(msg.Address, args)
	}
}

// osc /pig/destroy-forest
// -> ACK
//
func (s *OSCServer) remoteDestroyForest(msg *goosc.Message) {
	ClearError()
	op.DestroyForest()
	s.client.Ack(msg.Address, empty)
}


// osc /pig/print-forest
// -> ACK
//
func (s *OSCServer) remotePrintForest(msg *goosc.Message) {
	ClearError()
	for _, root := range op.RootOperators() {
		fmt.Println("\nOperator tree:")
		root.PrintTree()
	}
	fmt.Println()
	s.client.Ack(msg.Address, empty)
}


// osc /pig/q-is-parent <parent> <child>
// -> ACK bool
// -> ERROR
//
func (s *OSCServer) remoteQueryIsParent(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString, xpString}	
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		parent, err1 := op.GetOperator(args[0])
		if err1 != nil {
			s.sendError(err1, msg)
			return
		}
		child, err2 := op.GetOperator(args[1])
		if err2 != nil {
			s.sendError(err2, msg)
			return
		}
		flag := child.IsChildOf(parent)
		s.client.Ack(msg.Address, toSlice(flag))
	}
}
		
// osc /pig/q-midi-inputs
// -> ACK list of MIDI input devices
//
func (s *OSCServer) remoteQueryMIDIInputs(msg *goosc.Message) {
	ClearError()
	ids := midi.InputIDs()
	acc := make([]string, len(ids))
	fmt.Println("MIDI Input devices:")
	for i, id := range ids {
		info := portmidi.Info(id)
		fmt.Printf("\t%s\n", info.Name)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	s.client.Ack(msg.Address, acc)
}


// osc /pig/q-midi-outputs
// -> ACK list of MIDI output devices
//
func (s *OSCServer) remoteQueryMIDIOutputs(msg *goosc.Message) {
	ClearError()
	ids := midi.OutputIDs()
	acc := make([]string, len(ids))
	fmt.Println("MIDI Output devices:")
	for i, id := range ids {
		info := portmidi.Info(id)
		fmt.Printf("\t%s\n", info.Name)
		acc[i] = fmt.Sprintf("\"%s\" ", info.Name)
	}
	s.client.Ack(msg.Address, acc)
}

// osc /pig/q-operators
// -> ACK list of operators
//
func (s *OSCServer) remoteQueryOperators(msg *goosc.Message) {
	ClearError()
	oplist := op.Operators()
	acc := make([]string, len(oplist))
	for i, p := range oplist {
		acc[i] = p.Name()
	}
	s.client.Ack(msg.Address, acc)
}

// osc /pig/q-roots
// -> ACK list of root operators
//
func (s *OSCServer) remoteQueryRoots(msg *goosc.Message) {
	ClearError()
	oplist := op.RootOperators()
	acc := make([]string, len(oplist))
	for i, p := range oplist {
		acc[i] = p.Name()
	}
	s.client.Ack(msg.Address, acc)
}

// osc /pig/q-children <operator>
// -> ACK list of operators
// -> ERROR
//
func (s *OSCServer) remoteQueryChildren(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		p, err2 := op.GetOperator(args[0])
		if err2 != nil {
			s.sendError(err, msg)
		}
		children := p.Children()
		acc := make([]string, len(children))
		i := 0
		for key, _ := range children {
			acc[i] = key
			i++
		}
		s.client.Ack(msg.Address, acc)
	}
}

// osc /pig/q-parents <operator>
// -> ACK list of operators
//
func (s *OSCServer) remoteQueryParents(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		p, err2 := op.GetOperator(args[0])
		if err2 != nil {
			s.sendError(err, msg)
		}
		parents := p.Parents()
		acc := make([]string, len(parents))
		i := 0
		for key, _ := range parents {
			acc[i] = key
			i++
		}
		s.client.Ack(msg.Address, acc)
	}
}

// osc /pig/panic
// -> ACK
//
func (s *OSCServer) remotePanic(msg *goosc.Message) {
	ClearError()
	op.PanicAll()
	s.client.Ack(msg.Address, empty)
}

// osc /pig/reset
// -> ACK
//
func (s *OSCServer) remoteReset(msg *goosc.Message) {
	ClearError()
	op.ResetAll()
	s.client.Ack(msg.Address, empty)
}


// osc /pig/help topic
// no response
//
func (s *OSCServer) remoteHelp(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString}
	args, err := expect(template, msg.Arguments)
	topic := ""
	if err != nil {
		topic = "help"
	} else {
		topic = args[0]
	}
	fn, flag := oscHelp[topic]
	if flag {
		fn()
	} else {
		fmt.Println("Invalid help topic: %s", topic)
		fmt.Println("Try  'help help'")
	}
}

// osc /pig.batch filename
// no direct osc response.
//
func (s *OSCServer) remoteBatchLoad(msg *goosc.Message) {
	ClearError()
	template := []expectType{xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		s.sendError(err, msg)
	} else {
		filename := args[0]
		LoadBatchFile(filename)
	}
}
