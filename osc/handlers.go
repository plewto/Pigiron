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
	"github.com/plewto/pigiron/op"
)



func (s *PigServer) ping(msg *goosc.Message) {
	fmt.Printf("%s\n", msg.Address)
	s.client.Ack(msg.Address, empty)
}


func (s *PigServer) exit(msg *goosc.Message) {
	s.client.Ack(msg.Address, empty)
	Exit = true
}

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
		
		
		
		
	


