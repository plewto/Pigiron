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
	"strings"
	"os"
	goosc "github.com/hypebeast/go-osc/osc"
	//"github.com/plewto/pigiron/config"
)

// var (
// 	globalClient *PigClient
// )


// func init() { 
// 	host := config.GlobalParameters.OSCClientHost
// 	port := int(config.GlobalParameters.OSCClientPort)
// 	root := config.GlobalParameters.OSCClientRoot
// 	filename := config.GlobalParameters.OSCClientFilename
// 	globalClient = NewClient(host, port, root, filename)
	
// }

// PigClient provides an OSC callback client.
//
// In response to an incoming OSC message, a PigServer sends either an 
// ACK (OK), or ERROR message back to the client.
// If the filename field is a non-empty string the response is written to a
// temporary file. 
//
type PigClient struct {
	backing *goosc.Client
	root string
	filename string
	verbose bool
	
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
func NewClient(ip string, port int, root string, filename string) *PigClient {
	client := PigClient{goosc.NewClient(ip, port), root, filename, true}
	return &client
}

// IP returns the client's host IP address.
//
func (c *PigClient) IP() string {
	return c.backing.IP()
}

// Port returns the client's port number.
//
func (c *PigClient) Port() int {
	return c.backing.Port()
}
	

// echo prints each transmitted OSC message to the terminal.
// Does nothing if verbose is false.
//
func (c *PigClient) echo(address string, payload string) {
	if c.verbose {
		fmt.Printf("OSC Client  : %s\n", c.root)
		fmt.Printf("response to : %s\n", address)
		for _, s := range strings.Split(payload, "\n") {
			fmt.Printf("            : %s\n", s)
		}
	}
}


// writeResponseFile creates a file for the most recently transmitted message.
// If the filename field is empty or it can not be created, the write is
// silently ignored.
//
func (c *PigClient) writeResponseFile(address string, payload string) {
	if len(c.filename) > 0 {
		file, err := os.Create(c.filename)
		if err == nil {
			defer file.Close()
			file.WriteString(fmt.Sprintf("%s\n", address))
			file.WriteString(payload)
		}
	}
}


// Ack transmits an 'Acknowledgment' message.
// The message is transmitted via OSC and saved to a temporary response
// file.
//
// sourceAddress - the OSC address this is an acknowledgment of.
// payload - optional values included in the response.
//
func (c *PigClient) Ack(sourceAddress string, payload []string) {
	address := fmt.Sprintf("/%s/ACK", c.root)
	msg := goosc.NewMessage(address)
	msg.Append(sourceAddress)
	acc := fmt.Sprintf("ACK\n%s\n", sourceAddress)
	for _, s := range payload {
		msg.Append(s)
		acc += fmt.Sprintf("%s\n", s)
	}
	c.backing.Send(msg)
	c.writeResponseFile(address, acc)
	c.echo(address, acc)
}

// Error transmits an 'Error' message.
// With exception of the OSC message address, Error is identical to Ack.
//
func (c *PigClient) Error(sourceAddress string, payload []string) {
	address := fmt.Sprintf("/%s/ERROR", c.root)
	msg := goosc.NewMessage(address)
	msg.Append(sourceAddress)
	acc := fmt.Sprintf("ERROR\n%s\n", sourceAddress)
	for _, s := range payload {
		msg.Append(s)
		acc += fmt.Sprintf("%s\n", s)
	}
	c.backing.Send(msg)
	c.writeResponseFile(address, acc)
	c.echo(address, acc)
}

// AckGlobal transmits an 'Acknowledgment' response to the global OSC client.
// It has identical usage to the PigClient.Ack method.
//
func AckGlobal(sourceAddress string, payload []string) {
	globalClient.Ack(sourceAddress, payload)
}

// ErrorGlobal transmits an Error response to the global OSC client.
// It has identical usage to the PigClient.Error method.
//
func ErrorGlobal(sourceAddress string, payload []string) {
	globalClient.Error(sourceAddress, payload)
}
