package osc


import (
	"errors"
	"fmt"
	"strconv"
)


// template s -> string
//          i -> int
//          f -> float
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
			_, err2 := strconv.Atoi(arg)
			if err2 != nil {
				msg := "Expected int at index %d, got %s"
				err = errors.New(fmt.Sprintf(msg, i, arg))
				return empty, err
			}
			acc[i] = arg
		case 'f':
			_, err2 := strconv.ParseFloat(arg, 64)
			if err2 != nil {
				msg := "Expecrted float at index %d, got %s"
				err = errors.New(fmt.Sprintf(msg, i, arg))
				return empty, err
			}
			acc[i] = arg
		default:
			acc[i] = arg
		}
	}
	return acc, err
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





	
