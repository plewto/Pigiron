package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)

/*
** baseXformOperator extends baseOperator to implement midi.Transform
**
*/
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

func (xop *baseXformOperator) TransformRange() (floor byte, ceiling byte) {
	return xop.xformTable.TransformRange()
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
	
func (xop *baseXformOperator) Plot() string {
	return xop.xformTable.Plot()
}


func initXformOperator(xop *baseXformOperator) {

	// cmd op name, q-table-range
	// osc /pig/op name, q-table-range
	// Returns table index range:
	//    [floor, ceiling]
	//
	remoteQueryRange := func(msg *goosc.Message)([]string, error) {
		var err error
		f, c := xop.TransformRange()
		sf, sc := fmt.Sprintf("0x%02X", f), fmt.Sprintf("0x%02X", c)
		return []string{sf, sc}, err
	}

	// cmd op name, q-table-value, index
	// osc /pig/op name, q-table-value, index
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

	// cmd op name, set-table-value, index, value, [values ...]
	// osc /pig/op name, set-table-value, index, value, [values ...]
	// Returns error if either index or value are out of bounds.
	//    0 <= index < ceiling
	//    0 <= value < 0x80
	//
	remoteSetValue := func(msg *goosc.Message)([]string, error) {
		template := "osii"
		for i := 4; i < len(msg.Arguments); i++ {
			template += "i"
		}
		args, err := ExpectMsg(template, msg)
		if err != nil {
			return empty, nil
		}
		var index byte = byte(args[2].I)
		count := len(template) - 3
		for i := 0; i < count; i++ {
			value := byte(args[i+3].I)
			err = xop.xformTable.SetValue(index, value)
			if err != nil {
				break
			}
			index++
		}		
		return empty, err
	}


	// cmd op name, print-table
	// osc /pig/op name, print-table
	//
	// Display hex-dump of transformation table.
	//
	remoteDumpTable := func(msg *goosc.Message)([]string, error) {
		var err error
		fmt.Printf("%s\n", xop.Dump())
		fmt.Printf("%s", xop.Plot())
		return empty, err
	}
	
	xop.addCommandHandler("q-table-range", remoteQueryRange)
	xop.addCommandHandler("q-table-value", remoteQueryValue)
	xop.addCommandHandler("set-table-value", remoteSetValue)
	xop.addCommandHandler("print-table", remoteDumpTable)
	
}	
	
