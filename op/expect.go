package op

/*
 * Expect provide OSC message argument validation.
 *
*/

import (
	"fmt"
	"strconv"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)


// StringSlice Embeds all arguments into string slice.
//
func StringSlice(values ...interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}


// ToStringSlice converts []interface slice into a string slice.
//
func ToStringSlice(values []interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}


// ExpectValue struct holds return values for the Expect function.
// Each field has a different type and the Expect function sets the appropriate
// field based on it's template argument.  Only one field should ever have
// a non-nil value.
//
type ExpectValue struct {
	S string
	I int64
	F float64
	B bool
	C midi.MIDIChannel
	O Operator}


// Expect validates a list of values for appropriate type.
// Each character in the template indicates the expected value for the
// corresponding position of the values list.
// The possible template characters are:
//     s - string
//     i - int64
//     f - float64
//     b - bool (see strconv.ParseBool for excepted values)
//     c - MIDI channel (int 1 <= n <= 16)
//     o - Operator name
//
// The returned error is non-nil if either:
//    - len(template) > len(values)
//    - values[i] does not correspond to template[i]
//
// If no errors are detected the result is a list of ExpectValue of
// length len(template), with the appropriate ExpectType field set
// for each element.
//
func Expect(template string, values []interface{})([]ExpectValue, error) {
	var err error
	var acc []ExpectValue = make([]ExpectValue, len(template))
	if len(template) > len(values) {
		msg := "Expected at least %d arguments, got %d"
		err = fmt.Errorf(msg, len(template), len(values))
		return acc, err
	}
	for i, xtype := range template {
		arg := values[i]
		switch xtype {
		case 's':
			acc[i].S = arg.(string)
		case 'i':
			var s string = fmt.Sprintf("%d", arg)
			var n int64 = 0
			n, err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				msg := "Expected int at index %d, got %v"
				err = fmt.Errorf(msg, i, arg)
				return acc, err
			}
			acc[i].I = n
		case 'f':
			var s string = fmt.Sprintf("%f", arg)
			var n float64 = 0.0
			n, err = strconv.ParseFloat(s, 64)
			if err != nil {
				msg := "Expected float at index %d, got %v"
				err = fmt.Errorf(msg, i, arg)
				return acc, err
			}
			acc[i].F = n
		case 'b':
			var s string = fmt.Sprintf("%s", arg)
			var v bool = false
			v, err = strconv.ParseBool(s)
			if err != nil {
				msg := "Expected bool at index %d, got %s"
				err = fmt.Errorf(msg, i, s)
				return acc, err
			}
			acc[i].B = v
		case 'c':
			var s string = fmt.Sprintf("%d", arg)
			var n int64 = 0
			n, err = strconv.ParseInt(s, 10, 64)
			if err != nil || n < 1 || 16 < n {
				msg := "Expected MIDI channel at index %d, got %v"
				err = fmt.Errorf(msg, i, arg)
				return acc, err
			}
			acc[i].C = midi.MIDIChannel(n)
		case 'o':
			var s string = fmt.Sprintf("%s", arg)
			var op Operator
			op, err = GetOperator(s)
			if err != nil {
				msg := "Expected Operator name at index %d, got %s"
				err = fmt.Errorf(msg, i, arg)
				return acc, err
			}
			acc[i].O = op
		default:
			msg := "Unknown Expect template type '%s'"
			err = fmt.Errorf(msg, xtype)
			panic(err)
		}
	}
	return acc, err
}
			
// ExpectMsg is identical to Expect but is applied to osc message Arguments field.
//
func ExpectMsg(template string, msg *goosc.Message)([]ExpectValue, error) {
	return Expect(template, msg.Arguments)
}
