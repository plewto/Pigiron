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

// func (s *PigServer) newOperator(msg *goosc.Message) {
// 	template := []expectType{xpString, xpString}
// 	args, err := expect(template, msg.Arguments)                    // have expect return []string instead of []interface
// 	if err != nil {
// 		errmsg := []string{fmt.Sprintf("%v", err)}              // define string slicer builder    fn(args ...interface{}) --> []string
// 		sargs := "Arguments were: "
// 		for _, a := range msg.Arguments {
// 			sargs += fmt.Sprintf("%v, ", a)
// 		}
// 		errmsg = append(errmsg, sargs)
// 		s.client.Error(msg.Address, errmsg)
// 	} else {
// 		var opp op.Operator
// 		opType := fmt.Sprintf("%s", args[0])
// 		opName := fmt.Sprintf("%s", args[1])
// 		opp, err = op.NewOperator(opType, opName)
// 		if err != nil {
// 			errMsg := fmt.Sprintf("%v", err)
// 			s.client.Error(msg.Address, []string{errMsg})
// 		} else {
			
// 			s.client.Ack(msg.Address, []string{opp.Name()})
// 		}
// 	}
// }


func (s *PigServer) newOperator(msg *goosc.Message) {
	template := []expectType{xpString, xpString}
	args, err := expect(template, msg.Arguments)
	if err != nil {
		errmsg := []string{fmt.Sprintf("%v", err)}
		sargs := "Arguments were: "
		for _, a := range msg.Arguments {
			sargs += fmt.Sprintf("%v, ", a)
		}
		errmsg = append(errmsg, sargs)
		s.client.Error(msg.Address, errmsg)
	} else {
		var opp op.Operator
		opType := fmt.Sprintf("%s", args[0])
		opName := fmt.Sprintf("%s", args[1])
		opp, err = op.NewOperator(opType, opName)
		if err != nil {
			errMsg := fmt.Sprintf("%v", err)
			s.client.Error(msg.Address, []string{errMsg})
		} else {
			
			s.client.Ack(msg.Address, []string{opp.Name()})
		}
	}
}


