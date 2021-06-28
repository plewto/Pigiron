package op

/*
 * Defines general operator-specific OSC commands.  /pig/op/name/command
*/

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
)





func AddOpHandler(op Operator, command string, handler func(*goosc.Message)([]string, error)) {
	//name := op.Name()
	server := op.Server()
	//address := fmt.Sprintf("/%s/op/%s/%s", server.Root(), name, command)
	address := op.FormatOSCAddress(command)
	fmt.Printf("DEBUG address is '%s'\n", address)
	var wrappedHandler = func(msg *goosc.Message) {
		status, err := handler(msg)
		if err != nil {
			server.GetResponder().Error(address, status, err)
			server.GetREPLResponder().Error(address, status, err)
		} else {
			server.GetResponder().Ack(address, status)
			server.GetREPLResponder().Ack(address, status)
		}
	}
	server.AddMsgHandler(address, wrappedHandler)
}
	
func (op *baseOperator) remotePing(msg *goosc.Message)([] string, error) {
	var err error
	fmt.Printf("op %s ping\n", op.Name())
	return empty, err
}
	
