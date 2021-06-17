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


func expect(template []expectType, arguments []interface{}) ([]interface{}, error) {
	var err error
	acc := make([]interface{}, len(arguments))
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
				n, err2 := strconv.Atoi(s)
				if err2 != nil {
					msg := "Expected int at index %d, got %s"
					err = errors.New(fmt.Sprintf(msg, i, s))
					break
				} else {
					acc[i] = int64(n)
				}
			case xpFloat:
				n, err2 := strconv.ParseFloat(s, 64)
				if err2 != nil {
					msg := "Expected float at index %d, got %s"
					err = errors.New(fmt.Sprintf(msg, i, s))
					break
				} else {
					acc[i] = float64(n)
				}
			default:
				acc[i] = s
			} // end switch
		} // end for i, xtype
	}

	return acc, err
}
	
	
