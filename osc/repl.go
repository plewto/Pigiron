package osc

import (
	"fmt"
	"strings"
	"os"
	"bufio"
	"time"
	
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
	
)

var (
	internalClient *goosc.Client
	reader *bufio.Reader
	 
)


func init() {
	host := config.GlobalParameters.OSCServerHost
	port := int(config.GlobalParameters.OSCServerPort)
	internalClient = goosc.NewClient(host, port)
	reader = bufio.NewReader(os.Stdin)
	
}
	


func Read() string {
	s, _ := reader.ReadString('\n')
	return s
}


func split(s string) (string, []string) {
	var command string = ""
	var args []string = make([]string, 0)
	words := strings.Split(s, ",")
	if len(words) > 0 {
		command = strings.TrimSpace(words[0])
		if len(words) > 1 {
			for _, w := range words[1:] {
				args = append(args, strings.TrimSpace(w))
			}
		}
	}
	return command, args
}
	
func Eval(command string, args []string) {
	switch command {
	case "":  // ignore blank lines
	case "#": // ignore comment lines
	default:
		root := config.GlobalParameters.OSCServerRoot
		address := fmt.Sprintf("/%s/%s", root, command)
		msg := goosc.NewMessage(address)
		for _, s := range args {
			msg .Append(s)
		}
		internalClient.Send(msg)
	}
}


func REPL() {
	for {
		fmt.Print(config.GlobalParameters.TextColor)
		Prompt()
		raw := Read()
		command, args := split(raw)
		Eval(command, args)
		time.Sleep(10 * time.Millisecond)
	}
}
		
