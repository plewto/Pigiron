package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/piglog"	
)


// PigServer interface defines server-side OSC interface.
//
// Root() string
//    Returns the OSC address prefix, default is '/pig'
//
// SetRoot(s string)
//    Sets OSC server address prefix.
//
// GetResponder() Responder
//    Returns the OSC responder.
//
// GetREPLResponder() Responder
//    Returns the responder used for printing to the terminal.
//
// AddMsgHandler(address string, handler func(msg *go-osc.Message))
//    Adds an OSC handler function.
//    
// ListenAndServe()
//    Start server.
//
// IP() string
//    Returns server IP address.
//
// Port() int
//    Returns server port number.
//
// Close()
//    Close the server.  Close should only be called on termination of the
//    program.
//
// Commands() []string
//    Returns list of defined OSC commands.
//
type PigServer interface {
	Root() string
	SetRoot(string)
	GetResponder() Responder
	GetREPLResponder() Responder
	AddMsgHandler(address string, handler func(msg *goosc.Message))
	ListenAndServe()
	IP() string
	Port() int
	Close()
	Commands() []string
}


// OSCServer struct implements the PigServer interface.
//
type OSCServer struct {
	backingServer *goosc.Server
	dispatcher *goosc.StandardDispatcher
	root string
	responder Responder
	replResponder Responder
	ip string
	port int
	commands []string
}


// NewServer() returns a new PigServer.
//
func NewServer(ip string, port int, root string) PigServer {
	server := new(OSCServer)
	server.ip = ip
	server.port = port
	server.root = root
	server.responder = globalResponder
	server.replResponder = replResponder
	addr := fmt.Sprintf("%s:%d", ip, port)
	server.dispatcher = goosc.NewStandardDispatcher()
	server.backingServer = &goosc.Server {
		Addr: addr,
		Dispatcher: server.dispatcher,
	}
	server.commands = make([]string, 0, 16)
	return server
}

func (s *OSCServer) GetResponder() Responder {
	return s.responder
}

func (s *OSCServer) GetREPLResponder() Responder {
	return s.replResponder
}

func (s *OSCServer) AddMsgHandler(address string, handler func(msg *goosc.Message)) {

	logger := func(msg *goosc.Message) {
		piglog.Print(msg.Address)
		for i, a := range msg.Arguments {
			piglog.Print(fmt.Sprintf("[%2d] %v", i, a))
		}
		handler(msg)
	}
	s.dispatcher.AddMsgHandler(address, logger)
	s.commands = append(s.commands, address)
}

func (s *OSCServer) Commands() []string {
	return s.commands
}

func (s *OSCServer) ListenAndServe() {
	fmt.Printf("OSC Listening: %s:%d  /%s\n", s.ip, s.port, s.root)
	go s.backingServer.ListenAndServe()
}

func (s *OSCServer) Root() string {
	return s.root
}

func (s *OSCServer) SetRoot(root string) {
	s.root = root
}

	
func (s *OSCServer) IP() string {
	return s.ip
}

func (s *OSCServer) Port() int {
	return s.port
}

func (s *OSCServer) Close() {
	s.backingServer.CloseConnection()
}

// AddHandler()  adds new OSC handler function to server s.
// The OSC address prefix is automatically added to the command argument.
// The command "foo" --> becomes "/pig/foo"
// 
func AddHandler(s PigServer, command string, handler func(*goosc.Message)([]string, error)) {
	commands[command] = true
	address := fmt.Sprintf("/%s/%s", s.Root(), command)
	var result = func(msg *goosc.Message) {
		status, err := handler(msg)
		if err != nil {
			s.GetResponder().Error(address, status, err)
			s.GetREPLResponder().Error(address, status, err)
		} else {
			s.GetResponder().Ack(address, status)
			s.GetREPLResponder().Ack(address, status)
		}
	}
	s.AddMsgHandler(address, result)
}
