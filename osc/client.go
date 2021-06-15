package osc

import (
	"fmt"
	"strings"
	"os"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)

var (
	globalClient *PigClient
)


func init() { 
	host := config.GlobalParameters.OSCClientHost
	port := int(config.GlobalParameters.OSCClientPort)
	root := config.GlobalParameters.OSCClientRoot
	filename := config.GlobalParameters.OSCClientFilename
	globalClient = NewClient(host, port, root, filename)
	
}


type PigClient struct {
	backing *goosc.Client
	root string
	filename string
	verbose bool
	
}


func NewClient(ip string, port int, root string, filename string) *PigClient {
	client := PigClient{goosc.NewClient(ip, port), root, filename, true}
	return &client
}


func (c *PigClient) IP() string {
	return c.backing.IP()
}


func (c *PigClient) Port() int {
	return c.backing.Port()
}
	

func (c *PigClient) echo(address string, payload string) {
	if c.verbose {
		fmt.Printf("OSC Client  : %s\n", c.root)
		fmt.Printf("response to : %s\n", address)
		for _, s := range strings.Split(payload, "\n") {
			fmt.Printf("            : %s\n", s)
		}
	}
}


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


func AckGlobal(sourceAddress string, payload []string) {
	globalClient.Ack(sourceAddress, payload)
}


func ErrorGlobal(sourceAddress string, payload []string) {
	globalClient.Error(sourceAddress, payload)
}
