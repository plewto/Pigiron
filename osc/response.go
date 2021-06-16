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
)


func (s *PigServer) ping(msg *goosc.Message) {
	fmt.Printf("%s\n", msg.Address)
	s.client.Ack(msg.Address, empty)
}


func (s *PigServer) exit(msg *goosc.Message) {
	s.client.Ack(msg.Address, empty)
	Exit = true
}

