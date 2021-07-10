package smf

import "fmt"

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

	metaTypeTable map[MetaType]string = map[MetaType]string {
		MetaSequenceNumber: "SEQ-NUMBER",
		MetaText: "TEXT",
		MetaCopyright: "COPYRIGHT",
		MetaTrackName: "TRACK-NAME",
		MetaInstrumentName: "INSTRUMENT-NAME",
		MetaLyric: "LYRIC",
		MetaMarker: "MARKER",
		MetaCuePoint: "CUE",
		MetaChannelPrefix: "CHANNEL-PREFIX",
		MetaEndOfTrack: "END-OF-TRACK",
		MetaTempo: "TEMPO",
		MetaSMPTE: "SMPTE",
		MetaTimeSignature: "TIME-SIG",
		MetaKeySignature: "KEY-SIG",
		MetaSequencerEvent: "SEQ-EVENT",
	}
)


func (t MetaType) String() string {
	s, flag := metaTypeTable[t]
	if !flag {
		s = "?"
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
	bytes []byte
}

func (m *MetaMessage) Status() StatusByte {
	return StatusByte(m.bytes[0])
}

func (m *MetaMessage) Bytes() []byte {
	return m.bytes
}

func (m *MetaMessage) String() string {
	if len(m.bytes) < 2 {
		return "MetaMessage ??"
	} else {
		return fmt.Sprintf("MetaMessage %s", MetaType(m.bytes[1]))
	}
}

func (m *MetaMessage) Dump() {
	if m == nil {
		fmt.Println("MetaMessage <nil>")
		return
	}
	fmt.Printf("%s\n", m)
	fmt.Printf("[ 0] 0x%02x   - status\n", m.bytes[0])
	fmt.Printf("[ 1] 0x%02x   - meta type '%s'\n", m.bytes[1], MetaType(m.bytes[1]))
	counter := 0
	for i, b := range m.bytes[2:] {
		fmt.Printf("[%2d] 0x%02x   - vlq-%d\n", i+2, b, i)
		counter++
		if b & 0x80 == 0 {
			break
		}
	}
	isText := isMetaTextType(m.bytes[1])
	start := counter + 2
	for i, b := range m.bytes[start:] {
		fmt.Printf("[%2d] 0x%02x   - data-%d", i+start, b, i)
		if isText {
			fmt.Printf(" '%c'", b)
		}
		fmt.Println()
	}
}
			

func (m *MetaMessage) MetaType() MetaType {
	return MetaType(m.bytes[1])
}


// Returns only data bytes.  IE all values after final VLQ byte.
//
func (m *MetaMessage) Data() []byte {
	index := 2
	bytes := m.bytes
	for index < len(bytes) {
		b := bytes[index]
		index++
		if b & 0x80 == 0 {
			break
		}
	}
	return m.bytes[index:]
}

func newMetaMessage(mtype MetaType, data []byte) *MetaMessage {
	head := []byte{byte(MetaStatus), byte(mtype)}
	vlq := NewVLQ(len(data))
	bytes := append(head, vlq.Bytes()...)
	mm := &MetaMessage{append(bytes, data...)}
	return mm
}

func NewMetaSequenceNumber(n byte) *MetaMessage {
	data := []byte{n}
	return newMetaMessage(MetaSequenceNumber, data)
}


// NOTE: invalid non-text MetaType silently converted to MetaText
//
func NewMetaText(mtype MetaType, text string) *MetaMessage {
	if !isMetaTextType(byte(mtype)) {
		mtype = MetaText
	}
	data := []byte(text)
	return newMetaMessage(mtype, data)
}

func NewMetaChannelPrefix(prefix byte) *MetaMessage {
	data := []byte{prefix}
	return newMetaMessage(MetaChannelPrefix, data)
}

var (
	eot = newMetaMessage(MetaEndOfTrack, []byte{})
)

func NewMetaEndOfTrack() *MetaMessage {
	return eot
}

func NewMetaTempo(usec [4]byte) *MetaMessage {
	return newMetaMessage(MetaTempo, usec[:])
}

func NewMetaSMPTE(values [5]byte) *MetaMessage {
	return newMetaMessage(MetaSMPTE, values[:])
}

func NewMetaTimeSignature(num byte, dnom byte, args... byte) *MetaMessage {
	cpc, n32 := byte(24), byte(8)
	if len(args) > 0 {
		cpc = args[0]
	}
	if len(args) > 1 {
		n32 = args[1]
	}
	return newMetaMessage(MetaTimeSignature, []byte{num, dnom, cpc, n32})
}
		
func NewMetaKeySignature(key byte, minor bool) *MetaMessage {
	var flag = byte(0)
	if minor {
		flag = 1
	}
	return newMetaMessage(MetaKeySignature, []byte{key, flag})
}
	
func NewSequencerSpecificMessage(data []byte) *MetaMessage {
	return newMetaMessage(MetaSequencerEvent, data)
}

