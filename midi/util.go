package midi

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

const (
	CHANNEL_MASK (int64) = 0xF0
	STATUS_MASK (int64) = 0x0F
	SYSEX (int64) = 0xF0
)

var mnemonics = make(map[int64]string)

func init() {
	mnemonics[0x80] = "NOTE-OFF"   
	mnemonics[0x90] = "NOTE-ON "    
	mnemonics[0xA0] = "TOUCH   "      
	mnemonics[0xB0] = "CTRL    "       
	mnemonics[0xC0] = "PROG    "       
	mnemonics[0xD0] = "PRESS   "      
	mnemonics[0xE0] = "PITCH   "      
	mnemonics[0xF0] = "SYSEX   "      
	mnemonics[0xF1] = "TIMECODE"   
	mnemonics[0xF2] = "SONG-POS"        
	mnemonics[0xF3] = "SONG-SEL"    
	mnemonics[0xF6] = "TUNE    "       
	mnemonics[0xF8] = "CLOCK   "      
	mnemonics[0xFA] = "START   "      
	mnemonics[0xFB] = "CONTINUE"   
	mnemonics[0xFC] = "STOP    "       
	mnemonics[0xFE] = "ASENS   "      
	mnemonics[0xFF] = "RESET   "  
}
	

func StatusMnemonic(n int64) string {
	var mn string
	if n >= 0xF0 {
		mn, _  = mnemonics[n]
	} else {
		mn, _  = mnemonics[n & CHANNEL_MASK]
	}
	return mn
}

func IsChannelStatus(s int64) bool {
	return s & CHANNEL_MASK != 0xF0
}

func IsSystemStatus(s int64) bool {
	return s >= 0xF0
}

func StatusChannelCommand(s int64) int64 {
	return s & CHANNEL_MASK
}

func StatusChannelIndex(s int64) MIDIChannelIndex {
	return MIDIChannelIndex(s & STATUS_MASK)
}

func SetStatusChannelIndex(s int64, ci MIDIChannelIndex) int64 {
	return StatusChannelCommand(s) | int64(ci)
}

func ChannelEventToString(event portmidi.Event) string {
	s := event.Status
	cmd := s & CHANNEL_MASK
	c := (s & STATUS_MASK) + 1
	scmd, _ := mnemonics[cmd]
	str := fmt.Sprintf("%s chan %2d  data: %2x %2x", scmd, c, event.Data1, event.Data2)
	return str
}


func SystemEventToString(event portmidi.Event) string {
	s := event.Status
	acc := fmt.Sprintf("%s ", mnemonics[s])
	if s == SYSEX {
		bytes := event.SysEx
		for i:=0; i<len(bytes) && i<10; i++ {
			acc += fmt.Sprintf("%02x ", bytes[i])
		}
		if len(bytes) > 10 {
			acc += "..."
		}
	} else {
		acc += fmt.Sprintf("%02x, %02x", event.Data1, event.Data2)
	}
	return acc
		
}

func EventToString(event portmidi.Event) string {
	if IsChannelStatus(event.Status) {
		return ChannelEventToString(event)
	} else {
		return SystemEventToString(event)
	}
}
