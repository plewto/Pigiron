package osc

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/config"
)

var (
	globalClient *PigClient
	globalServer *PigServer
	empty []string
	Exit bool = false  // exit application if true
)


func init() {
	// Create global OSC client
	host := config.GlobalParameters.OSCClientHost
	port := int(config.GlobalParameters.OSCClientPort)
	root := config.GlobalParameters.OSCClientRoot
	filename := config.GlobalParameters.OSCClientFilename
	globalClient = NewClient(host, port, root, filename)
	// Create global OSC server
	host = config.GlobalParameters.OSCServerHost
	port = int(config.GlobalParameters.OSCServerPort)
	root = config.GlobalParameters.OSCServerRoot
	globalServer = NewServer(host, port, root)
}

type PigServer struct {
	backingServer *goosc.Server
	dispatcher *goosc.StandardDispatcher
	root string
	client *PigClient
}


func NewServer(ip string, port int, root string) *PigServer {
	server := new(PigServer)
	server.root = root
	server.client = globalClient
	addr := fmt.Sprintf("%s:%d", ip, port)
	server.dispatcher = goosc.NewStandardDispatcher()
	server.backingServer = &goosc.Server {
		Addr: addr,
		Dispatcher: server.dispatcher,
	}
	server.AddMsgHandler("ping", server.ping)
	server.AddMsgHandler("exit", server.exit)
	return server
}

func (s *PigServer) AddMsgHandler(command string, handler func(msg *goosc.Message)) {
	addr := fmt.Sprintf("/%s/%s", s.root, command)
	s.dispatcher.AddMsgHandler(addr, handler)
}

func (s *PigServer) ListenAndServe() {
	go s.backingServer.ListenAndServe()
}

func Listen() {
	fmt.Println("OSC Listening....")
	globalServer.ListenAndServe()
}




