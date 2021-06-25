package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/osc"
)

var (
	empty []string
)

// Add general op-related handlers to global OSC server
//
func Init() {
	fmt.Println("opserver.Init() executing -- REMOVE THIS LIEN")
	server := osc.GlobalServer
	osc.AddOSCHandler(server, "new", remoteNewOperator)
	// new            : type name <optional ...   -> name
	// remove         : name
        // q-op-types     : -> list
	// q-forest       : -> string
	// destroy-forest 
}



// osc /pig/new optype name
// -> name
//
func remoteNewOperator(msg *goosc.Message)([]string, error) {
	template := []osc.ExpectType{osc.XpString, osc.XpString}
	args, err := osc.Expect(template, msg.Arguments)
	if err != nil {
		return empty, err
	}
	otype, name := args[0], args[1]
	op, err := NewOperator(otype, name)
	if err != nil {
		return empty, err
	}
	return osc.StringSlice(op.Name()), err
}
	
	
