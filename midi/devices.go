package midi

import (
	"fmt"
	"strings"
	"errors"
	
	"github.com/rakyll/portmidi"
)


func must(err error, msg string) {
	if err != nil {
		if len(msg) > 0 {
			fmt.Printf("ERROR: %s\n", msg)
		}
		panic(err)
	}
}


func init() {
	must(portmidi.Initialize(),
		"Can not initialize portmidi, nothing can be done without it!")
}


// Cleanup terminates portamidi.
// Only call on application exit.
//
func Cleanup() {
	err := portmidi.Terminate()
	if err != nil {
		fmt.Println("Could not terminate portmidi")
	}
}


// InputIDs returns list of portmidi.DeviceID for all MIDI inputs.
//
func InputIDs() []portmidi.DeviceID {
	maxCount := portmidi.CountDevices()
	acc := make([]portmidi.DeviceID, 0, maxCount)
	for i := 0; i < maxCount; i++ {
		id := portmidi.DeviceID(i)
		info := portmidi.Info(id)
		if info.IsInputAvailable {
			acc = append(acc, id)
		}
	}
	return acc
}


// OutputIDs returns list of portmidi.DevicveID for all MIDI outputs
//
func OutputIDs() []portmidi.DeviceID {
	maxCount := portmidi.CountDevices()
	acc := make([]portmidi.DeviceID, 0, maxCount)
	for i := 0; i < maxCount; i++ {
		id := portmidi.DeviceID(i)
		info := portmidi.Info(id)
		if info.IsOutputAvailable {
			acc = append(acc, id)
		}
	}
	return acc
}


func padString(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}


// ReprDeviceInfo returns formatted string representation for MIDI device.
//
func ReprDeviceInfo(info *portmidi.DeviceInfo) string {
	name := fmt.Sprintf("\"%s\"", info.Name)
	s := fmt.Sprintf(" %s  ", padString(name, 32))
	if info.IsInputAvailable {
		s += "I"
	} else {
		s += " "
	}
	if info.IsOutputAvailable {
		s += "O"
	} else {
		s += " "
	}
	if info.IsOpened {
		s += "  opened"
	} else {
		s += "  closed"
	}
	return s
}


func dumpInputs() {
	fmt.Println("Portmidi Inputs:")
	for _, id := range InputIDs() {
		info := portmidi.Info(id)
		fmt.Printf("\t[id = %2d] %s\n", id, ReprDeviceInfo(info))
	}
}


func dumpOutputs() {
	fmt.Println("Portmidi Outputs:")
	for _, id := range OutputIDs() {
		info := portmidi.Info(id)
		fmt.Printf("\t[id = %2d] %s\n", id, ReprDeviceInfo(info))
	}
}


// DumpDevices displays list for all portmidi MIDI IO devices.
//
func DumpDevices() {
	dumpInputs()
	dumpOutputs()
}


// GetInputID searches for specific portmidi input DeviceID.
// pattern is matched against the device name, The id for the first device
// whose name contains pattern as a sub-string is returned.
//
// If there are no matching devices, the default id is returned with a
// non-nil error.
//
func GetInputID(pattern string) (portmidi.DeviceID, error) {
	var result portmidi.DeviceID
	var err error
	for _, id := range InputIDs() {
		info := portmidi.Info(id)
		name := info.Name
		if strings.Contains(name, pattern) {
			result = id
			return result, err
		}
	}
	msg := fmt.Sprintf("MIDI input device \"%s\" does not exists, using default.", pattern)
	err = errors.New(msg)
	result = portmidi.DefaultInputDeviceID()
	return result, err
}


// GetOutputID searches for specific portmidi output DeviceID.
// pattern is matched against the device name, The id for the first device
// whose name contains pattern as a sub-string is returned.
//
// If there are no matching devices, the default id is returned with a
// non-nil error.
//
func GetOutputID(pattern string) (portmidi.DeviceID, error) {
	var result portmidi.DeviceID
	var err error
	for _, id := range OutputIDs() {
		info := portmidi.Info(id)
		name := info.Name
		if strings.Contains(name, pattern) {
			result = id
			return result, err
		}
	}
	msg := fmt.Sprintf("MIDI output device \"%s\" does not exists, using default.", pattern)
	err = errors.New(msg)
	result = portmidi.DefaultOutputDeviceID()
	return result, err
}
	
		
