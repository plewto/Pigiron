package midi

/*
** smf.go defines general structure for Standard MIDI Files.
**
*/

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigerr"
)

// SMF struct defines a Standard MIDI File.
//
type SMF struct {
	filename string
	header *SMFHeader
	tracks[] *SMFTrack
}

func (smf *SMF) String() string {
	return fmt.Sprintf("SMF '%s'", smf.filename)
}

// smf.Dump() display MIDI file contents.
//
func (smf *SMF) Dump() {
	fmt.Println("SMF")
	fmt.Printf("  filename : %s\n", smf.filename)
	fmt.Printf("  format   : %d\n", smf.header.format)
	fmt.Printf("  tracks   : %d\n", smf.header.trackCount)
	fmt.Printf("  division : %d\n", smf.header.division)
	for i, trk := range smf.tracks {
		fmt.Printf("  TRACK: %d  events: %d\n", i, trk.Length())
		for j, evnt := range trk.events {
			fmt.Printf("  [%3d] %s\n", j, evnt)
		}
	}
}

// smf.Filename() returns the SMF filename.
//
func (smf *SMF) Filename() string {
	return smf.filename
}

// smf.Format() returns the SMF format.
// Three formats are defined.
//  0 - Single MIDI file (currently only format 0 files are supported).
//  1 - Multiple tracks, single song (possibly supported in future).
//  2 - Multiple songs - no plans to support.
//
func (smf *SMF) Format() int {
	return smf.header.format
}

// smf.TrackCount() returns number of tracks.
// For format 0 the result count is always 1.
//
func (smf *SMF) TrackCount() int {
	return len(smf.tracks)
}

// smf.Division() returns the MIDI clock resolution.
//
func (smf *SMF) Division() int {
	return smf.header.Division()
}

// smf.Track() returns the indicated track.
// For format 0, the argument must always be 0.
//
// Returns:
//  1. The indicated track.
//  2. Non-nil error if track number is out of bounds.
//
func (smf *SMF) Track(n int) (track *SMFTrack, err error) {
	if n < 0 || smf.TrackCount() <= n {
		errmsg := "SMF track number out of bounds %d, track count is %d"
		err = pigerr.New(fmt.Sprintf(errmsg, n, smf.TrackCount()))
		return
	}
	track = smf.tracks[n]
	return
}

// TickDuration() returns the duration of a single MIDI clock
//
// Arguments:
//    division - the SMF clock division.
//    tempo - the tempo in BPM.
//
// Returns
//    Tick duration in seconds.
//
func TickDuration(division int, tempo float64) float64 {
	division = division & 0x7FFF
	if tempo == 0 {
		dflt := 60.0
		errmsg := "MIDI tempo is 0, using default %f"
		pigerr.Warning(fmt.Sprintf(errmsg, dflt))
		tempo = dflt
	}
	var qdur float64 = 60.0/tempo
	return qdur/float64(division)
}

// smf.Duration() returns approximate MIDI file duration in seconds.
//
func (smf *SMF) Duration() float64 {
	if len(smf.tracks) == 0 {
		return 0.0
	}
	var tempo float64 = 120
	var tick = TickDuration(smf.Division(), tempo)
	var acc float64 = 0.0
	track := smf.tracks[0]
	for _, evnt := range track.events {
		if byte(evnt.metaType) == byte(META_TEMPO) {
			t, err := evnt.MetaTempoBPM()
			if err != nil {
				break
			}
			tempo = t
			tick = TickDuration(smf.Division(), tempo)
		}
		acc += float64(evnt.deltaTime) * tick
	}
	return acc
}

// ReadSMF() reads SMF from given filename.
//
// Returns:
//   1. *SMF
//   2. Non-nil error if filename could not be read as a MIDI file.
//
func ReadSMF(filename string) (smf *SMF, err error) {
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		errmsg := "Can not open MIDI file, filename = \"%s\""
		err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, filename))
		return
	}
	defer file.Close()
	smf = &SMF{}
	smf.filename = filename
	smf.header, err = readSMFHeader(file)
	// Due to the possible presence of non-track chunks, the
	// header.trackCount may initially be too high.
	// The count is adjusted once all chunks have been read.
	// non-track chunks are not-supported and discarded.
	if err != nil {
		errmsg := "MIDI file header malformed, filename =  %s"
		err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, filename))
		return
	}
	smf.tracks = make([]*SMFTrack, 0, smf.header.trackCount)
	for i := 0; i < smf.header.trackCount; i++ {
		var id chunkID
		var bytes []byte
		var track *SMFTrack
		id, bytes, err = readRawChunk(file)
		if err != nil {
			errmsg := "Error while reading MIDI file chunk %d, filename = %s"
			err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, i, filename))
			return
		}
		if !id.eq(trackID) {
			errmsg := "Ignoring non-track chunk number %d, type %s, filename = %s"
			pigerr.Warning(fmt.Sprintf(errmsg, i, id, filename))
			continue
		}
		track, err = convertTrackBytes(bytes)
		if err != nil {
			errmsg := "Error wile converting track %d, filename = %s"
			err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, i, filename))
			return
		}
		smf.tracks = append(smf.tracks, track)
	}
	smf.header.trackCount = len(smf.tracks)
	smf.filename = filename
	return
}
