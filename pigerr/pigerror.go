package pigerr

import (
	"fmt"
)

// Error returns new error object with given message.
// 
func New(message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	return fmt.Errorf(msg)
}

// CompoundError creates new error by composing previous error with new message.
//
func CompoundError(err1 error, message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s\n", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	msg += fmt.Sprintf("PREVIOUS ERROR:\n%s", err1)
	return fmt.Errorf(msg)
}


type warning interface {
	String() string
}

type baseWarning struct {
	message string
}

func (w *baseWarning) String() string {
	return w.message
}


func Warning(message string, more ...string) warning {
	msg := fmt.Sprintf("\nWARNING: %s", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nWARNING: %s", s)
	}
	warn := baseWarning{msg}
	fmt.Printf("%s\n", warn)
	return &warn
}


