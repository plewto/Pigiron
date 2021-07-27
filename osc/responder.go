package osc

import (
	"fmt"
	"os"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
	"github.com/plewto/pigiron/piglog"
	
)

// Responder interface defines how the OSC server sends responses back to the client.
//
// There are two possible responses: Ack and Error.
//
// Ack() sends an Acknowledgment that an OSC message has been received
// without error.  As part of it's "payload" it includes the original
// message and any requested data.
//
// An Error() response indicates the last received OSC message caused an
// error.  The response includes the offending OSC message and an additional
// error message. 
// 
type Responder interface {
	Ack(sourceAddress string, args []string)
	Error(sourceAddress string, args []string, err error)
	String() string
}

// BasciResponder is the primary implementation of the Responder interface.
// In addition to transmitting ACk and Error responses via OSC, it also
// writes identical information to a temporary file.  This is useful for
// clients which do not receive OSC.   The file is overwritten each time a
// new OSC message is received.
//
type BasicResponder struct {
	client *goosc.Client
	root string
	filename string
}


// NewBasicResponder() creates a new instance of basicResponder.
//
func NewBasicResponder(ip string, port int, root string, filename string) Responder {
	client := goosc.NewClient(ip, port)
	responder := BasicResponder{client, root, filename}
	return &responder
}

// r.writeResponseFile() creates a file for the most recently transmitted message.
// If the filename field is empty or it can not be created, the write is
// silently ignored.
//
func (r *BasicResponder) writeResponseFile(sourceAddress string, payload string) {
	if len(r.filename) > 0 {
		file, err := os.Create(r.filename)
		if err == nil {
			defer file.Close()
			file.WriteString(fmt.Sprintf("%s\n", sourceAddress))
			file.WriteString(payload)
		}
	}
}

// r.send() transmits OSC message to client.
//
func (r *BasicResponder) send(msg *goosc.Message) {
	r.client.Send(msg)
}

// f.Ack() transmits an Acknowledgment response to the client.
//
func (r *BasicResponder) Ack(sourceAddress string, args []string) {
	address := fmt.Sprintf("/%s/ACK", r.root)
	acc := fmt.Sprintf("ACK\n%s\n", sourceAddress)
	msg := goosc.NewMessage(address)
	msg.Append(sourceAddress)
	for i, a := range args {
		s := fmt.Sprintf("%s", a)
		msg.Append(s)
		acc += fmt.Sprintf("%s\n", s)
		piglog.Log(fmt.Sprintf("-> ACK [%3d] %s", i, a))
		
	}
	r.send(msg)
	r.writeResponseFile(sourceAddress, acc)
}

// f.Error() transmits an Error response to the client.
//
func (r *BasicResponder) Error(sourceAddress string, args []string, err error) {
	address := fmt.Sprintf("/%s/ERROR", r.root)
	acc := fmt.Sprintf("ERROR\n%s\n", sourceAddress)
	msg := goosc.NewMessage(address)
	msg.Append(sourceAddress)
	msg.Append(fmt.Sprintf("%s\n", err))
	acc += fmt.Sprintf("%s\n", err)
	piglog.Log(fmt.Sprintf("-> %s", acc))
	for i, a := range args {
		s := fmt.Sprintf("%s", a)
		msg.Append(s)
		acc += fmt.Sprintf("%s\n", s)
		piglog.Log(fmt.Sprintf("-> ERR [%3d] %s", i, a))
	}
	r.send(msg)
	r.writeResponseFile(sourceAddress, acc)
}

func (r *BasicResponder) String() string {
	host := r.client.IP()
	port := r.client.Port()
	root := r.root
	filename := r.filename
	acc := "BasicResponder "
	acc += fmt.Sprintf("root: %s,  IP %s:%d, filename '%s'", root, host, port, filename)
	return acc
}
	

// REPLResponder struct is a Responder which prints messages to the terminal.
//
type REPLResponder struct {}

// NewREPLResponder() creates new instance of REPLResponder.
//
func NewREPLResponder() Responder {
	return &REPLResponder{}
}

func setTextColor() {
	fmt.Print(config.GlobalParameters.TextColor)
}

func setErrorColor() {
	fmt.Print(config.GlobalParameters.ErrorColor)
}

func bar(text string) {
	fmt.Printf("\n-------------------------------- %s\n", text)
}

func (r *REPLResponder) Ack(sourceAddress string, args []string) {
	batchError = false
	if !inBatchMode {
		setTextColor()
		bar("OK")
		fmt.Println(sourceAddress)
		for i, a := range args {
			fmt.Printf("\t[%2d] %s\n", i, a)
		}
		Prompt()
	}
}


func (r *REPLResponder) Error(sourceAddress string, args []string, err error) {
	setErrorColor()
	bar("ERROR")
	fmt.Println(sourceAddress)
	fmt.Printf("%s\n", err)
	for i, a := range args {
		fmt.Printf("\t[%2d] %s\n", i, a)
	}
	setTextColor()
	Prompt()
	batchError = true
}


func (r *REPLResponder) String() string {
	return "REPLResponder"
}
	
