package osc


import (
	"errors"
	"fmt"
	"strconv"
)

type expectType int

const (
	xpString expectType = iota
	xpInt
	xpFloat
)

func (x expectType) String() string {
	return [...]string{"Any", "String", "Int", "Float"}[x]
}


func expect(template []expectType, arguments []interface{}) ([]string, error) {
	var err error
	acc := make([]string, len(arguments))
	if len(template) > len(arguments) {
		msg := "Expected at least %d arguments, got %d"
		err = errors.New(fmt.Sprintf(msg, len(template), len(arguments)))
	} else {
		for i, xtype := range template {
			s := fmt.Sprintf("%v", arguments[i])
			switch xtype {
			case xpString:
				acc[i] = s
			case xpInt:
				_, err2 := strconv.Atoi(s)
				if err2 != nil {
					msg := "Expected int at index %d, got %s"
					err = errors.New(fmt.Sprintf(msg, i, s))
					break
				} else {
					acc[i] = s
				}
			case xpFloat:
				_, err2 := strconv.ParseFloat(s, 64)
				if err2 != nil {
					msg := "Expected float at index %d, got %s"
					err = errors.New(fmt.Sprintf(msg, i, s))
					break
				} else {
					acc[i] = s
					}
			default:
				acc[i] = s
			} 
		} 
	}
	return acc, err
}


func expectLength(address string, args []string, index int) bool {
	if index < len(args) {
		return true
	} else {
		msg := "ERROR OSC message %s, expected at least %d arguments, got %v\n"
		fmt.Printf(msg, address, index+1, args)
		return false
	}
}


func getStringValue(address string, args []string, index int, fallback string) string {
	var s string
	if expectLength(address, args, index) {
		s = args[index]
	} else {
		s = fallback
	}
	return s
}


func getIntValue(address string, args []string, index int, fallback int64) int64 {
	var n int64
	var err error
	if expectLength(address, args, index) {
		s := args[index]
		n, err = strconv.ParseInt(s, 0, 64)
		if err != nil {
			msg := "ERROR OSC message %s, expected int at index %d, got %s\n"
			fmt.Println(msg, address, index, s)
			n = fallback
		}
	}
	return n
}
		
func getFloatValue(address string, args []string, index int, fallback float64) float64 {
	var n float64
	var err error
	if expectLength(address, args, index) {
		s := args[index]
		n, err = strconv.ParseFloat(s, 64)
		if err != nil {
			msg := "ERROR OSC message %s, expected float at index %d, got %s\n"
			fmt.Println(msg, address, index, s)
			n = fallback
		}
	}
	return n
}


func toSlice(values ...interface{}) []string {
	acc := make([]string, len(values))
	for i, v := range values {
		acc[i] = fmt.Sprintf("%v", v)
	}
	return acc
}


	
