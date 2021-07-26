package midi

/*
** devices.go provides an interface to portmidi MIDI devices.
**
*/

import (
	"fmt"
	"strconv"
	"strings"
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


// Cleanup() terminates portmidi.
// Only call on application exit.
//
func Cleanup() {
	err := portmidi.Terminate()
	if err != nil {
		fmt.Println("Could not terminate portmidi")
	}
}


// InputIDs() returns list of portmidi.DeviceID for all MIDI inputs.
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


// OutputIDs() returns list of portmidi.DevicveID for all MIDI outputs.
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


// deviceInfoString() returns string representation for MIDI device.
//
func deviceInfoString(info *portmidi.DeviceInfo) string {
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
		fmt.Printf("\t[id = %2d] %s\n", id, deviceInfoString(info))
	}
}


func dumpOutputs() {
	fmt.Println("Portmidi Outputs:")
	for _, id := range OutputIDs() {
		info := portmidi.Info(id)
		fmt.Printf("\t[id = %2d] %s\n", id, deviceInfoString(info))
	}
}


// DumpDevices() displays list for all portmidi devices.
//
func DumpDevices() {
	dumpInputs()
	dumpOutputs()
}

// getDeviceIdByIndex() returns portmidi DeviceID by location in device list.
//
// Returns non-nil error if string s can not be parsed as an int, or if it is
// out of bounds.
//
func getDeviceIdByIndex(s string, idList []portmidi.DeviceID) (portmidi.DeviceID, error) {
	var err error
	var index int
	var limit = len(idList)
	var id portmidi.DeviceID
	index, err = strconv.Atoi(s)
	if err != nil {
		return id, err
	}
	if limit < 0 || limit >= limit {
		err := fmt.Errorf("device id index out of bounds: %d", index)
		return id, err
	}
	id = idList[index]
	return id, err
}

// getDeviceIdByPattern() returns first id from list with matching pattern.
//
// The first device id which contains pattern as a sub-string is a match.
// Returns non-nil error if no matches are found.
//
func getDeviceIdByPattern(pattern string, idList []portmidi.DeviceID) (portmidi.DeviceID, error) {
	var err error
	var id portmidi.DeviceID
	for _, item := range idList {
		name := portmidi.Info(item).Name
		if strings.Contains(name, pattern) {
			id = item
			return id, err
		}
	}
	err = fmt.Errorf("Pattern '%s' did not match any MIDI device name", pattern)
	return id, err
}
		
// GetInputID() selects portmidi input DeviceID by either name or index.
//
// pattern may either be an integer index, or a sub-string of a device name.
// If no matches are found returns the portmidi default-id and a non-nil error.
//
func GetInputID(pattern string) (portmidi.DeviceID, error) {
	var err error
	var id portmidi.DeviceID
	var idList = InputIDs()
	id, err = getDeviceIdByIndex(pattern, idList)
	if err != nil {
		id, err = getDeviceIdByPattern(pattern, idList)
		if err != nil {
			errmsg := "Pattern '%s' did not match any MIDI input, using default"
			err = fmt.Errorf(errmsg, pattern)
			id = portmidi.DefaultInputDeviceID()
		}
	}
	return id, err
}

	
// GetOutputID() selects portmidi output DeviceID by either name or index.
//
// pattern may either be an integer index, or a sub-string of a device name.
// If no matches are found returns the portmidi default-id and a non-nil error.
//
func GetOutputID(pattern string) (portmidi.DeviceID, error) {
	var err error
	var id portmidi.DeviceID
	var idList = OutputIDs()
	id, err = getDeviceIdByIndex(pattern, idList)
	if err != nil {
		id, err = getDeviceIdByPattern(pattern, idList)
		if err != nil {
			errmsg := "Pattern '%s' did not match any MIDI output, using default"
			err = fmt.Errorf(errmsg, pattern)
			id = portmidi.DefaultOutputDeviceID()
		}
	}
	return id, err
}
