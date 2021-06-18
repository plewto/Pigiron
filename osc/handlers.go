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

func (s *PigServer) sendError(err error, msg *goosc.Message) {
	fmt.Printf("OSC ERROR: %s\n", msg.Address)
	fmt.Printf("%v\n", err)
	errmsg := toSlice(err)
	s.client.Error(msg.Address, errmsg)
}


func (s *PigServer) ping(msg *goosc.Message) {
	fmt.Printf("%s\n", msg.Address)
	s.client.Ack(msg.Address, empty)
}


func (s *PigServer) exit(msg *goosc.Message) {
	s.client.Ack(msg.Address, empty)
	Exit = true
}


// /pig/new-operator <opType> <name>
// --> actual name
//
func (s *PigServer) newOperator(msg *goosc.Message) {
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
func (s *PigServer) newMIDIOutput(msg *goosc.Message) {
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
func (s *PigServer) newMIDIInput(msg *goosc.Message) {
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

func (s *PigServer) deleteOperator(msg *goosc.Message) {
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
func (s *PigServer) connect(msg *goosc.Message) {
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
func (s *PigServer) disconnect(msg *goosc.Message) {
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
func (s *PigServer) disconnectAll(msg *goosc.Message) {
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

func (s *PigServer) destroyForest(msg *goosc.Message) {
	op.DestroyForest()
	s.client.Ack(msg.Address, empty)
}

func (s *PigServer) printForest(msg *goosc.Message) {
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
func (s *PigServer) queryIsParent(msg *goosc.Message) {
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
		

func (s *PigServer) queryMIDIInputs(msg *goosc.Message) {
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

func (s *PigServer) queryMIDIOutputs(msg *goosc.Message) {
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
func (s *PigServer) queryOperators(msg *goosc.Message) {
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
func (s *PigServer) queryRoots(msg *goosc.Message) {
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
func (s *PigServer) queryChildren(msg *goosc.Message) {
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
func (s *PigServer) queryParents(msg *goosc.Message) {
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
			

func (s *PigServer) panic(msg *goosc.Message) {
	op.PanicAll()
	s.client.Ack(msg.Address, empty)
}

func (s *PigServer) reset(msg *goosc.Message) {
	op.ResetAll()
	s.client.Ack(msg.Address, empty)
}

func (s *PigServer) help(msg *goosc.Message) {
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
	

