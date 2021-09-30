package backend

import (
	"fmt"
	"strings"
	"strconv"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
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


func PrintInputs() {
	fmt.Println("MIDI Inputs (rtmidi):")
	for i, s := range InputNames() {
		fmt.Printf("\t[%2d] %s\n", i, s)
	}
	fmt.Println()
}

func OutputNames() []string {
	acc := make([]string,0, 8)
	outs, _ := gomidi.Outs()
	for _, out := range outs {
		acc = append(acc, out.String())
	}
	return acc
}

func PrintOutputs() {
	fmt.Println("MIDI Outputs (rtmidi):")
	for i, n := range OutputNames() {
		fmt.Printf("\t[%2d] %s\n", i, n)
	}
	fmt.Println()
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
		

// pattern may be one of:
//   int index (as string)
//   substrng of device name.
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
	
