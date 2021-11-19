package backend

import (
	"fmt"
	"strings"
	"strconv"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // auto registers driver
	gomidi "gitlab.com/gomidi/midi/v2"
)


func InputNames() []string {
	acc := make([]string,0, 8)
	ins, _ := gomidi.Ins()
	for _, in := range ins {
		acc = append(acc, in.String())
	}
	return acc
}


func OutputNames() []string {
	acc := make([]string,0, 8)
	outs, _ := gomidi.Outs()
	for _, out := range outs {
		acc = append(acc, out.String())
	}
	return acc
}

func getOutputByIndex(s string) (gomidi.Out, error) {
	var err error
	var index int
	outs, _ := gomidi.Outs()
	index, err = strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	if index < 0 || len(outs) <= index {
		err = fmt.Errorf("MIDI output index out of bounds: %d", index)
		return nil, err
	}
	return outs[index], nil
}

func getOutputByName(pattern string) (gomidi.Out, error) {
	var err error
	outs, _ := gomidi.Outs()
	for _, port := range outs {
		name := port.String()
		if strings.Contains(name, pattern) {
			return port, nil
		}
	}
	err = fmt.Errorf("Pattern '%s' did not match any MIDI device", pattern)
	return nil, err
}
		

// GetMIDIOuput returns indicated MIDI output.
// pattern may be either the output's index (as string)
// or a sub-string of the port's name.
//
func GetMIDIOutput(pattern string) (gomidi.Out, error) {
	var err error
	var out gomidi.Out
	out, err = getOutputByIndex(fmt.Sprintf("%s", pattern))
	if err != nil {
		out, err = getOutputByName(pattern)
	}
	return out, err
}


func getInputByIndex(s string) (gomidi.In, error) {
	var err error
	var index int
	ins, _ := gomidi.Ins()
	index, err = strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	if index < 0 || len(ins) <= index {
		err = fmt.Errorf("MIDI input index in of bounds: %d", index)
		return nil, err
	}
	return ins[index], nil
}
		

func getInputByName(pattern string) (gomidi.In, error) {
	var err error
	ins, _ := gomidi.Ins()
	for _, port := range ins {
		name := port.String()
		if strings.Contains(name, pattern) {
			return port, nil
		}
	}
	err = fmt.Errorf("Pattern '%s' did not match any MIDI device", pattern)
	return nil, err
}
		

// GetMIDIInput returns selected MIDI input port.
// pattern may be either the output's index (as string)
// or a sub-string of the port's name.
//
func GetMIDIInput(pattern string) (gomidi.In, error) {
	var err error
	var in gomidi.In
	in, err = getInputByIndex(fmt.Sprintf("%s", pattern))
	if err != nil {
		in, err = getInputByName(pattern)
	}
	return in, err
}
	
