package midi

import (
	"fmt"
	"errors"
	"strings"
	
	gomidi "gitlab.com/gomidi/midi"
	driver "gitlab.com/gomidi/portmididrv"
	// driver "gitlab.com/gomidi/rtmididrv"       // TODO: Conditional import ?
)


var backend gomidi.Driver

func must(err error, msg string) {
	if err != nil {
		fmt.Printf("PANIC: %s\n", msg)
		panic(err)
	}
}


func init() {
	var err error
	backend, err = driver.New()
	must(err, "Can not initialize MIDi backend")
}

func inputs() []gomidi.In {
	result, err := backend.Ins()
	must(err, "Can not obtain MIDI inputs")
	return result
}

func outputs() []gomidi.Out {
	result, err := backend.Outs()
	must(err, "Can not obtain MIDI outputs")
	return result
}

func InputNames() []string {
	ary := inputs()
	var result []string = make([]string, 0, len(ary))
	for _, dev := range ary {
		result = append(result, fmt.Sprintf("%v", dev))
	}
	return result
}

func OutputNames() []string {
	ary := outputs()
	var result []string = make([]string, 0, len(ary))
	for _, dev := range ary {
		result = append(result, fmt.Sprintf("%v", dev))
	}
	return result
}

func DumpDevices() {
	fmt.Println("MIDI Inputs:")
	for _, name := range InputNames() {
		fmt.Printf("\t%s\n", name)
	}
	fmt.Println("MIDI Outputs:")
	for _, name := range OutputNames() {
		fmt.Printf("\t%s\n", name)
	}
}


func findDeviceIndex(pattern string, names []string) int {
	result := -1
	for i, n := range names {
		if strings.Contains(n, pattern) {
			result = i
			break
		}
	}
	return result
}
		
func GetInputDeviceName(pattern string) (string, error) {
	var err error
	names := InputNames()
	index := findDeviceIndex(pattern, names)
	if index == -1 {
		msg := fmt.Sprintf("No matching MIDI input: %s", pattern)
		err = errors.New(msg)
		return "", err
	} else {
		return names[index], nil
	}
}


func GetInputDevice(pattern string) (gomidi.In, error) {
	var device gomidi.In
	var err error
	index := findDeviceIndex(pattern, InputNames())
	if index >= 0 {
		device = inputs()[index]
	} else {
		msg := fmt.Sprintf("No matching MIDI input: %s", pattern)
		err = errors.New(msg)
	}
	return device, err
}

	
		

func GetOutputDeviceName(pattern string) (string, error) {
	var err error
	names := OutputNames()
	index := findDeviceIndex(pattern, names)
	if index == -1 {
		msg := fmt.Sprintf("No matching MIDI output: %s", pattern)
		err = errors.New(msg)
		return "", err
	} else {
		return names[index], nil
	}
}
	
func GetOutputDevice(pattern string) (gomidi.Out, error) {
	var device gomidi.Out
	var err error
	index := findDeviceIndex(pattern, OutputNames())
	if index >= 0 {
		device = outputs()[index]
	} else {
		msg := fmt.Sprintf("No matching MIDI output: %s", pattern)
		err = errors.New(msg)
	}
	return device, err
}	

func Cleanup() {
	backend.Close()
}
