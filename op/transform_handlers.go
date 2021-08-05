package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)


// initiTransformHandlers adds OSC handlers for Operators implementing midi.Transform.
//
// ISSUE: function signature could be cleaner due to being a newbi at Go.
// Both op and xform arguments are the same object.  Is there a way to 'cast' a
// general Operator to midi.Transform ?
//
func initTransformHandlers(op Operator, xform midi.Transform) {


	// cmd op name, q-xform-length
	// osc /pig/op name, q-xform-length
	// Returns:
	//    Transform table length
	//
	remoteQueryLength := func(msg *goosc.Message)([]string, error) {
		var err error
		s := fmt.Sprintf("%d", xform.Length())
		return []string{s}, err
	}

	// cmd op name, q-xform-value, index
	// osc /pig/op name, q-xform-value, index
	// Returns
	//    Indexed value from transform table
	//    Error if index is out of bounds
	//
	remoteQueryValue := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osi", msg)
		if err != nil {
			return empty, err
		}
		var index byte = byte(args[2].I) 
		var value byte
		value, err = xform.Value(index)
		s := fmt.Sprintf("0x%02X", value)
		return []string{s}, err
	}

	// cmd op name, set-xform-value, index, value
	// osc /pig/op name, set-xform-value, index, value
	// Returns error if either index or value are out of bounds.
	//    0 <= index < xform.Length()
	//    0 <= value < 0x80
	//
	remoteSetValue := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osii", msg)
		if err != nil {
			return empty, nil
		}
		var index byte = byte(args[2].I)
		var value byte = byte(args[3].I)
		err = xform.SetValue(index, value)
		return empty, err
	}


	// cmd op name, print-xform-table
	// osc /pig/op name, print-xform-table
	//
	// Display hex-dump of transformation table.
	//
	remoteDumpTable := func(msg *goosc.Message)([]string, error) {
		var err error
		fmt.Printf("%s\n", xform.Dump())
		return empty, err
	}
		
	
	op.addCommandHandler("q-xform-length", remoteQueryLength)
	op.addCommandHandler("q-xform-value", remoteQueryValue)
	op.addCommandHandler("set-xform-value", remoteSetValue)
	op.addCommandHandler("print-xform-table", remoteDumpTable)
	
	
}
