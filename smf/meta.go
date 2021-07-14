package smf

/*
 * Defines MetaMesasge implementation of MIDIMessage
 *
 */

import (
	"fmt"
	"math"
	"github.com/rakyll/portmidi"
)


// MetaType defines valid meta message type values.
//
type MetaType byte

const (
 	MetaSequenceNumber MetaType = 0x00  
 	MetaText MetaType = 0x01
 	MetaCopyright MetaType = 0x02
 	MetaTrackName MetaType = 0x03
 	MetaInstrumentName MetaType = 0x04
 	MetaLyric MetaType = 0x05
 	MetaMarker MetaType = 0x06
 	MetaCuePoint MetaType = 0x07
 	MetaChannelPrefix MetaType = 0x20
 	MetaEndOfTrack MetaType = 0x2F
 	MetaTempo MetaType = 0x51
 	MetaSMPTE MetaType = 0x54
 	MetaTimeSignature MetaType = 0x58
 	MetaKeySignature MetaType = 0x59    
 	MetaSequencerEvent MetaType = 0x7f  
)

var (
	metaTextTypes map[MetaType]bool = map[MetaType]bool{
		MetaText: true,
		MetaCopyright: true,
		MetaTrackName: true,
		MetaInstrumentName: true,
		MetaLyric: true,
		MetaMarker: true,
		MetaCuePoint: true,
	}

	// metaTypeTable maps MetaType to string mnemonic.
	//
	metaTypeTable map[MetaType]string = map[MetaType]string {
		MetaSequenceNumber: "SeqNum ",
		MetaText:           "Text   ",
		MetaCopyright:      "CpyRite",
		MetaTrackName:      "TrkName",
		MetaInstrumentName: "InsName",
		MetaLyric:          "Lyric  ",
		MetaMarker:         "Marker ",
		MetaCuePoint:       "Cue    ",
		MetaChannelPrefix:  "ChanPre",
		MetaEndOfTrack:     "EOT    ",
		MetaTempo:          "Tempo  ",
		MetaSMPTE:          "SMPTE  ",
		MetaTimeSignature:  "TSig   ",
		MetaKeySignature:   "KSig   ",
		MetaSequencerEvent: "SeqEvnt",
	}
)

func (mt MetaType) String() string {
	s, flag := metaTypeTable[mt]
	if !flag {
		s = "????"
	}
	return s
}


// isMetaType returns true if n is a valid meta type.
//
func isMetaType(n byte) bool {
	_, flag := metaTypeTable[MetaType(n)]
	return flag
}

// isMetaTextType returns true if n is one of the text based meta types.
//
func isMetaTextType(n byte) bool {
	_, flag := metaTextTypes[MetaType(n)]
	return flag
}


// Implements MIDIMessage
//
type MetaMessage struct {
	mtype MetaType
	data []byte
}

// Status returns meta status byte (always 0xFF).
//
func (m *MetaMessage) Status() StatusByte {
	return MetaStatus
}

// Bytes returns the byte sequence corresponding to this MetaMessage.
//
func (m *MetaMessage) Bytes() []byte { 
	prefix := []byte{byte(MetaStatus), byte(m.mtype)}
	vlq := NewVLQ(len(m.data))
	acc := append(prefix, vlq.Bytes()...)
	acc = append(acc, m.data...)
	return acc
}
		

// MetaType returns byte indicating the type of this MetaMessage.
//
func (m *MetaMessage) MetaType() MetaType {
	return m.mtype
}


func metaTempoToString(mm *MetaMessage) string {
	acc := ""
	if len(mm.data) != 3 {
		acc += fmt.Sprintf("<malformed, expected 3 data bytes, got %d>", len(mm.data))
		return acc
	}
	usec, _ := get3Bytes(mm.data, 0)
	acc += fmt.Sprintf("%d μSec,   ", usec)
	if usec != 0 {
		bpm := 60000000.0 / float64(usec)
		acc += fmt.Sprintf("%7.3f BPM", bpm)
	}
	return acc
}

func metaTimeSigToString(mm *MetaMessage) string {
	acc := ""
	if len(mm.data) != 4 {
		acc += fmt.Sprintf("<malformed, expected 4 data bytes, got %d>", len(mm.data))
		return acc
	}
	num := mm.data[0]
	exp := float64(mm.data[1])
	den := int(math.Pow(2, exp))
	acc += fmt.Sprintf("%d/%d", num, den)
	return acc
}


func (mm *MetaMessage) String() string {
	mtype := mm.mtype
	mnemonic, flag := metaTypeTable[mtype]
	if !flag {
		mnemonic = "?????  "
	}
	acc := fmt.Sprintf("META %s : ", mnemonic)
	if len(mm.data) > 0 {
		maxbytes := 8
		acc += "["
		for i := 0; i < len(mm.data) && i < maxbytes; i++ {
			acc += fmt.Sprintf("%02X ", mm.data[i])
		}
		if len(mm.data) > maxbytes {
			acc += fmt.Sprintf("... %d more] ", len(mm.data) - maxbytes)
		} else {
			acc += "] "
		}
	}
	switch {
	case isMetaTextType(byte(mtype)):
		acc += fmt.Sprintf("\"%s\"", string(mm.data))
	case mtype == MetaTempo:
		acc += metaTempoToString(mm)
	case mtype == MetaTimeSignature:
		acc += metaTimeSigToString(mm)
	default:
		// ignore
	}
	return acc
}

func (m *MetaMessage) Data() []byte {
	return m.data
}

func (m *MetaMessage) ConvertsToPortmidi() bool {
	return false
}

func (m *MetaMessage) ToPortmidiEvent() (portmidi.Event, error) {
	var pme portmidi.Event
	msg := "MetaMessage can not be converted to portmidi.Event. "
	msg += "MetaType is 0x%x '%s'"
	err := fmt.Errorf(msg, byte(m.MetaType()), m.MetaType())
	return pme, err
}

func newMetaMessage(mtype MetaType, data []byte)(*MetaMessage, error) {
	var err error
	var meta *MetaMessage
	if !isMetaType(byte(mtype)) {
		msg := "smf.newMetaMessage, invalid MetaType 0x%x"
		err = fmt.Errorf(msg, byte(mtype))
		return meta, err
	}
	// ISSUE: TODO validate data
	meta = &MetaMessage{mtype, data}
	return meta, err
}
