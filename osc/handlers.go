package osc

/* 
 * response.go contains global OSC response functions.
 *
*/


import (
	"fmt"
)

import (
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"
	
	"github.com/plewto/pigiron/op"
	"github.com/plewto/pigiron/midi"
)



func (s *PigServer) ping(msg *goosc.Message) {
	fmt.Printf("%s\n", msg.Address)
	s.client.Ack(msg.Address, empty)
}


func (s *PigServer) exit(msg *goosc.Message) {
	s.client.Ack(msg.Address, empty)
	Exit = true
}

// /pig/new-op <opType> <name>
// --> actual name
//
func (s *PigServer) newOperator(msg *goosc.Message) {
	template := []expectType{xpString, xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		fmt.Printf("OSC ERROR: %s\n", msg.Address)
		fmt.Printf("%v\n", err)
		errmsg := toSlice(err)
		s.client.Error(msg.Address, errmsg)
	} else {
		var newOp op.Operator
		opType := args[0]
		opName := args[1]
		newOp, err = op.NewOperator(opType, opName)
		if err != nil {
			fmt.Printf("OSC ERROR: %s  args: %v\n", msg.Address, args)
			errmsg := toSlice(err)
			s.client.Error(msg.Address, errmsg)
			
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
		fmt.Printf("OSC ERROR: %s\n", msg.Address)
		fmt.Printf("%v\n", err)
		errmsg := toSlice(err)
		s.client.Error(msg.Address, errmsg)
	} else {
		var newOp op.Operator
		dev := args[0]
		name := args[1]
		newOp, err = op.NewMIDIOutput(name, dev)
		if err != nil {
			fmt.Printf("OSC ERROR: %s  args: %v\n", msg.Address, args)
			errmsg := toSlice(err)
			s.client.Error(msg.Address, errmsg)
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
		fmt.Printf("OSC ERROR: %s\n", msg.Address)
		fmt.Printf("%v\n", err)
		errmsg := toSlice(err)
		s.client.Error(msg.Address, errmsg)
	} else {
		var newOp op.Operator
		dev := args[0]
		name := args[1]
		newOp, err = op.NewMIDIInput(name, dev)
		if err != nil {
			fmt.Printf("OSC ERROR: %s  args: %v\n", msg.Address, args)
			errmsg := toSlice(err)
			s.client.Error(msg.Address, errmsg)
		} else {
			s.client.Ack(msg.Address, toSlice(newOp.Name()))
		}
	}
	// START DEBUG
	fmt.Println("---------------------------------------")
	for _, p := range op.Operators() {
		fmt.Println("+++")
		fmt.Println(p.Info())
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
		


