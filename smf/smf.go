package smf

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigpath"
)

func init (
	fmt.Println("*** smf package has been depreciated. ****")
)


func ignore(...interface{}){} // For testing only

func exError(message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	return fmt.Errorf(msg)
}

func compoundError(err1 error, message string, more ...string) error {
	msg := fmt.Sprintf("ERROR: %s\n", message)
	for _, s := range more {
		msg += fmt.Sprintf("\nERROR: %s", s)
	}
	msg += fmt.Sprintf("PREVIOUS ERROR:\n%s", err1)
	return fmt.Errorf(msg)
}


// SMF struct Represents a Standard MIDI File.
//
type SMF struct {
	filename string
	header *Header
	tracks []*Track
}


func (smf *SMF) ClockDivision() (int, error) {
	bytes := smf.header.Bytes()
	division, err := getShort(bytes, 12)
	if err != nil {
		errmsg := "*smf.ClockDivision error"
		err = compoundError(err, errmsg)
		return division, err
	}
	if division < 12 {
		errmsg := "*smf.ClockDivision looks weird: %d"
		err = exError(fmt.Sprintf(errmsg, division))
		return division, err
	}
	return division, err
}

func (smf *SMF) TrackCount() int {
	return len(smf.tracks)
}

func (smf *SMF) GetTrack(n int) (*Track, error) {
	var err error
	var trk *Track
	tcount := smf.TrackCount()
	if n >= tcount {
		errmsg := "*smf.GetTrack index out of bounds\n"
		errmsg += "Track number %d >= track count %d"
		err = exError(fmt.Sprintf(errmsg, n, tcount))
		return trk, err
	}
	trk = smf.tracks[n]
	return trk, err
}


func ReadSMF(filename string) (*SMF, error) {
	var err error
	var smf = &SMF{}
	smf.filename = pigpath.SubSpecialDirectories(filename)
	var header *Header
	var file *os.File
	file, err = os.Open(filename)
	errmsg := "smf.ReadSMF can not open MIDI file\n"
	errmsg += fmt.Sprintf("filename is '%s'", filename)
	if err != nil {
		err2 := compoundError(err, errmsg)
		return smf, err2
	}
	defer file.Close()
	header, err = readHeader(file)
	if err != nil {
		err = compoundError(err, errmsg)
		return smf, err
	}
	smf.header = header
	for i := 0; i < header.chunkCount; i++ {
		var id [4]byte
		var length int
		id, length, err = readChunkPreamble(file)
		if err != nil {
			err = compoundError(err, errmsg)
			return smf, err
		}
		if id != trackID { // skip non-track chunks
			msg := "Skiping non-recognized chunk %v, smf file '%s'\n"
			fmt.Printf(msg, string(id[:]), filename)
			var junk = make([]byte, length)
			_, err = file.Read(junk)
			if err != nil {
				sid := string(id[:])
				msg = fmt.Sprintf("smf.ReadSMF file: '%s'", filename)
				msg2 := "An error occured while reading non-recognized chunk: %s"
				err = compoundError(err, msg, fmt.Sprintf(msg2, sid))
				return smf, err
			}
		} else { // read track chunk
			trk := &Track{make([]byte, 0, 1024)}
			trk.bytes, err = readTrackBytes(file, length)
			if err != nil {
				msg := fmt.Sprintf("smf.READSMF file: '%s'", filename)
				msg2 := "An error occured while reading track chunk %d"
				err = compoundError(err, msg, fmt.Sprintf(msg2, i))
				return smf, err
			}
			smf.tracks = append(smf.tracks, trk)
		}
	}
	return smf, err
}
				
func (smf *SMF) Filename() string {
	return smf.filename
}


func (smf *SMF) Dump() {
	fmt.Println("SMF")
	fmt.Printf("filename : '%s'\n", smf.filename)
	smf.header.Dump()
	for i, trk := range smf.tracks {
		fmt.Printf("---- Track [%02d]\n", i)
		trk.Dump()
	}
}
	
