package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)


// baseXformOperator extends baseOperator to implement midi.Transform
//
type baseXformOperator struct {
	baseOperator
	xformTable midi.DataTable
}


func (op *baseXformOperator) Reset() {
	base := &op.baseOperator
	xform := &op.xformTable
	base.Reset()
	xform.Reset()
}

func (xop *baseXformOperator) Length() int {
	return xop.xformTable.Length()
}

func (xop *baseXformOperator) Value(index byte) (value byte, err error) {
	return xop.xformTable.Value(index)
}

func (xop *baseXformOperator) SetValue(index byte, value byte) error {
	return xop.xformTable.SetValue(index, value)
}

func (xop *baseXformOperator) Dump() string {
	return xop.xformTable.Dump()
}
	
	
func initXformOperator(xop *baseXformOperator) {

	// cmd op name, q-xform-length
	// osc /pig/op name, q-xform-length
	// Returns:
	//    Transform table length
	//
	remoteQueryLength := func(msg *goosc.Message)([]string, error) {
		var err error
		s := fmt.Sprintf("%d", xop.Length())
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
		value, err = xop.Value(index)
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
		err = xop.SetValue(index, value)
		return empty, err
	}


	// cmd op name, print-xform-table
	// osc /pig/op name, print-xform-table
	//
	// Display hex-dump of transformation table.
	//
	remoteDumpTable := func(msg *goosc.Message)([]string, error) {
		var err error
		fmt.Printf("%s\n", xop.Dump())
		return empty, err
	}
		
	
	xop.addCommandHandler("q-xform-length", remoteQueryLength)
	xop.addCommandHandler("q-xform-value", remoteQueryValue)
	xop.addCommandHandler("set-xform-value", remoteSetValue)
	xop.addCommandHandler("print-xform-table", remoteDumpTable)
	
}	
	
