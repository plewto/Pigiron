package osc

import (
	"fmt"
	"os"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)

type Responder interface {
	Ack(sourceAddress string, args []string)
	Error(address string, args []string, err error)
	String() string
}

type BasicResponder struct {
	client *goosc.Client
	root string
	filename string
}

func NewBasicResponder(ip string, port int, root string, filename string) Responder {
	client := goosc.NewClient(ip, port)
	responder := BasicResponder{client, root, filename}
	return &responder
}

// writeResponseFile creates a file for the most recently transmitted message.
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

func (r *BasicResponder) send(msg *goosc.Message) {
	r.client.Send(msg)
}

func (r *BasicResponder) Ack(sourceAddress string, args []string) {
	address := fmt.Sprintf("/%s/ACK", r.root)
	acc := fmt.Sprintf("ACK\n%s\n", sourceAddress)
	msg := goosc.NewMessage(address)
	msg.Append(sourceAddress)
	for _, a := range args {
		s := fmt.Sprintf("%s", a)
		msg.Append(s)
		acc += fmt.Sprintf("%s\n", s)
	}
	r.send(msg)
	r.writeResponseFile(sourceAddress, acc)
}

func (r *BasicResponder) Error(sourceAddress string, args []string, err error) {
	address := fmt.Sprintf("/%s/ERROR", r.root)
	acc := fmt.Sprintf("ERROR\n%s\n", sourceAddress)
	msg := goosc.NewMessage(address)
	msg.Append(sourceAddress)
	msg.Append(fmt.Sprintf("%s\n", err))
	acc += fmt.Sprintf("%s\n", err)
	for _, a := range args {
		s := fmt.Sprintf("%s", a)
		msg.Append(s)
		acc += fmt.Sprintf("%s\n", s)
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
	acc += fmt.Sprintf("root: %s,  IP %s:%s, filename '%s'", root, host, port, filename)
	return acc
}
	


type REPLResponder struct {}

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
}


func (r *REPLResponder) String() string {
	return "REPLResponder"
}
	
