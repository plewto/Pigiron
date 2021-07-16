package smf

// NOTE: These test are executed prior to smf_test.go
// If smf_test.go fails it is likely event_test.go will likewise fail
//


import (
	"testing"
	"fmt"
	"github.com/plewto/pigiron/pigpath"
)

// func noResourcesAbort(fnName string) bool {
// 	_, err := pigpath.ResourceFilename("testFiles", "a.mid")
// 	if err != nil {
// 		fmt.Printf("\nWARNING: Can not read resource file required for %s\n", fnName)
// 		fmt.Println("WARNING: Error from pigpath.ResourceFilename was:")
// 		fmt.Printf("%s\n", err)
// 		fmt.Printf("WARNING: Aborting test.\n\n")
// 		return true
// 	}
// 	return false
// }


func TestCreateEventList(t *testing.T) {
	var err error
	var smf *SMF
	var track *Track
	var events *EventList
	var division int
	filename := pigpath.ResourceFilename("testFiles", "no-clocks.mid")
	fmt.Printf("test MIDI file is %s\n", filename)
	smf, err = ReadSMF(filename)
	if err != nil {
		msg := "smf.ReadSMF returned unexpected error for file %s\n"
		err = compoundError(err, fmt.Sprintf(msg, filename))
		t.Fatal(err)
	}
	division, _  = smf.ClockDivision()
	track, err = smf.GetTrack(0)
	// track.Dump()
	if err != nil {
		t.Fatal(err)
	}
	events, err = createEventList(division, track.Bytes())
	events.Dump()
	if err != nil {
		fmt.Println(err)
	}
	ignore(events)
}
