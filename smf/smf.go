package smf

/*
** smf.go defines general structure for Standard MIDI Files.
**
*/

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pigpath"
	"github.com/plewto/pigiron/pigerr"
)

// SMF struct defines a Standard MIDI File.
//
type SMF struct {
	filename string
	header *Header
	tracks[] Track
}

func NewSMF() *SMF {
	smf := new(SMF)
	smf.filename = ""
	smf.header = &Header{0, 1, 24}
	smf.tracks = make([]Track, 0, 1)
	return smf
}


func (smf *SMF) String() string {
	return fmt.Sprintf("SMF '%s'", smf.filename)
}

func (smf *SMF) Dump(verbose ...int) {
	fmt.Println("SMF")
	fmt.Printf("\tFilename : \"%s\"\n", smf.filename)
	fmt.Println("\tHeader")
	fmt.Printf("\t\tformat   : %d\n", smf.Format())
	fmt.Printf("\t\tdivision : %d\n", smf.Division())
	fmt.Printf("\t\ttracks   : %d\n", smf.TrackCount())
	if len(verbose) == 0 {
		return
	}
	switch verbose[0] {
	case 1:  // list track events only
		fmt.Println("\tTracks:")
		for i, trk := range smf.tracks {
			fmt.Printf("\t\t[trk %d] %d events\n", i, len(trk.events))
		}
	case 2: // list full track details
		fmt.Println("\tTracks:")
		for i, trk := range smf.tracks {
			fmt.Printf("\t\t[trk %d] %d events\n", i, len(trk.events))
			for j, ev := range trk.events {
				fmt.Printf("\t\t\t[%5d] %s\n", j, ev.String())
			}
		}
	default: // ignor
	}
}

func (smf *SMF) Format() int {
	return smf.header.format
}

func (smf *SMF) Division() int {
	return smf.header.division
}

func (smf *SMF) TrackCount() int {
	return len(smf.tracks)
}

func (smf *SMF) Track(n int) (track Track, err error) {
	if n < 0 || smf.TrackCount() <= n {
		err = fmt.Errorf("SMF Track number out of bounds: %d", n)
		return
	}
	track = smf.tracks[n]
	return
}

func (smf *SMF) Filename() string {
	return smf.filename
}

func ReadSMF(filename string) (smf *SMF, err error) {
	filename = pigpath.SubSpecialDirectories(filename)
	file, ferr := os.Open(filename)
	if ferr != nil {
		errmsg := "Can not open SMF file: '%s'\n%s"
		err = fmt.Errorf(errmsg, filename, ferr.Error())
		return
	}
	defer file.Close()
	smf = NewSMF()
	smf.header, err = readHeader(file)
	if err != nil {
		return smf, err
	}
	smf.tracks = make([]Track, 0, smf.header.trackCount)
	for i := 0; i < smf.header.trackCount; i++ {
		track, terr := readTrack(file)
		if terr != nil {
			errmsg := "Can not read track %d of smf file %s\n%s"
			err = fmt.Errorf(errmsg, i, filename, terr.Error)
			return
		}
		smf.tracks = append(smf.tracks, *track)
	}
	smf.filename = filename
	return
}

// TickDuration calculates duration of single clock tick.
// Args:
//   division is smf clock Division.
//   tempo in BPM.
//
func TickDuration(division int, tempo float64) float64 {
	division = division & 0x7FFFF
	if tempo == 0 {
		dflt := 60.0
		errmsg := "MIDI tempo is 0, using default %f"
		pigerr.Warning(fmt.Sprintf(errmsg, dflt))
		tempo = dflt
	}
	var qdur float64 = 60.0/tempo
	return qdur/float64(division)
}
	
// smf.Duration returns aproximate duration 0f track 0 in seconds.
//
func (smf *SMF) Duration() float64 {
	if len(smf.tracks) == 0 {
		return 0.0
	}
	var acc float64 = 0.0
	var tempo float64 = 120
	var tick = TickDuration(smf.Division(), tempo)
	var track = smf.tracks[0]
	for _, event := range track.events {
		msg := event.Message()
		if midi.IsTempoChange(msg) {
			tempo, _ := midi.MetaTempoBPM(msg)
			tick = TickDuration(smf.Division(), tempo)
		}
		acc += float64(event.deltaTime) * tick
	}
	return acc
}
			
			
		
		
		
		
			
		
