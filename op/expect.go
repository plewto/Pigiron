package op


import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	goosc "github.com/hypebeast/go-osc/osc"
)


var boolValues = map[string]string{
	"false" : "false",
	"0" : "false",
	"f" : "false",
	"off" : "false",
	"no" : "false",
	"n"  : "false",
	"disable" : "false",
	"true" : "true",
	"1" : "true",
	"t" : "true",
	"on" : "true",
	"yes" : "true",
	"y" : "true",
	"enable" : "true"}
	
func parseBool(s string)(string, error) {
	var err error
	v, flag := boolValues[strings.ToLower(s)]
	if !flag {
		msg := "Expected boolean, got %s"
		err = errors.New(fmt.Sprintf(msg, s))
		return "false", err
	} else {
		return v, err
	}
}


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


// template s -> string
//          i -> int
//          f -> float
//          b -> bool
//          c -> midi channel (1..16)
//          * -> any   (convert to string)
//
func Expect(template string, arguments []string)([]string, error) {
	var err error
	acc := make([]string, len(arguments))
	// Bug 001 fix
	// Removes spurious first element from arguments.
	if len(arguments) == 1 && arguments[0] == "" {
		arguments = arguments[1:]
	}

	if len(template) > len(arguments) {
		msg := "Expected at least %d arguments, got %d"
		err = errors.New(fmt.Sprintf(msg, len(template), len(arguments)))
		return empty, err
	}

	for i, xtype := range template {
		arg := arguments[i]
		switch xtype {
		case 's':
			acc[i] = arg
		case 'i':
			_, err := strconv.Atoi(arg)
			if err != nil {
				msg := "Expected int at index %d, got %s"
				err = errors.New(fmt.Sprintf(msg, i, arg))
				return empty, err
			}
			acc[i] = arg
		case 'f':
			_, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				msg := "Expecrted float at index %d, got %s"
				err = errors.New(fmt.Sprintf(msg, i, arg))
				return empty, err
			}
			acc[i] = arg
		case 'b':
			v, err := parseBool(arg)
			if err != nil {
				return empty, err
			}
			acc[i] = v
		case 'c':
			v, err := strconv.Atoi(arg)
			if err != nil {
				msg := "Expected MIDI channel at index %d, got %s"
				err = errors.New(fmt.Sprintf(msg, i, arg))
				return empty, err
			}
			if v < 1 || 16 < v {
				msg := "Expected MIDI channel at index %d, got %s"
				err = errors.New(fmt.Sprintf(msg, i, arg))
				return empty, err
			}
			acc[i] = arg
					
		case '*':
			acc[i] = fmt.Sprintf("%v", arg)
		default:
			acc[i] = arg
		}
	}
	return acc, err
}
		

func ExpectMsg(template string, msg *goosc.Message)([]string, error) {
	return Expect(template, ToStringSlice(msg.Arguments))
}



func ExpectLength(address string, args []string, index int) bool {
	if index < len(args) {
		return true
	} else {
		msg := "ERROR OSC message %s, Expected at least %d arguments, got %v\n"
		fmt.Printf(msg, address, index+1, args)
		return false
	}
}


func GetStringValue(address string, args []string, index int, fallback string) string {
	var s string
	if ExpectLength(address, args, index) {
		s = args[index]
	} else {
		s = fallback
	}
	return s
}


func GetIntValue(address string, args []string, index int, fallback int64) int64 {
	var n int64
	var err error
	if ExpectLength(address, args, index) {
		s := args[index]
		n, err = strconv.ParseInt(s, 0, 64)
		if err != nil {
			msg := "ERROR OSC message %s, Expected int at index %d, got %s\n"
			fmt.Println(msg, address, index, s)
			n = fallback
		}
	}
	return n
}
		
func GetFloatValue(address string, args []string, index int, fallback float64) float64 {
	var n float64
	var err error
	if ExpectLength(address, args, index) {
		s := args[index]
		n, err = strconv.ParseFloat(s, 64)
		if err != nil {
			msg := "ERROR OSC message %s, Expected float at index %d, got %s\n"
			fmt.Println(msg, address, index, s)
			n = fallback
		}
	}
	return n
}





	
