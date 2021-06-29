package op


import (
	//"errors"
	"fmt"
	"strconv"
	//"strings"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)




func StringSlice(values ...interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}

func ToStringSlice(values []interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}


type ExpectValue struct {
	S string
	I int64
	F float64
	B bool
	C midi.MIDIChannel
	O Operator}

 
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
			
			
func ExpectMsg(template string, msg *goosc.Message)([]ExpectValue, error) {
	return Expect(template, msg.Arguments)
}
