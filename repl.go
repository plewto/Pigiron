package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)

var prompt = "PIG: "

func parseAddress(s string) (command string, args string) {
	pos := strings.Index(s, " ")

	if pos == -1 {
		command = strings.TrimSpace(s)
		args = ""
	} else {
		command = s[:pos]
		args = strings.TrimSpace(s[pos:])
	}
	return command, args
}
	

// repl provides local interactive control.
//
// At the PIG: prompt the user may directly enter OSC commands.
// The format is almost identical to remote commands with the following
// exceptions:
//
//   1) The osc address root is not required, it is automatically supplied
//      i.e.  enter 'exit'  not /pig/exit
//
//   2) Arguments to the command are separated by commas.
//      i.e.  new-operator Monitor, alpha
//
func repl() {
	reader := bufio.NewReader(os.Stdin)
	ip := config.GlobalParameters.OSCServerHost
	port := int(config.GlobalParameters.OSCServerPort)
	root := config.GlobalParameters.OSCServerRoot
	client := goosc.NewClient(ip, port)
	for {
		fmt.Print(prompt)
		raw, _ := reader.ReadString('\n')
		raw = strings.TrimSpace(raw)
		command, args := parseAddress(raw)
		address := fmt.Sprintf("/%s/%s", root, command)
		msg := goosc.NewMessage(address)
		for _, arg := range strings.Split(args, ",") {
			msg.Append(strings.TrimSpace(arg))
		}
		client.Send(msg)
	}
}
