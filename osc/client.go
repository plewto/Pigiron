package osc

// PigClient provides for OSC return messages.
// There are 2 possible response messages:
//     1) ACK (Acknowledge)
//     2) ERROR
// with identical format except for the message address.
//
// If the filename field is a non-empty string, a duplicate of the 
// response is written to the file.   The response file is overwritten after
// each transmission.   No error message is reported if the response file
// can not be opened.
//
// There may be more then one PigClient but a single global client is
// provided as a default.
//

import (
	"fmt"
	//"strings"
	"os"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)


// PigClient interface defines OSC server return messages.
//
type PigClient interface {
	Send(*goosc.Message)
	Ack(address string, data []string)
	Error(address string, data []string)
	Info() string
	ForREPL() bool
	SetForREPL(flag bool)
}



// BasicClient is the default implementation for PigClient.
//
// In response to an incoming OSC message, a PigServer sends either an 
// ACK (OK), or ERROR message back to the client.
// If the filename field is a non-empty string the response is written to a
// temporary file. 
//
type BasicClient struct {
	backing *goosc.Client
	root string
	filename string
	forREPL bool
}

// NewClient creates a new instance of PigClient.
//
// ip - client host ip address
// port - client port number
// root - OSC address prefix.  For root 'foo' and command 'bar' the
// ultimate OSC address is '/foo/bar'
// filename - If non-empty a textural version of each OSC response message
// is written to a temporary file.   If the file can not be opened, it is
// silently ignored.
//
func NewClient(ip string, port int, root string, filename string) PigClient {
	client := BasicClient{goosc.NewClient(ip, port), root, filename, true}
	return &client
}

// IP returns the client's host IP address.
//
func (c *BasicClient) IP() string {
	return c.backing.IP()
}

// Port returns the client's port number.
//
func (c *BasicClient) Port() int {
	return c.backing.Port()
}
	
func (c *BasicClient) ForREPL() bool {
	return c.forREPL
}

func (c *BasicClient) SetForREPL(flag bool) {
	c.forREPL = flag
}



// writeResponseFile creates a file for the most recently transmitted message.
// If the filename field is empty or it can not be created, the write is
// silently ignored.
//
func (c *BasicClient) writeResponseFile(address string, payload string) {
	if c.ForREPL() {
		return
	}
	if len(c.filename) > 0 {
		file, err := os.Create(c.filename)
		if err == nil {
			defer file.Close()
			file.WriteString(fmt.Sprintf("%s\n", address))
			file.WriteString(payload)
		}
	}
}


func (c *BasicClient) Send(msg *goosc.Message) {
	if !c.ForREPL() {
		c.backing.Send(msg)
	}
}

// Ack transmits an 'Acknowledgment' message.
// The message is transmitted via OSC and saved to a temporary response
// file.
//
// sourceAddress - the OSC address this is an acknowledgment of.
// payload - optional values included in the response.
//
func (c *BasicClient) Ack(sourceAddress string, payload []string) {
	address := fmt.Sprintf("/%s/ACK", c.root)
	if c.ForREPL() && !inBatchMode {
		fmt.Printf("------------------------ ACK %s\n", address)
		for i, p := range payload {
			fmt.Printf("\t[%2d] %s\n", i, p)
		}
	} else {
		msg := goosc.NewMessage(address)
		msg.Append(sourceAddress)
		acc := fmt.Sprintf("ACK\n%s\n", sourceAddress)
		for _, s := range payload {
			msg.Append(s)
			acc += fmt.Sprintf("%s\n", s)
		}
		c.backing.Send(msg)
		c.writeResponseFile(address, acc)
	}
}

// Error transmits an 'Error' message.
// With exception of the OSC message address, Error is identical to Ack.
//
func (c *BasicClient) Error(sourceAddress string, payload []string) {
	address := fmt.Sprintf("/%s/ERROR", c.root)
	if c.ForREPL() {
		fmt.Print(config.GlobalParameters.ErrorColor)
		fmt.Printf("------------------------ ERROR %s\n", address)
		for i, p := range payload {
			fmt.Printf("\t[%2d] %s\n", i, p)
		}
	} else {
		msg := goosc.NewMessage(address)
		msg.Append(sourceAddress)
		acc := fmt.Sprintf("ERROR\n%s\n", sourceAddress)
		for _, s := range payload {
			msg.Append(s)
			acc += fmt.Sprintf("%s\n", s)
		}
		c.backing.Send(msg)
		c.writeResponseFile(address, acc)
	}
	batchError = true
}


func (c *BasicClient) Info() string {
	acc := "BasicClient\n"
	acc += fmt.Sprintf("\troot     : \"%s\"\n", c.root)
	acc += fmt.Sprintf("\tfilename : \"%s\"\n", c.filename)
	// acc += fmt.Sprintf("\tverbose  : %v\n", c.verbose)
	return acc
}
	
