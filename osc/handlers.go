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


func (s *OSCServer) ping(msg *goosc.Message) {
	ClearError()
	fmt.Printf("%s\n", msg.Address)
	s.client.Ack(msg.Address, empty)
}


func (s *OSCServer) exit(msg *goosc.Message) {
	ClearError()
	s.client.Ack(msg.Address, empty)
	Exit = true
}


// /pig/new-operator <opType> <name>
// --> actual name
//
func (s *OSCServer) newOperator(msg *goosc.Message) {
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


// /pig/new-midi-output <device> <name>
// --> actual name
//
func (s *OSCServer) newMIDIOutput(msg *goosc.Message) {
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

// /pig/new-midi-input <device> <name>
// --> actual name
//
func (s *OSCServer) newMIDIInput(msg *goosc.Message) {
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

func (s *OSCServer) deleteOperator(msg *goosc.Message) {
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

// /pig/connect <parent-name> <child-name>
//
func (s *OSCServer) connect(msg *goosc.Message) {
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

// /pig/disconnect <parent> <child>
//
func (s *OSCServer) disconnect(msg *goosc.Message) {
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

// /pig/disconnect-all <operator>
// --> Ack | Error
//
func (s *OSCServer) disconnectAll(msg *goosc.Message) {
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

func (s *OSCServer) destroyForest(msg *goosc.Message) {
	ClearError()
	op.DestroyForest()
	s.client.Ack(msg.Address, empty)
}

func (s *OSCServer) printForest(msg *goosc.Message) {
	ClearError()
	for _, root := range op.RootOperators() {
		fmt.Println("\nOperator tree:")
		root.PrintTree()
	}
	fmt.Println()
	s.client.Ack(msg.Address, empty)
}


// /pig/q-is-parent <parent> <child>
// --> bool
//
func (s *OSCServer) queryIsParent(msg *goosc.Message) {
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
		

func (s *OSCServer) queryMIDIInputs(msg *goosc.Message) {
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

func (s *OSCServer) queryMIDIOutputs(msg *goosc.Message) {
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

// /pig/q-operators
// --> ACK <list-of-operators>
//
func (s *OSCServer) queryOperators(msg *goosc.Message) {
	ClearError()
	oplist := op.Operators()
	acc := make([]string, len(oplist))
	for i, p := range oplist {
		acc[i] = p.Name()
	}
	s.client.Ack(msg.Address, acc)
}

// /pig/q-roots
// --> ACK <list-of-operators>
//
func (s *OSCServer) queryRoots(msg *goosc.Message) {
	ClearError()
	oplist := op.RootOperators()
	acc := make([]string, len(oplist))
	for i, p := range oplist {
		acc[i] = p.Name()
	}
	s.client.Ack(msg.Address, acc)
}

// /pig/q-children <operator>
// --> Ack <list-of-operators>
//
func (s *OSCServer) queryChildren(msg *goosc.Message) {
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

// /pig/q-parents <operator>
// --> Ack <list-of-operators>
//
func (s *OSCServer) queryParents(msg *goosc.Message) {
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
			

func (s *OSCServer) panic(msg *goosc.Message) {
	ClearError()
	op.PanicAll()
	s.client.Ack(msg.Address, empty)
}

func (s *OSCServer) reset(msg *goosc.Message) {
	ClearError()
	op.ResetAll()
	s.client.Ack(msg.Address, empty)
}

func (s *OSCServer) help(msg *goosc.Message) {
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
	
func (s *OSCServer) batchLoad(msg *goosc.Message) {
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
