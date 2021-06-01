package midi

import (
	"fmt"
	"strconv"
	"errors"
	
	gomidi "gitlab.com/gomidi/midi"
	// driver "gitlab.com/gomidi/rtmididrv"       // TODO: Conditional import ?
	driver "gitlab.com/gomidi/portmididrv"

)

var drv gomidi.Driver

func must(err error, msg string) {
	if err != nil {
		fmt.Printf("PANIC: %s\n", msg)
		panic(err)
	}
}


func init() {
	var err error
	drv, err = driver.New()
	must(err, "Can not initialize MIDI driver")
}


// inputDevices returns list of system MIDI inputs.
//
func inputDevices() []gomidi.In {
	result, err := drv.Ins()
	must(err, "Can not obtain MIDI inputs")
	return result
}

// outputDevices returns list of system MIDI outputs.
//
func outputDevices() []gomidi.Out {
	result, err := drv.Outs()
	must(err, "Can not obtain MIDI outputs")
	return result
}

// InputDeviceNames returns list of system MIDI input device names.
//
func InputDeviceNames() []string {
	ary := inputDevices()
	var result []string = make([]string, 0, len(ary))
	for _, dev := range ary {
		result = append(result, fmt.Sprintf("%v", dev))
	}
	return result
}


// OutputDeviceNames returns list of system MIDI output device names.
//
func OutputDeviceNames() []string {
	ary := outputDevices()
	var result []string = make([]string, 0, len(ary))
	for _, dev := range ary {
		result = append(result, fmt.Sprintf("%v", dev))
	}
	return result
}


// DumpDevices displays list of system MIDI devices.
//
func DumpDevices() {
	fmt.Println("MIDI Inputs:")
	for i, d := range InputDeviceNames() {
		fmt.Printf("\t[%2d] %s\n", i, d)
	}
	fmt.Println("MIDI Outputs:")
	for i, d := range OutputDeviceNames() {
		fmt.Printf("\t[%2d] %s\n", i, d)
	}
}


// toInt converts string to int and returns it.
// If string can not be converted or out side the closed interval [0, max-1], return -1.
//
func toInt(s string, max int) int {
	a, err := strconv.Atoi(s)
	if err != nil || a < 0 || a >= max {
		return -1
	} else {
		return a
	}
}



// findDeviceIndex locates index of target device within names array.
// target may either be an integer index or an exact match for an element in names.
// If no matching target is found, return -1, error.
//
func findDeviceIndex(target string, names []string) (int, error) {
	var err error
	n := toInt(target, len(names))
	if n != -1 {
		return n, err
	} else {
		for i, n := range names {
			if target == n {
				return i, err
			}
		}
	}
	msg := fmt.Sprintf("MIDI Device %s does not exists", target)
	err = errors.New(msg)
	return -1, err
}
		
// GetInputDevice returns system MIDI input device.
// target may either be an index into the device names array, or the exact
// name of a MIDI device.  Returns error if target does not match any device.
//
func GetInputDevice(target string) (gomidi.In, error) {
	var device gomidi.In
	names := InputDeviceNames()
	index, err := findDeviceIndex(target, names)
	if err == nil {
		device = inputDevices()[index]
	}
	return device, err
}


// GetOutputDevice returns system MIDI output device.
// target may either be an index into the device names array, or the exact
// name of a MIDI device.  Returns error if target does not match any device.
//
func GetOutputDevice(target string) (gomidi.Out, error) {
	var device gomidi.Out
	names := OutputDeviceNames()
	index, err := findDeviceIndex(target, names)
	if err == nil {
		device = outputDevices()[index]
	}
	return device, err
}
		

// Closes backing MIDI driver.
// 
func Cleanup() {
	drv.Close()
}
